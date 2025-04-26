package main

import (
	"back/infra/postgres"
	"back/vitrina/config"
	"back/vitrina/internal/delivery"
	"back/vitrina/internal/middleware"
	"back/vitrina/internal/repo"
	"back/vitrina/internal/usecase"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err)
	}
}

func main() {
	cfg := config.Load()
	db := postgres.Init(os.Getenv("PG_CONN"))
	defer db.Close()

	infoRepo := repo.NewRest(db)
	infoUsecase := usecase.NewRest(infoRepo)
	infoHandler := delivery.NewRestHandler(infoUsecase)

	foodRepo := repo.NewFood(db)
	orderRepo := repo.NewOrder(db)
	orderUsecase := usecase.NewOrder(orderRepo, foodRepo)
	orderHandler := delivery.NewOrderHandler(orderUsecase)

	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	r.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("query to path: " + r.URL.String())
		w.WriteHeader(http.StatusOK)
	})
	r.Use(middleware.CorsMiddleware)

	r.HandleFunc("/info", infoHandler.GetInfo).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/menu", infoHandler.GetMenu).Methods(http.MethodGet, http.MethodOptions)
	order := r.PathPrefix("/order").Subrouter()
	{
		order.HandleFunc("/food", orderHandler.ChangeFoodCountInBasket).Methods(http.MethodPost, http.MethodOptions)
		order.HandleFunc("/basket", orderHandler.GetUserBasket).Methods(http.MethodGet, http.MethodOptions)
		order.HandleFunc("/info", orderHandler.UpdateBasketInfo).Methods(http.MethodPost, http.MethodOptions)
		order.HandleFunc("/pay", orderHandler.Pay).Methods(http.MethodPost, http.MethodOptions)
		order.HandleFunc("/current", orderHandler.GetCurrent).Methods(http.MethodGet, http.MethodOptions)
		order.HandleFunc("/archive", orderHandler.GetArchive).Methods(http.MethodGet, http.MethodOptions)
		order.HandleFunc("/{id}", orderHandler.GetOrderById).Methods(http.MethodGet, http.MethodOptions)
	}

	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%s", cfg.Server.Host, os.Getenv("SERVER_PORT")),
		Handler:           r,
		ReadTimeout:       cfg.Server.ReadTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
		IdleTimeout:       cfg.Server.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("listen: %s\\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown err:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
