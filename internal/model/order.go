package model

import "github.com/tgkzz/gateway/pkg/grpc/order/dto"

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

func (o *Order) FromDtoOrder(req *dto.Order) {
	o.Id = req.Id
	o.Username = req.Username
	o.TotalPrice = req.TotalPrice

	o.Items = make([]Item, len(req.Items))
	for i, item := range req.Items {
		o.Items[i] = Item{
			ItemId:   item.ItemId,
			Name:     item.Name,
			Price:    item.Price,
			Currency: item.Currency,
		}
	}
}

func (o *Order) ToDtoOrder() dto.Order {
	it := make([]dto.Item, len(o.Items))
	for i := range o.Items {
		it[i] = dto.Item{
			ItemId: o.Items[i].ItemId,
			Name:   o.Items[i].Name,
			Price:  o.Items[i].Price,
		}
	}

	return dto.Order{
		Id:         o.Id,
		Username:   o.Username,
		TotalPrice: o.TotalPrice,
		Items:      it,
	}
}
