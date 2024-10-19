package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/tgkzz/gateway/config"
	"github.com/tgkzz/gateway/internal/app/http/server"
)

type App struct {
	httpServer *server.HttpServer
	logger     *slog.Logger
}

func New(cfg config.Config, logger *slog.Logger) (*App, error) {
	httpServer, err := server.NewHttpServer(cfg, logger)
	if err != nil {
		return nil, fmt.Errorf("new http server: %w", err)
	}

	return &App{
		httpServer: httpServer,
		logger:     logger,
	}, nil
}

func (a *App) MustRun(cfg *config.Config) {
	const op = "app.MustRun"

	log := a.logger.With(
		slog.String("op", op),
	)

	if err := a.httpServer.Run(cfg.Port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("%s: %w", op, err)
		panic(err)
	}
}

func (a *App) Stop(ctx context.Context) error {
	const op = "app.Stop"

	log := a.logger.With(slog.String("op", op))

	if err := a.httpServer.Stop(ctx); err != nil {
		log.Error("%s: %w", op, err)
		return err
	}

	return nil
}
