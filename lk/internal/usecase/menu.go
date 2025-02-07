package usecase

import (
	"back/lk/internal/entity"
	"back/lk/internal/repo"
	"context"
)

type MenuInterface interface {
	AddFood(ctx context.Context, food entity.Food) (entity.Food, error)
	DeleteFood(ctx context.Context, id uint64) error
	AddCategory(ctx context.Context, category entity.Category) (entity.Category, error)
	DeleteCategory(ctx context.Context, id uint64) error
}

type Menu struct {
	repo repo.MenuInterface
}

func NewMenu(r repo.MenuInterface) MenuInterface {
	return &Menu{repo: r}
}

func (m *Menu) AddFood(ctx context.Context, food entity.Food) (entity.Food, error) {
	res, err := m.repo.AddFood(ctx, food)
	if err != nil {
		return entity.Food{}, err
	}
	return res, nil
}

func (m *Menu) DeleteFood(ctx context.Context, id uint64) error {
	return m.repo.DeleteFood(ctx, id)
}

func (m *Menu) AddCategory(ctx context.Context, category entity.Category) (entity.Category, error) {
	res, err := m.repo.AddCategory(ctx, category)
	if err != nil {
		return entity.Category{}, err
	}
	return res, nil
}

func (m *Menu) DeleteCategory(ctx context.Context, id uint64) error {
	return m.repo.DeleteCategory(ctx, id)
}
