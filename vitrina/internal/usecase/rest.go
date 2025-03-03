package usecase

import (
	"back/vitrina/entity"
	"back/vitrina/internal/repo"
	"context"
)

type RestInterface interface {
	GetInfo(ctx context.Context, id uint64) (entity.Rest, error)
	GetMenu(ctx context.Context, id uint64) (entity.CategoryList, error)
}

type Rest struct {
	repo repo.RestInterface
}

func NewRest(r repo.RestInterface) RestInterface {
	return &Rest{repo: r}
}

func (u *Rest) GetInfo(ctx context.Context, id uint64) (entity.Rest, error) {
	res, err := u.repo.GetBaseInfo(ctx, id)
	if err != nil {
		return entity.Rest{}, err
	}
	schedule, err := u.repo.GetSchedule(ctx, id)
	if err != nil {
		return entity.Rest{}, err
	}
	res.Schedule = schedule
	return res, nil
}

func (u *Rest) GetMenu(ctx context.Context, id uint64) (entity.CategoryList, error) {
	// получаю список категорий
	categories, err := u.repo.GetCategories(ctx, id)
	if err != nil {
		return categories, err
	}
	for i, c := range categories {
		// для каждой категории получаю ее блюда
		food, err := u.repo.GetFoodForCategory(ctx, c.ID, "in")
		if err != nil {
			return entity.CategoryList{}, err
		}
		categories[i].Items = food
	}
	return categories, nil
}
