package handler

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (eh *EchoHandler) login(ctx echo.Context) error {
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

func (eh *EchoHandler) register(ctx echo.Context) error {
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
