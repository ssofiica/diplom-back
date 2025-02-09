package repo

import (
	"back/lk/internal/entity"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MenuInterface interface {
	AddFood(ctx context.Context, food entity.Food) (entity.Food, error)
	DeleteFood(ctx context.Context, id uint64) error
	AddCategory(ctx context.Context, category entity.Category) (entity.Category, error)
	DeleteCategory(ctx context.Context, id uint64) error
}

type Menu struct {
	db *pgxpool.Pool
}

func NewMenu(db *pgxpool.Pool) MenuInterface {
	return &Menu{db: db}
}

func (m *Menu) AddFood(ctx context.Context, food entity.Food) (entity.Food, error) {
	query := `
			insert into food(name, restaurant_id, category_id, weight, price, img_url, status)
			values ($1, $2, $3, $4, $5, $6, $7)
			returning id, name, restaurant_id, category_id, weight, price, img_url, status
	`
	var res entity.Food
	row := m.db.QueryRow(ctx, query,
		food.Name, food.RestaurantID, food.CategoryID, food.Weight, food.Price, food.Img, food.Status)
	err := row.Scan(&res.ID, &res.Name, &res.RestaurantID, &res.CategoryID, &res.Weight,
		&res.Price, &res.Img, &res.Status)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (m *Menu) AddCategory(ctx context.Context, category entity.Category) (entity.Category, error) {
	query := `
			insert into category(name, restaurant_id)
			values ($1, $2)
			returning id, name, restaurant_id
	`
	var res entity.Category
	row := m.db.QueryRow(ctx, query, category.Name, category.RestaurantID)
	err := row.Scan(&res.ID, &res.Name, &res.RestaurantID)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (m *Menu) DeleteFood(ctx context.Context, id uint64) error {
	query := `update food set status='delete' where id=$1;`
	_, err := m.db.Exec(ctx, query, id)
	return err
}

func (m *Menu) DeleteCategory(ctx context.Context, id uint64) error {
	query := `delete from category where id=$1;`
	_, err := m.db.Exec(ctx, query, id)
	return err
}
