package server

import (
	"context"
	"github.com/tgkzz/gateway/config"
	"github.com/tgkzz/gateway/internal/handler"
	"log/slog"
)

type HttpServer struct {
	Handler handler.Handler
}

func NewHttpServer(cfg config.Config, logger *slog.Logger) (*HttpServer, error) {
	h, err := handler.NewEchoHandler(cfg.AuthServer.Port, logger)
	if err != nil {
		return nil, err
	}
	return &HttpServer{
		Handler: h,
	}, nil
}

func (s *HttpServer) Run(port string) error {
	return s.Handler.Start(port)
}

func (s *HttpServer) Stop(ctx context.Context) error {
	return s.Handler.Stop(ctx)
}
