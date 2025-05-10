package repo

import (
	"context"
	"database/sql"
	"time"

	"back/lk/internal/entity"
)

type AnalyticsInterface interface {
	SetOrder(ctx context.Context, ord entity.Order) error
	GetLinnerCharts(ctx context.Context, restId uint64, start, end string) (entity.LinnerChartRepo, error)
	GetTopFood(ctx context.Context, restId uint64, start, end string) (entity.TopBar, error)
}

type Analytics struct {
	conn *sql.DB
}

func NewAnalytics(c *sql.DB) AnalyticsInterface {
	return &Analytics{
		conn: c,
	}
}

var statuses = map[entity.OrderStatus]int{
	entity.OrderStatusFinished: 1,
	entity.OrderStatusCanceled: 2,
}

var types = map[entity.OrderType]int{
	entity.OrderTypePickup:   1,
	entity.OrderTypeDelivery: 2,
}

func (r *Analytics) GetLinnerCharts(ctx context.Context, restId uint64, start, end string) (entity.LinnerChartRepo, error) {
	query := `
        SELECT
			toDate(created_at) AS date,
            SUM(sum),
            SUM(sum)/COUNT(order_id),
            (SUM(CASE WHEN status = 'finished' THEN 1 ELSE 0 END) / COUNT(order_id)) * 100,
            AVG(ready_at - accepted_at)
        FROM orders
        WHERE restaurant_id = ? AND created_at >= ? AND created_at <= ?
		GROUP BY date
		ORDER BY date;
    `

	var (
		date        []string
		revenue     []int
		avgCheck    []float64
		conversion  []float64
		avgPrepTime []time.Duration
		d           time.Time
		rev         int
		aC          float64
		c           float64
		aPT         time.Duration
	)

	rows, err := r.conn.Query(query, restId, start, end)
	if err != nil {
		return entity.LinnerChartRepo{}, err
	}
	for rows.Next() {
		if err := rows.Scan(&d, &rev, &aC, &c, &aPT); err != nil {
			return entity.LinnerChartRepo{}, err
		}
		date = append(date, d.Format("02-01-2006"))
		revenue = append(revenue, rev)
		avgCheck = append(avgCheck, aC)
		conversion = append(conversion, c)
		avgPrepTime = append(avgPrepTime, aPT)
	}
	return entity.LinnerChartRepo{
		Date:        date,
		Revenue:     revenue,
		AvgCheck:    avgCheck,
		Conversion:  conversion,
		AvgPrepTime: avgPrepTime,
	}, nil
}

func (r *Analytics) GetTopFood(ctx context.Context, restId uint64, start, end string) (entity.TopBar, error) {
	query := `
		SELECT food_name, SUM(count) AS total_orders FROM order_food
        WHERE 
            ordered_at >= ? AND ordered_at <= ?
            AND order_status = 1 AND restaurant_id = ?
        GROUP BY food_name
        ORDER BY total_orders DESC
        LIMIT 5;`
	var (
		name  []string
		count []int
		n     string
		c     int
	)

	rows, err := r.conn.Query(query, start, end, uint32(restId))
	if err != nil {
		return entity.TopBar{}, err
	}
	for rows.Next() {
		err := rows.Scan(&n, &c)
		if err != nil {
			return entity.TopBar{}, err
		}
		name = append(name, n)
		count = append(count, c)
	}
	return entity.TopBar{
		Name:  name,
		Count: count,
	}, nil
}

func (r *Analytics) SetOrder(ctx context.Context, o entity.Order) error {
	tx, err := r.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `INSERT INTO orders (order_id, user_id, status, type, sum, restaurant_id, created_at, accepted_at, ready_at, finished_at, canceled_at) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	_, err = tx.Exec(query, uint64(o.Id), o.User.Id, statuses[o.Status], types[o.Type], o.Sum, o.RestaurantID,
		o.CreatedAt, o.AcceptedAt, o.ReadydAt, o.FinishedAt, o.CanceledAt)
	if err != nil {
		return err
	}

	query1 := `INSERT INTO order_food (order_id, food_id, count, food_name, food_price, food_weight, restaurant_id, category, category_id, ordered_at, order_status) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	stmt, err := tx.Prepare(query1)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, f := range o.Food {
		_, err = stmt.Exec(uint64(o.Id), uint32(f.Food.ID), uint16(f.Count), f.Food.Name,
			f.Food.Price, f.Food.Weight, o.RestaurantID, f.Food.CategoryName, f.Food.CategoryID, o.CreatedAt, statuses[o.Status],
		)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
