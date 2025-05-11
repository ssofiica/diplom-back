package repo

import (
	"back/lk/internal/entity"
	"context"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RestInterface interface {
	GetBaseInfo(ctx context.Context, id uint64) (entity.Rest, error)
	GetSchedule(ctx context.Context, id uint64) ([]entity.Schedule, error)
	PutLogoImage(ctx context.Context, url string, id uint64) error
	UploadBaseInfo(ctx context.Context, info entity.BaseInfo, id uint64) (entity.BaseInfo, error)
	UpdateDescrip(ctx context.Context, str string, index uint8, id uint64) error
	PutImgUrl(ctx context.Context, url string, index uint8, id uint64) error
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

func (r *Rest) UploadBaseInfo(ctx context.Context, info entity.BaseInfo, id uint64) (entity.BaseInfo, error) {
	var (
		res entity.BaseInfo
		sb  strings.Builder
	)
	sb.WriteString(
		`update restaurant set`)
	queryFilters, args := r.inQuery(entity.Rest{
		Name:    info.Name,
		Address: info.Address,
		Phone:   info.Phone,
		Email:   info.Email,
	})
	sb.WriteString(queryFilters)
	sb.WriteString(` where id=@id returning id, name, address, phone, email, logo_url`)
	namedArgs := pgx.NamedArgs{
		"name":    args["name"],
		"address": args["address"],
		"phone":   args["phone"],
		"email":   args["email"],
		"id":      id,
	}
	row := r.db.QueryRow(ctx, sb.String(), namedArgs)
	err := row.Scan(&res.Id, &res.Name, &res.Address, &res.Phone, &res.Email, &res.Logo)
	if err != nil {
		return entity.BaseInfo{}, err
	}
	return res, nil
}

func (r *Rest) UpdateDescrip(ctx context.Context, str string, index uint8, id uint64) error {
	var err error
	if str == "" {
		query := `UPDATE restaurant 
		SET description_array = 
			description_array[1:$1-1] || description_array[$1+1:array_length(description_array, 1)]
		WHERE id = $2`
		_, err = r.db.Exec(ctx, query, index, id)
	} else {
		query := `update restaurant set description_array[$1]=$2 where id=$3`
		_, err = r.db.Exec(ctx, query, index, str, id)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *Rest) PutImgUrl(ctx context.Context, url string, index uint8, id uint64) error {
	var err error
	if url == "" {
		query := `UPDATE restaurant 
		SET img_urls = img_urls[1:$1-1] || img_urls[$1+1:array_length(img_urls, 1)]
		WHERE id = $2`
		_, err = r.db.Exec(ctx, query, index, id)
	} else {
		query := `update restaurant set img_urls[$1]=$2 where id=$3`
		_, err = r.db.Exec(ctx, query, index, url, id)
	}
	if err != nil {
		return err
	}
	return nil
}

func (r *Rest) inQuery(params entity.Rest) (string, map[string]any) {
	var (
		args map[string]any = map[string]any{}
		sb   strings.Builder
		arr  []string
	)
	if params.Name != "" {
		arr = append(arr, ` name=@name`)
		args["name"] = params.Name
	}
	if params.Address != "" {
		arr = append(arr, ` address=@address`)
		args["address"] = params.Address
	}
	if params.Phone != "" {
		arr = append(arr, ` phone=@phone`)
		args["phone"] = params.Phone
	}
	if params.Email != "" {
		arr = append(arr, ` email=@email`)
		args["email"] = params.Email
	}
	if params.Logo != "" {
		arr = append(arr, ` logo=@logo`)
		args["logo"] = params.Logo
	}
	if params.Description != nil {
		arr = append(arr, ` description_array=@descrip_array`)
		args["descrip_array"] = params.Description
	}
	if params.Img != nil {
		arr = append(arr, ` img_urls=@img_array`)
		args["img_array"] = params.Img
	}
	sb.WriteString(arr[0])
	for i := 1; i < len(arr); i++ {
		sb.WriteString("," + arr[i])
	}
	return sb.String(), args
}
