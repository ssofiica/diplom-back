package repo

import (
	"context"
	"errors"

	"back/lk/internal/entity"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderInterface interface {
	GetMiniOrdersByStatus(ctx context.Context, restId uint32, status string) (entity.MiniOrderList, error)
	GetOrderById(ctx context.Context, orderId uint32) (entity.Order, error)
	GetOrderFood(ctx context.Context, id uint32) ([]entity.OrderFood, error)
	UpdateStatus(ctx context.Context, id uint32, status entity.OrderStatus) error
}

type Order struct {
	db    *pgxpool.Pool
	click *driver.Conn
}

func NewOrder(db *pgxpool.Pool, click *driver.Conn) OrderInterface {
	return &Order{
		db:    db,
		click: click,
	}
}

func (r *Order) GetMiniOrdersByStatus(ctx context.Context, restId uint32, status string) (entity.MiniOrderList, error) {
	query := `select id, status, type, sum, created_at from "order" where restaurant_id=$1 and status=$2 order by created_at DESC`
	var res entity.MiniOrderList
	rows, err := r.db.Query(ctx, query, restId, status)
	if err != nil {
		return entity.MiniOrderList{}, nil
	}
	for rows.Next() {
		ord := entity.MiniOrder{}
		var address pgtype.Text
		err = rows.Scan(
			&ord.Id,
			&ord.Status,
			&ord.Type,
			&ord.Sum,
			&ord.CreatedAt,
		)
		if err != nil {
			return entity.MiniOrderList{}, nil
		}
		ord.Address = address.String
		res = append(res, ord)
	}
	return res, nil
}

func (r *Order) GetOrderById(ctx context.Context, orderId uint32) (entity.Order, error) {
	query := `select id, user_id, status, address, type, sum, comment, restaurant_id,
				created_at, accepted_at, ready_at, finished_at, canceled_at
			 	from "order" where id=$1`
	var res entity.Order
	var address, Type, comment pgtype.Text
	var acceptedAt, readyAt, finishedAt, canceledAt pgtype.Timestamptz
	var userID uint32
	err := r.db.QueryRow(ctx, query, orderId).Scan(
		&res.Id,
		&userID,
		&res.Status,
		&address,
		&Type,
		&res.Sum,
		&comment,
		&res.RestaurantID,
		&res.CreatedAt,
		&acceptedAt,
		&readyAt,
		&finishedAt,
		&canceledAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Order{Id: 0}, nil
		}
		return entity.Order{}, err
	}
	res.Address = address.String
	res.Type = entity.OrderType(Type.String)
	res.Comment = comment.String
	res.AcceptedAt = acceptedAt.Time
	res.ReadydAt = readyAt.Time
	res.FinishedAt = finishedAt.Time
	res.CanceledAt = canceledAt.Time

	query = `select id, name, phone from "user" where id=$1`
	user := entity.OrderUser{}
	err = r.db.QueryRow(ctx, query, userID).Scan(&user.Id, &user.Name, &user.Phone)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Order{Id: 0}, nil
		}
		return entity.Order{}, err
	}
	res.User = user
	return res, nil
}

func (r *Order) GetOrderFood(ctx context.Context, id uint32) ([]entity.OrderFood, error) {
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

func (r *Order) UpdateStatus(ctx context.Context, id uint32, status entity.OrderStatus) error {
	query := `update "order" set status=$1`
	if status == entity.OrderStatusAccepted {
		query += `, accepted_at=now()`
	} else if status == entity.OrderStatusReady {
		query += `, ready_at=now()`
	} else if status == entity.OrderStatusFinished {
		query += `, finished_at=now()`
	} else if status == entity.OrderStatusCanceled {
		query += `, canceled_at=now()`
	} else {
		return errors.New("wrong status")
	}
	query += ` where id=$2`
	_, err := r.db.Exec(ctx, query, status, id)
	if err != nil {
		return err
	}
	return nil
}
