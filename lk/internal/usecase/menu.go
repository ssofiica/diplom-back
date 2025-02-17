package usecase

import (
	"back/lk/internal/entity"
	"back/lk/internal/repo"
	"context"
)

type MenuInterface interface {
	GetMenu(ctx context.Context, restId uint64) (entity.CategoryList, error)
	GetFoodByStatus(ctx context.Context, status entity.FoodStatus, categoryId uint64) (entity.FoodList, error)
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

func (m *Menu) GetMenu(ctx context.Context, restId uint64) (entity.CategoryList, error) {
	// получаю список категорий
	categories, err := m.repo.GetCategories(ctx, restId)
	if err != nil {
		return categories, err
	}
	for i, c := range categories {
		// для каждой категории получаю ее блюда
		food, err := m.repo.GetFoodForCategory(ctx, c.ID, "in")
		if err != nil {
			return entity.CategoryList{}, err
		}
		categories[i].Items = food
	}
	return categories, nil
}

func (m *Menu) GetFoodByStatus(ctx context.Context, status entity.FoodStatus, categoryId uint64) (entity.FoodList, error) {
	return m.repo.GetFoodForCategory(ctx, categoryId, string(status))
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
