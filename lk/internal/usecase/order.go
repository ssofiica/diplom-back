package usecase

import (
	"context"

	"back/lk/internal/entity"
	"back/lk/internal/repo"
)

type OrderInterface interface {
	GetMiniOrders(ctx context.Context, restId uint32, status string) (entity.MiniOrderList, error)
	GetOrderById(ctx context.Context, id uint32) (entity.Order, error)
	UpdateStatus(ctx context.Context, orderId uint32, status string) error
}

type Order struct {
	repo  repo.OrderInterface
	click repo.ClickHouseInterface
}

func NewOrder(r repo.OrderInterface, c repo.ClickHouseInterface) OrderInterface {
	return &Order{
		repo:  r,
		click: c,
	}
}

func (u *Order) GetMiniOrders(ctx context.Context, restId uint32, status string) (entity.MiniOrderList, error) {
	orders, err := u.repo.GetMiniOrdersByStatus(ctx, restId, status)
	if err != nil {
		return entity.MiniOrderList{}, err
	}
	return orders, nil
}

func (u *Order) GetOrderById(ctx context.Context, id uint32) (entity.Order, error) {
	order, err := u.repo.GetOrderById(ctx, id)
	if err != nil {
		return entity.Order{}, err
	}
	if order.Id == 0 {
		return order, nil
	}
	food, err := u.repo.GetOrderFood(ctx, order.Id)
	if err != nil {
		return entity.Order{}, err
	}
	order.Food = food
	return order, nil
}

func (u *Order) UpdateStatus(ctx context.Context, orderId uint32, status string) error {
	err := u.repo.UpdateStatus(ctx, orderId, entity.OrderStatus(status))
	if err != nil {
		return err
	}
	if status == string(entity.OrderStatusFinished) || status == string(entity.OrderStatusCanceled) {
		order, err := u.GetOrderById(ctx, orderId)
		if err != nil {
			return err
		}
		if err = u.click.SetOrder(ctx, order); err != nil {
			return err
		}
	}
	return nil
}
