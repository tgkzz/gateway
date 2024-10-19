package auth

import (
	"context"
	"log/slog"

	"github.com/tgkzz/gateway/pkg/grpc/auth"
)

type IAuthService interface {
	Login(ctx context.Context, login, password string) (string, error)
	CreateNewUser(ctx context.Context, username, psw string) error
}

type AuthService struct {
	cli auth.AuthClient
}

func NewAuthService(authHost, authPort string, logger *slog.Logger) (IAuthService, error) {
	authCli, err := auth.NewAuthClient(authHost, authPort, logger)
	if err != nil {
		return nil, err
	}

	return &AuthService{cli: authCli}, nil
}

func (a *AuthService) Login(ctx context.Context, login, password string) (string, error) {
	t, err := a.cli.Login(ctx, login, password)
	if err != nil {
		return "", err
	}
	return t, nil
}

func (a *AuthService) CreateNewUser(ctx context.Context, username, psw string) error {
	if _, err := a.cli.Register(ctx, username, psw); err != nil {
		return err
	}

	return nil
}
