package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Init(connUrl string) *pgxpool.Pool {
	db, err := pgxpool.New(context.Background(), connUrl)
	if err != nil {
		fmt.Println("error wih db", err)
	}
	return db
}
