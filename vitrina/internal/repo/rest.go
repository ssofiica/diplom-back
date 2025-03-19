package repo

import (
	"back/vitrina/internal/entity"
	"context"

	//"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RestInterface interface {
	GetBaseInfo(ctx context.Context, id uint64) (entity.Rest, error)
	GetSchedule(ctx context.Context, id uint64) ([]entity.Schedule, error)
	GetCategories(ctx context.Context, restId uint64) (entity.CategoryList, error)
	GetFoodForCategory(ctx context.Context, categoryId uint64, status string) (entity.FoodList, error)
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

func (r *Rest) GetCategories(ctx context.Context, restId uint64) (entity.CategoryList, error) {
	query := `select id, name from category where restaurant_id=$1`
	var res entity.CategoryList
	rows, err := r.db.Query(ctx, query, restId)
	if err != nil {
		return entity.CategoryList{}, err
	}
	for rows.Next() {
		var c entity.Category
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return entity.CategoryList{}, err
		}
		res = append(res, c)
	}
	return res, nil
}

func (r *Rest) GetFoodForCategory(ctx context.Context, categoryId uint64, status string) (entity.FoodList, error) {
	query := `select id, name, weight, price, img_url from food where category_id=$1 and status=$2;`
	var res entity.FoodList
	rows, err := r.db.Query(ctx, query, categoryId, status)
	if err != nil {
		return entity.FoodList{}, err
	}
	for rows.Next() {
		var f entity.Food
		err := rows.Scan(&f.ID, &f.Name, &f.Weight, &f.Price, &f.Img)
		if err != nil {
			return entity.FoodList{}, err
		}
		res = append(res, f)
	}
	return res, nil
}
