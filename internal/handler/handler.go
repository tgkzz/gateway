package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/tgkzz/gateway/config"
	"github.com/tgkzz/gateway/internal/service/auth"
	"github.com/tgkzz/gateway/internal/service/order"
	"golang.org/x/time/rate"
)

type Handler interface {
	// server
	Start(port string) error
	Stop(ctx context.Context) error

	// order
	createOrder(c echo.Context) error
	getOrder(c echo.Context) error

	// auth
	login(ctx echo.Context) error
	register(ctx echo.Context) error
}

type EchoHandler struct {
	authService  auth.IAuthService
	orderService order.IOrderService
	logger       *slog.Logger
	echoInstance *echo.Echo
}

func NewEchoHandler(cfg config.Config, logger *slog.Logger) (Handler, error) {
	authService, err := auth.NewAuthService(cfg.AuthServer.Host, cfg.AuthServer.Port, logger)
	if err != nil {
		return nil, err
	}
	orderService, err := order.NewOrderService(cfg.OrderServer.Host, cfg.OrderServer.Port, logger)
	if err != nil {
		return nil, err
	}
	return &EchoHandler{
		authService:  authService,
		logger:       logger,
		orderService: orderService,
	}, nil
}

func (eh *EchoHandler) Stop(ctx context.Context) error {
	if eh.echoInstance == nil {
		return errors.New("echo instance not initialized")
	}

	if err := eh.echoInstance.Shutdown(ctx); err != nil {
		return err
	}

	return nil
}

func (eh *EchoHandler) Start(port string) error {
	e := eh.routes()

	eh.echoInstance = e

	return e.Start(":" + port)
}

func (eh *EchoHandler) routes() *echo.Echo {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		Skipper:   middleware.DefaultSkipper,
		StackSize: 8 << 10,
		LogLevel:  log.ERROR,
	}))

	e.Use(middleware.RateLimiterWithConfig(middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: rate.Limit(10), Burst: 50, ExpiresIn: 5 * time.Minute},
		),
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}))

	// prometheus?

	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		Timeout:      time.Second * 30,
		ErrorMessage: RequestTimeout,
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			op := fmt.Sprintf("%s", "handler.timeout "+c.Path())
			deadline, _ := c.Request().Context().Deadline()

			l := eh.logger.With(
				slog.String("op", op),
				slog.Time("now", time.Now().UTC()),
				slog.Time("deadline", deadline),
			)

			l.Error("request timed out for path", slog.String("path", c.Path()))
		},
	}))

	e.Use(middleware.Secure())

	v1 := e.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", eh.login)
			auth.POST("/register", eh.register)
		}
		order := v1.Group("/order")
		{
			order.POST("/create", eh.createOrder)
			order.GET("/get/:id", eh.getOrder)
		}
	}

	return e
}

func (eh *EchoHandler) errorHandler(err error) error {
	return echo.NewHTTPError(http.StatusInternalServerError, err)
}
