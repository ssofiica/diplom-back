package usecase

import (
	"back/vitrina/internal/repo"
	"context"
)

type RestInterface interface {
	GetInfo(ctx context.Context, id uint64) (int, error)
	GetMenu(ctx context.Context, id uint64) (int, error)
}

type Rest struct {
	repo repo.RestInterface
}

func NewRest(r repo.RestInterface) RestInterface {
	return &Rest{repo: r}
}

func (u *Rest) GetInfo(ctx context.Context, id uint64) (int, error) {
	return 0, nil
}

func (u *Rest) GetMenu(ctx context.Context, id uint64) (int, error) {
	return 0, nil
}
