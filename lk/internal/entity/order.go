package entity

import (
	"time"
)

type OrderStatus string
type OrderType string

var (
	OrderStatusDraft    OrderStatus = "draft"
	OrderStatusCreated  OrderStatus = "created"
	OrderStatusAccepted OrderStatus = "accepted"
	OrderStatusReady    OrderStatus = "ready"
	OrderStatusFinished OrderStatus = "finished"
	OrderStatusCanceled OrderStatus = "canceled"

	OrderTypeDelivery OrderType = "delivery"
	OrderTypePickup   OrderType = "pickup"
)

type OrderUser struct {
	Id    uint32
	Name  string
	Phone string
}

type UserDTO struct {
	Id    uint32 `json:"id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func (u *OrderUser) ToDTO() UserDTO {
	return UserDTO{
		Id:    u.Id,
		Name:  u.Name,
		Phone: u.Phone,
	}
}

type Order struct {
	Id           uint32
	User         OrderUser
	Status       OrderStatus
	Address      string
	Type         OrderType
	Sum          uint32
	RestaurantID uint32
	Comment      string
	Food         OrderFoodList
	CreatedAt    time.Time
	AcceptedAt   time.Time
	ReadydAt     time.Time
	FinishedAt   time.Time
	CanceledAt   time.Time
}

type OrderDTO struct {
	Id           uint32         `json:"id"`
	User         UserDTO        `json:"user"`
	Status       string         `json:"status"`
	Address      string         `json:"address,omitempty"`
	Sum          uint32         `json:"sum"`
	RestaurantID uint32         `json:"restaurant_id"`
	Comment      string         `json:"comment,omitempty"`
	Type         string         `json:"type,omitempty"`
	Food         []OrderFoodDTO `json:"food"`
	CreatedAt    time.Time      `json:"created_at,omitempty"`
	AcceptedAt   time.Time      `json:"accepted_at,omitempty"`
	ReadydAt     time.Time      `json:"ready_at,omitempty"`
	FinishedAt   time.Time      `json:"finished_at,omitempty"`
	CanceledAt   time.Time      `json:"canceled_at,omitempty"`
}

func (o *Order) ToDTO() OrderDTO {
	f := o.Food.ToDTO()
	return OrderDTO{
		Id:           o.Id,
		User:         o.User.ToDTO(),
		Status:       string(o.Status),
		Type:         string(o.Type),
		Address:      o.Address,
		Sum:          o.Sum,
		RestaurantID: o.RestaurantID,
		Comment:      o.Comment,
		Food:         f,
		CreatedAt:    o.CreatedAt,
		AcceptedAt:   o.AcceptedAt,
		ReadydAt:     o.ReadydAt,
		FinishedAt:   o.FinishedAt,
		CanceledAt:   o.CanceledAt,
	}
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

type MiniOrder struct {
	Id           uint32
	UserID       uint32
	Status       OrderStatus
	Address      string
	Type         OrderType
	Sum          uint32
	RestaurantID uint32
	CreatedAt    time.Time
}

func (o *MiniOrder) ToDTO() MiniOrderDTO {
	return MiniOrderDTO{
		Id:           o.Id,
		UserID:       o.UserID,
		Status:       o.Status,
		Address:      o.Address,
		Type:         o.Type,
		Sum:          o.Sum,
		RestaurantID: o.RestaurantID,
		CreatedAt:    o.CreatedAt,
	}
}

type MiniOrderDTO struct {
	Id           uint32      `json:"id"`
	UserID       uint32      `json:"user_id"`
	Status       OrderStatus `json:"status"`
	Address      string      `json:"address,omitempty"`
	Type         OrderType   `json:"type"`
	Sum          uint32      `json:"sum"`
	RestaurantID uint32      `json:"restaurant_id"`
	CreatedAt    time.Time   `json:"created_at"`
}

type MiniOrderList []MiniOrder

func (o *MiniOrderList) ToDTO() []MiniOrderDTO {
	length := len(*o)
	if length == 0 {
		return []MiniOrderDTO{}
	}
	res := make([]MiniOrderDTO, length)
	for i, tmp := range *o {
		res[i] = tmp.ToDTO()
	}
	return res
}
