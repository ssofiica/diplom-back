package usecase

import (
	"back/lk/internal/entity"
	"back/lk/internal/repo"
	"context"
)

type RestInterface interface {
	GetInfo(ctx context.Context, id uint64) (entity.Rest, error)
}

type Rest struct {
	repo repo.RestInterface
}

func NewRest(r repo.RestInterface) RestInterface {
	return &Rest{repo: r}
}

func (u *Rest) GetInfo(ctx context.Context, id uint64) (entity.Rest, error) {
	res, err := u.repo.GetBaseInfo(ctx, id)
	if err != nil {
		return entity.Rest{}, err
	}
	schedule, err := u.repo.GetSchedule(ctx, id)
	if err != nil {
		return entity.Rest{}, err
	}
	res.Schedule = schedule
	return res, nil
}
