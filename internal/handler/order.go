package handler

import (
	"github.com/labstack/echo/v4"
	"github.com/tgkzz/gateway/internal/model"
	"net/http"
)

func (eh *EchoHandler) createOrder(c echo.Context) error {
	var req model.Order
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, CouldNotReadBody)
	}

	return c.NoContent(http.StatusNotImplemented)
}

func (eh *EchoHandler) getOrder(c echo.Context) error {
	res := model.Order{}

	return c.JSONPretty(http.StatusOK, res, "  ")
}
