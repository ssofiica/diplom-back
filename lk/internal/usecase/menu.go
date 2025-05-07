package usecase

import (
	"back/lk/internal/entity"
	"back/lk/internal/repo"
	"context"
	"fmt"
)

var defaultImg = "food/default.jpg"

type MenuInterface interface {
	GetMenu(ctx context.Context, restId uint64) (entity.CategoryList, error)
	GetCategoryList(ctx context.Context, restId uint64) (entity.CategoryList, error)
	GetFoodByStatus(ctx context.Context, status entity.FoodStatus, categoryId uint64) (entity.FoodList, error)
	AddFood(ctx context.Context, food entity.Food) (entity.Food, error)
	DeleteFood(ctx context.Context, id uint64) error
	AddCategory(ctx context.Context, category entity.Category) (entity.Category, error)
	DeleteCategory(ctx context.Context, id uint64) error
	EditFood(ctx context.Context, id uint32, params entity.EditFood) (entity.Food, error)
	ChangeStatus(ctx context.Context, id uint32, status string) error
	UploadFoodLogo(ctx context.Context, file []byte, extention string, mimeType string, foodId uint64, restId uint64) (string, error)
}

type Menu struct {
	repo  repo.MenuInterface
	minio repo.RepoMinio
}

func NewMenu(r repo.MenuInterface, m repo.RepoMinio) MenuInterface {
	return &Menu{repo: r, minio: m}
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

func (m *Menu) GetCategoryList(ctx context.Context, restId uint64) (entity.CategoryList, error) {
	categories, err := m.repo.GetCategories(ctx, restId)
	if err != nil {
		return entity.CategoryList{}, err
	}
	return categories, nil
}

func (m *Menu) GetFoodByStatus(ctx context.Context, status entity.FoodStatus, categoryId uint64) (entity.FoodList, error) {
	return m.repo.GetFoodForCategory(ctx, categoryId, string(status))
}

func (m *Menu) AddFood(ctx context.Context, food entity.Food) (entity.Food, error) {
	food.Img = defaultImg
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

func (m *Menu) EditFood(ctx context.Context, id uint32, params entity.EditFood) (entity.Food, error) {
	return m.repo.EditFood(ctx, id, params)
}

func (m *Menu) ChangeStatus(ctx context.Context, id uint32, status string) error {
	return m.repo.ChangeStatus(ctx, id, status)
}

func (m *Menu) UploadFoodLogo(ctx context.Context, file []byte, extention string, mimeType string, foodId uint64, restId uint64) (string, error) {
	path := fmt.Sprintf("%s/%d/%d%s", ImageTypeFood, restId, foodId, extention)
	_, err := m.minio.UploadImage(ctx, file, path, mimeType)
	if err != nil {
		return "", err
	}
	err = m.repo.UpdateFoodImg(ctx, path, foodId)
	if err != nil {
		return "", err
	}
	return path, nil
}
