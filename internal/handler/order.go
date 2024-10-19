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

	orderId, err := eh.orderService.CreateOrder(c.Request().Context(), req)
	if err != nil {
		return eh.errorHandler(err)
	}

	type response struct {
		OrderId string `json:"orderId"`
	}

	return c.JSON(http.StatusCreated, response{
		OrderId: orderId,
	})
}

func (eh *EchoHandler) getOrder(c echo.Context) error {
	orderId := c.Param("orderId")

	res, err := eh.orderService.GetOrderById(c.Request().Context(), orderId)
	if err != nil {
		return eh.errorHandler(err)
	}

	return c.JSONPretty(http.StatusOK, res, "  ")
}
