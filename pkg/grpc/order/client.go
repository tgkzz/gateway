package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/tgkzz/gateway/pkg/grpc/order/dto"
	order1 "github.com/tgkzz/order/gen/go/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log/slog"
)

type OrderClient interface {
	CreateOrder(ctx context.Context, order dto.Order) (string, error)
	GetOrderById(ctx context.Context, orderId string) (*dto.Order, error)
	DeleteOrderById(ctx context.Context, orderId string) error
}

type Client struct {
	client order1.OrderServiceClient
	logger *slog.Logger
}

var (
	ErrEmptyOrderId        = errors.New("order id is empty")
	ErrCouldNotCreateOrder = errors.New("could not create order, bad arguments")
)

func NewOrderClient(host, port string, logger *slog.Logger) (OrderClient, error) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", host, port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	c := order1.NewOrderServiceClient(conn)

	return &Client{client: c, logger: logger}, err
}

func (c *Client) CreateOrder(ctx context.Context, order dto.Order) (string, error) {
	const op = "grpcOrder.CreateOrder"

	log := c.logger.With(
		slog.String("op", op),
		slog.Any("order", order),
	)

	resp, err := c.client.CreateOrder(ctx, order.FromDtoToCliOrder())
	if err != nil {
		log.Error(err.Error())
		if st, ok := status.FromError(err); ok {
			log.Error(st.Message(), slog.String("grpc_code", st.Code().String()))

			if st.Code() == codes.InvalidArgument {
				return "", ErrCouldNotCreateOrder
			}
		}
		return "", err
	}

	if resp.OrderId == "" {
		log.Error("empty order id received")
		return "", ErrEmptyOrderId
	}

	log.Info("order created successfully in gRPC client", slog.String("orderId", resp.GetOrderId()))

	return resp.GetOrderId(), err
}

func (c *Client) GetOrderById(ctx context.Context, orderId string) (*dto.Order, error) {
	const op = "grpcOrder.GetOrderById"

	log := c.logger.With(
		slog.String("op", op),
		slog.Any("orderId", orderId),
	)

	var res dto.Order
	resp, err := c.client.GetOrderById(ctx, &order1.GetOrderRequest{OrderId: orderId})
	if err != nil {
		log.Error(err.Error())
		// TODO: add error handling of codes in both services
		return nil, err
	}

	r := res.FromCliOrderToDto(resp)

	return r, nil
}

func (c *Client) DeleteOrderById(ctx context.Context, orderId string) error {
	const op = "grpcOrder.DeleteOrderById"

	log := c.logger.With(
		slog.String("op", op),
		slog.Any("orderId", orderId),
	)

	if _, err := c.client.DeleteOrderById(
		ctx,
		&order1.DeleteOrderRequest{
			OrderId: orderId,
		},
	); err != nil {
		log.Error(err.Error())
		return err
	}

	return nil
}
