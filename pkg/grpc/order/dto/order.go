package dto

import order1 "github.com/tgkzz/order/gen/go/order"

type Order struct {
	Id         string  `json:"id"`
	Username   string  `json:"username"`
	TotalPrice float64 `json:"total_price"`
	Items      []Item  `json:"items"`
}

type Item struct {
	ItemId   string  `json:"item_id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
}

func (o Order) FromDtoToCliOrder() *order1.CreateOrderRequest {
	items := make([]*order1.CreateOrderItemRequest, len(o.Items))
	for i, item := range o.Items {
		items[i] = &order1.CreateOrderItemRequest{
			Name:     item.Name,
			Price:    float32(item.Price),
			Currency: item.Currency,
		}
	}

	return &order1.CreateOrderRequest{
		Username:   o.Username,
		TotalPrice: float32(o.TotalPrice),
		Items:      items,
	}
}

func (o Order) FromCliOrderToDto(req *order1.GetOrderResponse) *Order {
	items := make([]Item, len(o.Items))
	for i, item := range req.Items {
		items[i] = Item{
			ItemId:   item.ItemId,
			Name:     item.Name,
			Price:    float64(item.Price),
			Currency: item.Currency,
		}
	}

	return &Order{
		Id:         req.GetOrderId(),
		Username:   req.GetUsername(),
		TotalPrice: float64(req.GetPrice()),
		Items:      items,
	}
}
