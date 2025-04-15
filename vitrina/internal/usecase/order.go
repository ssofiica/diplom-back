package usecase

import (
	"back/vitrina/internal/entity"
	"back/vitrina/internal/repo"
	"errors"

	"context"
)

var ErrFoodStoped = errors.New("Блюда нет в наличии")

type OrderInterface interface {
	AddFoodToOrder(ctx context.Context, userId, restId uint32, data entity.RequestAddFood) error
	GetBasketId(ctx context.Context, userId uint32) (uint32, error)
	GetBasket(ctx context.Context, userId, id uint32) (entity.Order, error)
	GetOrderById(ctx context.Context, id uint32) (entity.Order, error)
	UpdateBasketInfo(ctx context.Context, userId uint32, info entity.RequestBasketInfo) (entity.Order, error)
}

type Order struct {
	repoOrder repo.OrderInterface
	repoFood  repo.FoodInterface
}

func NewOrder(o repo.OrderInterface, f repo.FoodInterface) OrderInterface {
	return &Order{repoOrder: o, repoFood: f}
}

func (u *Order) AddFoodToOrder(ctx context.Context, userId, restId uint32, data entity.RequestAddFood) error {
	// получаем корзину
	order, err := u.repoOrder.GetUserBasket(ctx, userId, 0)
	if err != nil {
		return err
	}
	var id uint32
	if order.Id == 0 {
		// если у чела нет корзины, создаем ее
		id, err = u.repoOrder.CreateBasket(ctx, userId, restId)
		if err != nil {
			return err
		}
	} else {
		id = order.Id
	}
	// проверяем, что блюдо в наличии
	is, err := u.repoFood.IsInStock(ctx, data.FoodId)
	if err != nil {
		return err
	}
	if !is {
		return ErrFoodStoped
	}
	// берем инфу о блюде, из нее цену
	food, err := u.repoFood.GetFoodById(ctx, data.FoodId)
	if err != nil {
		return err
	}
	// начинаем транзакцию
	tx, err := u.repoFood.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	countNow, err := u.repoFood.GetFoodCountInBasket(ctx, tx, data.FoodId, id)
	if err != nil {
		return err
	}
	if data.Count == 0 {
		err = u.repoFood.DeleteFoodFromBasket(ctx, tx, data.FoodId, id)
		if err != nil {
			return err
		}
		err = u.repoOrder.UpdateBasketSum(ctx, tx, id, uint16(countNow)*food.Price, false)
		if err != nil {
			return err
		}
		err = tx.Commit(ctx)
		if err != nil {
			return err
		}
		return nil
	}
	// добавляем блюдо
	err = u.repoFood.AddToOrder(ctx, tx, id, data.FoodId, data.Count)
	if err != nil {
		return err
	}
	if num := int16(data.Count) - int16(countNow); num > 0 {
		err = u.repoOrder.UpdateBasketSum(ctx, tx, id, uint16(num)*food.Price, true)
	} else {
		dif := countNow - data.Count
		err = u.repoOrder.UpdateBasketSum(ctx, tx, id, uint16(dif)*food.Price, false)
	}
	if err != nil {
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (u *Order) GetBasketId(ctx context.Context, userId uint32) (uint32, error) {
	order, err := u.repoOrder.GetUserBasket(ctx, userId, 0)
	if err != nil {
		return 0, err
	}
	return order.Id, nil
}

func (u *Order) GetBasket(ctx context.Context, userId, id uint32) (entity.Order, error) {
	order, err := u.repoOrder.GetUserBasket(ctx, userId, id)
	if err != nil {
		return entity.Order{}, err
	}
	if order.Id == 0 {
		return order, nil
	}
	food, err := u.repoFood.GetOrderFood(ctx, order.Id)
	if err != nil {
		return entity.Order{}, err
	}
	order.Food = food
	return order, nil
}

func (u *Order) GetOrderById(ctx context.Context, id uint32) (entity.Order, error) {
	order, err := u.repoOrder.GetOrderById(ctx, id)
	if err != nil {
		return entity.Order{}, err
	}
	if order.Id == 0 {
		return order, nil
	}
	food, err := u.repoFood.GetOrderFood(ctx, order.Id)
	if err != nil {
		return entity.Order{}, err
	}
	order.Food = food
	return order, nil
}

func (u *Order) UpdateBasketInfo(ctx context.Context, userId uint32, info entity.RequestBasketInfo) (entity.Order, error) {
	order, err := u.repoOrder.GetUserBasket(ctx, userId, 0)
	if err != nil {
		return entity.Order{}, err
	}
	if order.Id == 0 {
		return order, nil
	}
	order, err = u.repoOrder.UpdateBasketInfo(ctx, order.Id, info)
	if err != nil {
		return entity.Order{}, err
	}
	return order, nil
}
