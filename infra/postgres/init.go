package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func Init(connUrl string) *pgxpool.Pool {
	db, err := pgxpool.New(context.Background(), connUrl)
	if err != nil {
		fmt.Println("error wih db", err)
	}
	if err := db.Ping(context.Background()); err != nil {
		log.Fatal("Failed to connect to PostgreSQL", err)
		return nil
	}
	return db
}
