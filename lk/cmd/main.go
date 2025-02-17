package main

import (
	"back/infra/postgres"
	"back/lk/config"
	"back/lk/internal/delivery"
	"back/lk/internal/repo"
	"back/lk/internal/usecase"
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

	menuRepo := repo.NewMenu(db)
	menuUsecase := usecase.NewMenu(menuRepo)
	menuHandler := delivery.NewMenuHandler(menuUsecase)

	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("query to path: " + r.URL.String())
		w.WriteHeader(http.StatusNotFound)
	})
	//r.Use() //тут мидлвары вставить надо аргументами
	menu := r.PathPrefix("/menu").Subrouter()
	{
		menu.HandleFunc("", menuHandler.GetMenu).Methods(http.MethodGet, http.MethodOptions)
		menu.HandleFunc("/{category_id}", menuHandler.GetFoodByStatus).Methods(http.MethodGet, http.MethodOptions)
		menu.HandleFunc("/food/add", menuHandler.AddFood).Methods(http.MethodPost, http.MethodOptions)
		menu.HandleFunc("/food/{id}", menuHandler.DeleteFood).Methods(http.MethodDelete, http.MethodOptions)
		menu.HandleFunc("/category/add", menuHandler.AddCategory).Methods(http.MethodPost, http.MethodOptions)
		menu.HandleFunc("/category/{id}", menuHandler.DeleteCategory).Methods(http.MethodDelete, http.MethodOptions)
	}

	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
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
