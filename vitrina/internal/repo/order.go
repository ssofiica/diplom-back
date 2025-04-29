package repo

import (
	"back/vitrina/internal/entity"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderInterface interface {
	GetUserBasket(ctx context.Context, userId, id uint32) (entity.Order, error)
	CreateBasket(ctx context.Context, userId uint32, restId uint32) (uint32, error)
	UpdateBasketSum(ctx context.Context, tx pgx.Tx, id uint32, value uint16, plus bool) error
	GetOrderById(ctx context.Context, orderId uint32) (entity.Order, error)
	UpdateBasketInfo(ctx context.Context, id uint32, data entity.RequestBasketInfo) (entity.Order, error)
	UpdateStatus(ctx context.Context, id uint32, status entity.OrderStatus) error
	GetMiniOrdersByStatus(ctx context.Context, userId uint32, statuses []entity.OrderStatus) (entity.MiniOrderList, error)
}

type Order struct {
	db *pgxpool.Pool
}

func NewOrder(db *pgxpool.Pool) OrderInterface {
	return &Order{db: db}
}

func (r *Order) GetUserBasket(ctx context.Context, userId, orderId uint32) (entity.Order, error) {
	query := `select id, user_id, created_at, status, type, address, comment, sum, restaurant_id
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
	var address, comment pgtype.Text
	err := r.db.QueryRow(ctx, query, id).Scan(
		&res.Id,
		&res.UserID,
		&res.CreatedAt,
		&res.Status,
		&res.Type,
		&address,
		&comment,
		&res.Sum,
		&res.RestaurantID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Order{Id: 0}, nil
		}
		return entity.Order{}, err
	}
	res.Address = address.String
	res.Comment = comment.String
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

func (r *Order) UpdateBasketSum(ctx context.Context, tx pgx.Tx, id uint32, value uint16, plus bool) error {
	query := `update "order" set sum = sum`
	if plus {
		query = query + `+ $1 where id=$2`
	} else {
		query = query + `- $1 where id=$2`
	}
	_, err := tx.Exec(ctx, query, value, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Order) GetOrderById(ctx context.Context, orderId uint32) (entity.Order, error) {
	query := `select id, user_id, status, address, type, sum, comment, restaurant_id,
				created_at, accepted_at, ready_at, finished_at, canceled_at
			 	from "order" where id=$1`
	var res entity.Order
	var address, Type, comment pgtype.Text
	var acceptedAt, readyAt, finishedAt, canceledAt pgtype.Timestamptz
	err := r.db.QueryRow(ctx, query, orderId).Scan(
		&res.Id,
		&res.UserID,
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
	return res, nil
}

func (r *Order) UpdateBasketInfo(ctx context.Context, id uint32, data entity.RequestBasketInfo) (entity.Order, error) {
	var (
		res entity.Order
		sb  strings.Builder
	)
	sb.WriteString(
		`update "order" set`)
	queryFilters, args := r.inQuery(data)
	sb.WriteString(queryFilters)
	sb.WriteString(` where id=@id returning id, address, type, comment`)
	namedArgs := pgx.NamedArgs{
		"address": args["address"],
		"type":    args["type"],
		"comment": args["comment"],
		"id":      id,
	}
	var address, Type, comment pgtype.Text
	row := r.db.QueryRow(ctx, sb.String(), namedArgs)
	err := row.Scan(&res.Id, &address, &Type, &comment)
	if err != nil {
		return entity.Order{}, err
	}
	res.Address = address.String
	res.Type = entity.OrderType(Type.String)
	res.Comment = comment.String
	return res, err
}

func (r *Order) inQuery(params entity.RequestBasketInfo) (string, map[string]any) {
	var (
		args map[string]any = map[string]any{}
		sb   strings.Builder
		arr  []string
	)
	if params.Address != "" {
		arr = append(arr, ` address=@address`)
		args["address"] = params.Address
	}
	if params.Comment != "" {
		arr = append(arr, ` comment=@comment`)
		args["comment"] = params.Comment
	}
	if params.Type != "" {
		arr = append(arr, ` type=@type`)
		args["type"] = params.Type
	}
	sb.WriteString(arr[0])
	if len(arr) == 1 {
		return sb.String(), args
	}
	for i := 1; i < len(arr); i++ {
		sb.WriteString("," + arr[i])
	}
	return sb.String(), args
}

func (r *Order) UpdateStatus(ctx context.Context, id uint32, status entity.OrderStatus) error {
	if status != entity.OrderStatusCreated && status != entity.OrderStatusCanceled {
		return errors.New("Неверный статус")
	}
	query := `update "order" set status=$1`
	if status == entity.OrderStatusCreated {
		query += `, created_at=now()`
	} else if status == entity.OrderStatusCanceled {
		query += `, canceled_at=now()`
	}
	query += ` where id=$2`
	_, err := r.db.Exec(ctx, query, status, id)
	if err != nil {
		return err
	}
	return nil
}

func (r *Order) GetMiniOrdersByStatus(ctx context.Context, userId uint32, statuses []entity.OrderStatus) (entity.MiniOrderList, error) {
	if len(statuses) == 0 {
		return nil, errors.New("Нет статусов")
	}
	placeholders := make([]string, len(statuses))
	for i, st := range statuses {
		placeholders[i] = fmt.Sprintf("'%s'", st)
	}

	query := fmt.Sprintf(`select id, user_id, status, address, type, sum, restaurant_id,
				created_at from "order" where user_id=$1 and status IN (%s) order by created_at DESC`, strings.Join(placeholders, ", "))
	var res entity.MiniOrderList
	rows, err := r.db.Query(ctx, query, userId)
	if err != nil {
		return entity.MiniOrderList{}, nil
	}
	for rows.Next() {
		ord := entity.MiniOrder{}
		var address pgtype.Text
		err = rows.Scan(
			&ord.Id,
			&ord.UserID,
			&ord.Status,
			&address,
			&ord.Type,
			&ord.Sum,
			&ord.RestaurantID,
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
