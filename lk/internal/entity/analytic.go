package entity

import "time"

//    Графики
// Выручка	SUM(sum)	Общая/дневная/по ресторану
// Средний чек	SUM(sum)/COUNT(order_id)	Тренды по дням/ресторанам
// Конверсия	(finished orders)/(all orders)*100%	Эффективность работы
// Время приготовления	AVG(ready_at - accepted_at)	Оптимизация кухни
//
//    Показатели
// Топ-5 блюд

// currentTime := time.Now()
// onlyDate := currentTime.Format("02-01-2006")

type LinnerChart struct {
	Title string
	X     string
	Y     string
	XArr  any
	YArr  any
}

type LinnerChartRepo struct {
	Date        []string
	Revenue     []int
	AvgCheck    []float64
	Conversion  []float64
	AvgPrepTime []time.Duration
}

type Analytics struct {
	Revenue     LinnerChart `json:"revenue"`
	AvgCheck    LinnerChart `json:"avg_check"`
	Conversion  LinnerChart `json:"conversion"`
	AvgPrepTime LinnerChart `json:"avg_prep_time"`
	TopFood     TopBar      `json:"top_food"`
}

type TopBar struct {
	Title string
	Name  []string
	Count []int
}
