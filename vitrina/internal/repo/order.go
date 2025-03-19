package repo

import (
	"back/vitrina/internal/entity"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderInterface interface {
	GetUserBasket(ctx context.Context, userId, id uint32) (entity.Order, error)
	CreateBasket(ctx context.Context, userId uint32, restId uint32) (uint32, error)
	UpdateBasketSum(ctx context.Context, tx pgx.Tx, id uint32, value uint16) error
}

type Order struct {
	db *pgxpool.Pool
}

func NewOrder(db *pgxpool.Pool) OrderInterface {
	return &Order{db: db}
}

func (r *Order) GetUserBasket(ctx context.Context, userId, orderId uint32) (entity.Order, error) {
	query := `select id, user_id, created_at, status, sum, restaurant_id
			 	from "order" where`
	var id uint32
	if userId > 0 {
		query += ` user_id=$1 and status='draft'`
		id = userId
	} else if id > 0 {
		query += ` id=$1 and status='draft'`
		id = orderId
	}
	var res entity.Order
	err := r.db.QueryRow(ctx, query, id).Scan(
		&res.Id,
		&res.UserID,
		&res.CreatedAt,
		&res.Status,
		&res.Sum,
		&res.RestaurantID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Order{Id: 0}, nil
		}
		return entity.Order{}, err
	}
	return res, nil
}

func (r *Order) CreateBasket(ctx context.Context, userId uint32, restId uint32) (uint32, error) {
	query := `
			insert into "order" (user_id, status, sum, restaurant_id, created_at)
			values ($1, $2, 0, $3, now())
			returning id
	`
	var res uint32
	err := r.db.QueryRow(ctx, query, userId, entity.OrderStatusDraft, restId).Scan(&res)
	if err != nil {
		return 0, err
	}
	return res, nil
}

func (r *Order) UpdateBasketSum(ctx context.Context, tx pgx.Tx, id uint32, value uint16) error {
	query := `update "order" set sum = sum + $1 where id=$2`
	_, err := tx.Exec(ctx, query, value, id)
	if err != nil {
		return err
	}
	return nil
}
