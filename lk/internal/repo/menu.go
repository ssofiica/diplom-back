package repo

import (
	"back/lk/internal/entity"
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MenuInterface interface {
	AddFood(ctx context.Context, food entity.Food) (entity.Food, error)
	DeleteFood(ctx context.Context, id uint64) error
	AddCategory(ctx context.Context, category entity.Category) (entity.Category, error)
	DeleteCategory(ctx context.Context, id uint64) error
	GetCategories(ctx context.Context, restId uint64) (entity.CategoryList, error)
	GetFoodForCategory(ctx context.Context, categoryId uint64, status string) (entity.FoodList, error)
	ChangeStatus(ctx context.Context, id uint32, status string) error
	EditFood(ctx context.Context, id uint32, params entity.EditFood) (entity.Food, error)
}

type Menu struct {
	db *pgxpool.Pool
}

func NewMenu(db *pgxpool.Pool) MenuInterface {
	return &Menu{db: db}
}

func (m *Menu) GetCategories(ctx context.Context, restId uint64) (entity.CategoryList, error) {
	query := `select id, name, restaurant_id from category where restaurant_id=$1`
	var res entity.CategoryList
	rows, err := m.db.Query(ctx, query, restId)
	if err != nil {
		return entity.CategoryList{}, err
	}
	for rows.Next() {
		var c entity.Category
		err := rows.Scan(&c.ID, &c.Name, &c.RestaurantID)
		if err != nil {
			return entity.CategoryList{}, err
		}
		res = append(res, c)
	}
	return res, nil
}

func (m *Menu) GetFoodForCategory(ctx context.Context, categoryId uint64, status string) (entity.FoodList, error) {
	query := `select id, name, weight, price, img_url, status from food where category_id=$1 and status=$2;`
	var res entity.FoodList
	rows, err := m.db.Query(ctx, query, categoryId, status)
	if err != nil {
		return entity.FoodList{}, err
	}
	for rows.Next() {
		var f entity.Food
		err := rows.Scan(&f.ID, &f.Name, &f.Weight, &f.Price, &f.Img, &f.Status)
		if err != nil {
			return entity.FoodList{}, err
		}
		res = append(res, f)
	}
	return res, nil
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

func (m *Menu) ChangeStatus(ctx context.Context, id uint32, status string) error {
	query := `update food set status=$1 where id=$2;`
	_, err := m.db.Exec(ctx, query, status, id)
	return err
}

func (m *Menu) EditFood(ctx context.Context, id uint32, params entity.EditFood) (entity.Food, error) {
	var (
		res entity.Food
		sb  strings.Builder
	)
	sb.WriteString(
		`update food set`)
	queryFilters, args := m.inQuery(params)
	sb.WriteString(queryFilters)
	sb.WriteString(` where id=@id returning id, name, price, weight, img_url, status, category_id`)
	namedArgs := pgx.NamedArgs{
		"name":        args["name"],
		"price":       args["price"],
		"weight":      args["weight"],
		"img_url":     args["img_url"],
		"category_id": args["category_id"],
		"status":      args["status"],
		"id":          id,
	}
	row := m.db.QueryRow(ctx, sb.String(), namedArgs)
	err := row.Scan(&res.ID, &res.Name, &res.Price, &res.Weight, &res.Img, &res.Status,
		&res.CategoryID)
	if err != nil {
		return entity.Food{}, err
	}
	return res, nil
}

func (m *Menu) inQuery(params entity.EditFood) (string, map[string]any) {
	var (
		args map[string]any = map[string]any{}
		sb   strings.Builder
		arr  []string
	)
	if params.Name != "" {
		arr = append(arr, ` name=@name`)
		args["name"] = params.Name
	}
	if params.Price > 0 {
		arr = append(arr, ` price=@price`)
		args["price"] = params.Price
	}
	if params.Weight > 0 {
		arr = append(arr, ` weight=@weight`)
		args["weight"] = params.Weight
	}
	if params.Img != "" {
		arr = append(arr, ` img_url=@img_url`)
		args["img_url"] = params.Img
	}
	if params.CategoryID > 0 {
		arr = append(arr, ` category_id=@category_id`)
		args["category_id"] = params.CategoryID
	}
	if params.Status != "" {
		arr = append(arr, ` status=@status`)
		args["status"] = params.Status
	}
	sb.WriteString(arr[0])
	for i := 1; i < len(arr); i++ {
		sb.WriteString("," + arr[i])
	}
	return sb.String(), args
}
