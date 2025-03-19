package entity

import (
	"time"
)

type OrderStatus string

var (
	OrderStatusDraft    OrderStatus = "draft"
	OrderStatusCreated  OrderStatus = "created"
	OrderStatusAccepted OrderStatus = "accepted"
	OrderStatusReady    OrderStatus = "ready"
	OrderStatusFinished OrderStatus = "finished"
	OrderStatusCanceled OrderStatus = "canceled"
)

type Order struct {
	Id           uint32
	UserID       uint32
	Status       OrderStatus
	Address      string
	Sum          uint32
	RestaurantID uint32
	Food         OrderFoodList
	CreatedAt    time.Time
	AcceptedAt   time.Time
	ReadydAt     time.Time
	FinishedAt   time.Time
	CanceledAt   time.Time
}

type OrderDTO struct {
	Id           uint32         `json:"id"`
	UserID       uint32         `json:"user_id"`
	CreatedAt    time.Time      `json:"created_at"`
	Status       string         `json:"status"`
	Address      string         `json:"address"`
	Sum          uint32         `json:"sum"`
	RestaurantID uint32         `json:"restaurant_id"`
	Food         []OrderFoodDTO `json:"food"`
}

func (o *Order) ToDTO() OrderDTO {
	f := o.Food.ToDTO()
	return OrderDTO{
		Id:           o.Id,
		UserID:       o.UserID,
		CreatedAt:    o.CreatedAt,
		Status:       string(o.Status),
		Address:      o.Address,
		Sum:          o.Sum,
		RestaurantID: o.RestaurantID,
		Food:         f,
	}
}

type RequestAddFood struct {
	FoodId uint32 `json:"food_id"`
	Count  uint8  `json:"count"`
}

func (r *RequestAddFood) Valid() bool {
	return r.FoodId > 0 && r.Count > 0
}

type OrderFood struct {
	Food  Food
	Count uint8
}

type OrderFoodList []OrderFood

type OrderFoodDTO struct {
	Food  FoodDTO `json:"item"`
	Count uint8   `json:"count"`
}

func (o *OrderFood) ToDTO() OrderFoodDTO {
	return OrderFoodDTO{
		Food:  o.Food.ToDTO(),
		Count: o.Count,
	}
}

func (o *OrderFoodList) ToDTO() []OrderFoodDTO {
	length := len(*o)
	if length == 0 {
		return []OrderFoodDTO{}
	}
	res := make([]OrderFoodDTO, length)
	for i, tmp := range *o {
		res[i] = tmp.ToDTO()
	}
	return res
}
