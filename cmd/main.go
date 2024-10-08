package main

import (
	"context"
	"github.com/tgkzz/gateway/config"
	"github.com/tgkzz/gateway/internal/app"
	logger2 "github.com/tgkzz/gateway/pkg/logger"
	"log/slog"
	"os"
	"os/signal"
	"time"
)

const defaultConfigPath = "./gateway.json"

func main() {
	const op = "main.Start"

	c := context.Background()

	ctx, stop := signal.NotifyContext(c, os.Interrupt)
	defer stop()

	cfg := config.MustRead(defaultConfigPath)

	var logger *slog.Logger
	switch cfg.Env != "" {
	case true:
		logger = logger2.SetupLogger(cfg.Env)
	default:
		logger = logger2.SetupLogger("local")
	}

	a, err := app.New(*cfg, logger)
	if err != nil {
		logger.Error("%s: %s", op, err)
	}

	go func() {
		a.MustRun(cfg)
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(c, 10*time.Second)
	defer cancel()
	if err = a.Stop(ctx); err != nil {
		logger.Error("%s: %s", op, err)
	}
}
