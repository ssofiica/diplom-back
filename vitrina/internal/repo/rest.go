package repo

import (
	"context"

	//"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type RestInterface interface {
	GetInfo(ctx context.Context, id uint64) (int, error)
	GetMenu(ctx context.Context, id uint64) (int, error)
}

type Rest struct {
	db *pgxpool.Pool
}

func NewRest(db *pgxpool.Pool) RestInterface {
	return &Rest{db: db}
}

func (r *Rest) GetInfo(ctx context.Context, id uint64) (int, error) {
	return 0, nil
}

func (r *Rest) GetMenu(ctx context.Context, id uint64) (int, error) {
	return 0, nil
}
