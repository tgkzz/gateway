package handler

import (
	"context"
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/tgkzz/gateway/internal/service/auth"
	"log/slog"
	"net/http"
)

type Handler interface {
	Login(ctx echo.Context) error
	Register(ctx echo.Context) error
	Start(port string) error
	Stop(ctx context.Context) error
}

type EchoHandler struct {
	authService  auth.IAuthService
	echoInstance *echo.Echo
}

func NewEchoHandler(authPort string, logger *slog.Logger) (Handler, error) {
	authService, err := auth.NewAuthService(authPort, logger)
	if err != nil {
		return nil, err
	}
	return &EchoHandler{authService: authService}, nil
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
	e := echo.New()

	v1 := e.Group("/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", eh.Login)
			auth.POST("/register", eh.Register)
		}
	}

	eh.echoInstance = e

	return e.Start(":" + port)
}

func (eh *EchoHandler) Login(ctx echo.Context) error {
	username := ctx.FormValue("username")
	if username == "" {
		return eh.errorHandler(fmt.Errorf("%s", UsernameIsEmpty))
	}

	password := ctx.FormValue("password")
	if password == "" {
		return eh.errorHandler(fmt.Errorf("%s", PasswordIsEmpty))
	}

	token, err := eh.authService.Login(ctx.Request().Context(), username, password)
	if err != nil {
		return eh.errorHandler(err)
	}

	return ctx.JSON(http.StatusNonAuthoritativeInfo, token)
}

func (eh *EchoHandler) Register(ctx echo.Context) error {
	username := ctx.FormValue("username")
	if username == "" {
		return eh.errorHandler(fmt.Errorf("%s", UsernameIsEmpty))
	}

	password := ctx.FormValue("password")
	if password == "" {
		return eh.errorHandler(fmt.Errorf("%s", PasswordIsEmpty))
	}

	if err := eh.authService.CreateNewUser(ctx.Request().Context(), username, password); err != nil {
		return eh.errorHandler(err)
	}

	return ctx.NoContent(http.StatusCreated)
}

func (eh *EchoHandler) errorHandler(err error) error {
	return echo.NewHTTPError(http.StatusInternalServerError, err)
}
