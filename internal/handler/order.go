package handler

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tgkzz/gateway/internal/model"
	"github.com/tgkzz/gateway/internal/service/order"
)

func (eh *EchoHandler) createOrder(c echo.Context) error {
	var req model.Order
	if err := c.Bind(&req); err != nil {
		return c.String(http.StatusBadRequest, CouldNotReadBody)
	}

	orderId, err := eh.orderService.CreateOrder(c.Request().Context(), req)
	if err != nil {
		if errors.Is(err, order.ErrInvalidArguments) {
			return c.String(http.StatusBadRequest, err.Error())
		}
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
		if errors.Is(err, order.ErrOrderNotFound) {
			return c.String(http.StatusNotFound, err.Error())
		}
		return eh.errorHandler(err)
	}

	return c.JSONPretty(http.StatusOK, res, "  ")
}
