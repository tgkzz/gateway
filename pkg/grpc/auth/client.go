package auth

import (
	"context"
	"fmt"
	authPkg "github.com/tgkzz/auth/gen/go/auth"
	"github.com/tgkzz/auth/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log/slog"
)

type AuthClient interface {
	Register(ctx context.Context, username, password string) (int, error)
	Login(ctx context.Context, username, password string) (string, error)
}

type Auth struct {
	client authPkg.AuthServiceClient
	logger *slog.Logger
}

func NewAuthClient(port string, logger *slog.Logger) (AuthClient, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf(":%s", port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	c := authPkg.NewAuthServiceClient(conn)

	return &Auth{client: c, logger: logger}, nil
}

func (a *Auth) Register(ctx context.Context, username, password string) (int, error) {
	const op = "grpcAuthService.Register"

	log := a.logger.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	resp, err := a.client.Register(ctx, &authPkg.RegisterRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Error("failed to register", logger.Err(err))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return int(resp.GetUserId()), nil
}

func (a *Auth) Login(ctx context.Context, username, password string) (string, error) {
	const op = "grpcAuthService.Login"

	log := a.logger.With(
		slog.String("op", op),
		slog.String("username", username),
	)

	resp, err := a.client.Login(ctx, &authPkg.LoginRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		log.Error("failed to login", logger.Err(err))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.GetToken(), nil
}
