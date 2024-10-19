package order

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/tgkzz/gateway/internal/model"
	"github.com/tgkzz/gateway/pkg/grpc/order"
	"github.com/tgkzz/gateway/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type IOrderService interface {
	CreateOrder(c context.Context, order model.Order) (string, error)
	GetOrderById(c context.Context, orderId string) (model.Order, error)
}

type OrderService struct {
	cli    order.OrderClient
	logger *slog.Logger
}

var (
	ErrOrderNotFound    = errors.New("order not found")
	ErrInvalidArguments = errors.New("invalid arguments")
)

func NewOrderService(orderHost, orderPort string, logger *slog.Logger) (IOrderService, error) {
	orderCli, err := order.NewOrderClient(orderHost, orderPort, logger)
	if err != nil {
		return nil, err
	}

	return &OrderService{cli: orderCli, logger: logger}, nil
}

func (or *OrderService) CreateOrder(c context.Context, order model.Order) (string, error) {
	const op = "order.CreateOrder"

	log := or.logger.With(
		slog.String("op", op),
		slog.Any("order", order),
	)

	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()

	orderId, err := or.cli.CreateOrder(ctx, order.ToDtoOrder())
	if err != nil {
		log.Error("failed to create order in orderService", logger.Err(err))
		if st, ok := status.FromError(err); ok {
			log.Error(st.Message(), slog.String("grpc_code", st.Code().String()))

			if st.Code() == codes.InvalidArgument {
				return "", ErrInvalidArguments
			}
		}
		return "", err
	}

	log.Info("order created successfully", slog.String("orderId", orderId))

	return orderId, nil
}

func (or *OrderService) GetOrderById(c context.Context, orderId string) (model.Order, error) {
	const op = "order.GetOrderById"

	log := or.logger.With(
		slog.String("op", op),
		slog.String("orderId", orderId),
	)

	ctx, cancel := context.WithTimeout(c, 30*time.Second)
	defer cancel()

	var res model.Order
	resp, err := or.cli.GetOrderById(ctx, orderId)
	if err != nil {
		log.Error("failed to get order in orderService", logger.Err(err))
		if st, ok := status.FromError(err); ok {
			log.Error(st.Message(), slog.String("grpc_code", st.Code().String()))

			if st.Code() == codes.NotFound {
				return model.Order{}, ErrOrderNotFound
			}
		}
		return model.Order{}, err
	}

	res.FromDtoOrder(resp)

	return res, nil
}
