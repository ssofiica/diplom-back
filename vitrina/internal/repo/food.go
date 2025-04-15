package repo

import (
	"back/vitrina/internal/entity"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FoodInterface interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	IsInStock(ctx context.Context, id uint32) (bool, error)
	GetFoodById(ctx context.Context, id uint32) (entity.Food, error)
	AddToOrder(ctx context.Context, tx pgx.Tx, orderId, foodId uint32, count uint8) error
	GetOrderFood(ctx context.Context, id uint32) ([]entity.OrderFood, error)
	GetFoodCountInBasket(ctx context.Context, tx pgx.Tx, foodId uint32, orderId uint32) (uint8, error)
	DeleteFoodFromBasket(ctx context.Context, tx pgx.Tx, foodId uint32, orderId uint32) error
}

type Food struct {
	db *pgxpool.Pool
}

func NewFood(db *pgxpool.Pool) FoodInterface {
	return &Food{db: db}
}

func (r *Food) Begin(ctx context.Context) (pgx.Tx, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *Food) IsInStock(ctx context.Context, id uint32) (bool, error) {
	query := `select status from food where id=$1`
	var res entity.FoodStatus
	err := r.db.QueryRow(ctx, query, id).Scan(&res)
	if err != nil {
		return false, err
	}
	if res != entity.FoodStatusIn {
		return false, nil
	}
	return true, nil
}

func (r *Food) GetFoodById(ctx context.Context, id uint32) (entity.Food, error) {
	query := `select id, name, price, weight, img_url, category_id, restaurant_id from food where id=$1`
	res := entity.Food{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&res.ID,
		&res.Name,
		&res.Price,
		&res.Weight,
		&res.Img,
		&res.CategoryID,
		&res.RestaurantID,
	)
	if err != nil {
		return entity.Food{}, err
	}
	return res, nil
}

func (r *Food) AddToOrder(ctx context.Context, tx pgx.Tx, orderId, foodId uint32, count uint8) error {
	query1 := `insert into order_food (order_id, food_id, count)
				values ($1, $2, $3)
				ON CONFLICT (food_id, order_id) DO UPDATE 
				SET count = EXCLUDED.count`
	_, err := tx.Exec(ctx, query1, orderId, foodId, count)
	if err != nil {
		return err
	}
	return nil
}

func (r *Food) GetOrderFood(ctx context.Context, id uint32) ([]entity.OrderFood, error) {
	query := `select food_id, count from order_food where order_id=$1`
	res := []entity.OrderFood{}
	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return []entity.OrderFood{}, err
	}
	query1 := `select id, name, price, weight, img_url from food where id=$1`
	for rows.Next() {
		tmp := entity.OrderFood{}
		var id uint32
		err := rows.Scan(&id, &tmp.Count)
		if err != nil {
			return []entity.OrderFood{}, err
		}
		err = r.db.QueryRow(ctx, query1, id).Scan(
			&tmp.Food.ID,
			&tmp.Food.Name,
			&tmp.Food.Price,
			&tmp.Food.Weight,
			&tmp.Food.Img,
		)
		if err != nil {
			return []entity.OrderFood{}, err
		}
		res = append(res, tmp)
	}
	return res, nil
}

func (r *Food) GetFoodCountInBasket(ctx context.Context, tx pgx.Tx, foodId uint32, orderId uint32) (uint8, error) {
	query := `select count from order_food where order_id=$1 and food_id=$2`
	var res uint8
	err := tx.QueryRow(ctx, query, orderId, foodId).Scan(&res)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}
	return res, nil
}

func (r *Food) DeleteFoodFromBasket(ctx context.Context, tx pgx.Tx, foodId uint32, orderId uint32) error {
	query := `delete from order_food where order_id=$1 and food_id=$2`
	_, err := tx.Exec(ctx, query, orderId, foodId)
	if err != nil {
		return err
	}
	return nil
}
