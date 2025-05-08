package repo

import (
	"context"
	"database/sql"
	"fmt"

	"back/lk/internal/entity"
)

type ClickHouseInterface interface {
	SetOrder(ctx context.Context, ord entity.Order) error
}

type ClickHouse struct {
	conn *sql.DB
}

func NewClickHouse(c *sql.DB) ClickHouseInterface {
	return &ClickHouse{
		conn: c,
	}
}

var statuses = map[entity.OrderStatus]int{
	entity.OrderStatusFinished: 1,
	entity.OrderStatusCanceled: 2,
}

var types = map[entity.OrderType]int{
	entity.OrderTypePickup:   1,
	entity.OrderTypeDelivery: 2,
}

func (r *ClickHouse) SetOrder(ctx context.Context, o entity.Order) error {
	tx, err := r.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO orders (order_id, user_id, status, type, sum, restaurant_id, created_at, accepted_at, ready_at, finished_at, canceled_at) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	_, err = tx.Exec(query, uint64(o.Id), o.User.Id, statuses[o.Status], types[o.Type], o.Sum, o.RestaurantID,
		o.CreatedAt, o.AcceptedAt, o.ReadydAt, o.FinishedAt, o.CanceledAt)
	if err != nil {
		return err
	}
	fmt.Println("insert orders", err)

	query1 := `INSERT INTO order_food (order_id, food_id, count, food_name, food_price, food_weight, restaurant_id, category, category_id, ordered_at, order_status) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	stmt, err := tx.Prepare(query1)
	fmt.Println("prepare", err, o.Food)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, f := range o.Food {
		_, err = stmt.Exec(uint64(o.Id), uint32(f.Food.ID), uint16(f.Count), f.Food.Name,
			f.Food.Price, f.Food.Weight, o.RestaurantID, "", f.Food.CategoryID, o.CreatedAt, statuses[o.Status],
		)
		fmt.Println("ok", err)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
