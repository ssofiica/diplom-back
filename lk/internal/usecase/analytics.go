package usecase

import (
	"context"
	"time"

	"back/lk/internal/entity"
	"back/lk/internal/repo"
)

type AnalyticsInterface interface {
	GetAnalytics(ctx context.Context, restId uint64, start, end time.Time) (entity.Analytics, error)
}

type Analytics struct {
	click repo.AnalyticsInterface
	menu  repo.MenuInterface
}

func NewAnalytics(c repo.AnalyticsInterface, m repo.MenuInterface) AnalyticsInterface {
	return &Analytics{click: c, menu: m}
}

func (u *Analytics) GetAnalytics(ctx context.Context, restId uint64, start, end time.Time) (entity.Analytics, error) {
	linner, err := u.click.GetLinnerCharts(ctx, restId, start, end)
	if err != nil {
		return entity.Analytics{}, err
	}
	res := entity.Analytics{
		Revenue: entity.LinnerChart{
			Title: "Дневная выручка",
			X:     "Выручка",
			Y:     "руб",
			XArr:  linner.Date,
			YArr:  linner.Revenue,
		},
		AvgCheck: entity.LinnerChart{
			Title: "Средний чек",
			X:     "Чек",
			Y:     "руб.",
			XArr:  linner.Date,
			YArr:  linner.AvgCheck,
		},
		Conversion: entity.LinnerChart{
			Title: "Процент завершенных заказов",
			X:     "Процент",
			Y:     "%",
			XArr:  linner.Date,
			YArr:  linner.Conversion,
		},
		AvgPrepTime: entity.LinnerChart{
			Title: "Время приготовления",
			X:     "Время",
			Y:     "мин.",
			XArr:  linner.Date,
			YArr:  linner.AvgPrepTime,
		},
	}
	food, err := u.click.GetTopFood(ctx, restId, start, end)
	if err != nil {
		return entity.Analytics{}, err
	}
	food.Title = "Топ-5 блюд по кол-ву заказов"
	res.TopFood = food
	return res, nil
}
