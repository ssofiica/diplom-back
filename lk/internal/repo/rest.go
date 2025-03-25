package repo

import (
	"back/lk/internal/entity"
	"context"

	//"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RestInterface interface {
	GetBaseInfo(ctx context.Context, id uint64) (entity.Rest, error)
	GetSchedule(ctx context.Context, id uint64) ([]entity.Schedule, error)
	PutLogoImage(ctx context.Context, url string, id uint64) error
}

type Rest struct {
	db *pgxpool.Pool
}

func NewRest(db *pgxpool.Pool) RestInterface {
	return &Rest{db: db}
}

func (r *Rest) GetBaseInfo(ctx context.Context, id uint64) (entity.Rest, error) {
	query := `select id, name, address, logo_url, description_array, img_urls, phone, email from restaurant where id=$1`
	var res entity.Rest
	err := r.db.QueryRow(ctx, query, id).Scan(
		&res.Id,
		&res.Name,
		&res.Address,
		&res.Logo,
		&res.Description,
		&res.Img,
		&res.Phone,
		&res.Email,
	)
	if err != nil {
		return entity.Rest{}, err
	}
	return res, nil
}

func (r *Rest) GetSchedule(ctx context.Context, id uint64) ([]entity.Schedule, error) {
	query := `select day, open_time, close_time from schedule where restaurant_id=$1`
	res := []entity.Schedule{}
	rows, err := r.db.Query(ctx, query, id)
	if err != nil {
		return res, err
	}
	for rows.Next() {
		var s entity.Schedule
		err := rows.Scan(&s.Day, &s.Open, &s.Close)
		if err != nil {
			return []entity.Schedule{}, err
		}
		res = append(res, s)
	}
	return res, nil
}

func (r *Rest) PutLogoImage(ctx context.Context, url string, id uint64) error {
	query := `update restaurant set logo_url=$1 where id=$2`
	_, err := r.db.Exec(ctx, query, url, id)
	if err != nil {
		return err
	}
	return nil
}
