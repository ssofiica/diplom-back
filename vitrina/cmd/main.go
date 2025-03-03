// package main

// import (
// 	"back/infra/postgres"
// 	"back/vitrina/config"
// 	"back/vitrina/internal/delivery"
// 	"back/vitrina/internal/repo"
// 	"back/vitrina/internal/usecase"
// 	"context"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"os"
// 	"os/signal"
// 	"syscall"

// 	"github.com/gorilla/mux"
// 	"github.com/joho/godotenv"
// )

// func init() {
// 	if err := godotenv.Load(); err != nil {
// 		fmt.Println(err)
// 	}
// }

// func main() {
// 	cfg := config.Load()
// 	db := postgres.Init(os.Getenv("PG_CONN"))
// 	defer db.Close()

// 	restRepo := repo.NewRest(db)
// 	restUC := usecase.NewRest(restRepo)
// 	restHandler := delivery.NewRestHandler(restUC)

// 	r := mux.NewRouter().PathPrefix("/api").Subrouter()

// 	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("query to path: " + r.URL.String())
// 		w.WriteHeader(http.StatusNotFound)
// 	})
// 	r.HandleFunc("/info", restHandler.GetInfo).Methods(http.MethodGet, http.MethodOptions)
// 	r.HandleFunc("/menu", restHandler.GetMenu).Methods(http.MethodGet, http.MethodOptions)
// 	//r.Use() //тут мидлвары вставить надо аргументами
// 	// order := r.PathPrefix("/order").Subrouter()
// 	// {
// 	// 	order.HandleFunc("", infoHandler.GetInfo).Methods(http.MethodGet, http.MethodOptions)
// 	// }
// 	fmt.Println(cfg.Server.Host, os.Getenv("SERVER_PORT"))
// 	srv := &http.Server{
// 		Addr:              fmt.Sprintf("%s:%s", cfg.Server.Host, os.Getenv("SERVER_PORT")),
// 		Handler:           r,
// 		ReadTimeout:       cfg.Server.ReadTimeout,
// 		WriteTimeout:      cfg.Server.WriteTimeout,
// 		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
// 		IdleTimeout:       cfg.Server.IdleTimeout,
// 	}

// 	go func() {
// 		if err := srv.ListenAndServe(); err != nil {
// 			log.Fatalf("listen: %s\\n", err)
// 		}
// 	}()

// 	// Wait for interrupt signal to gracefully shutdown the server with
// 	// a timeout of 5 seconds.
// 	quit := make(chan os.Signal, 1)
// 	// kill (no param) default send syscall.SIGTERM
// 	// kill -2 is syscall.SIGINT
// 	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
// 	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
// 	<-quit
// 	log.Println("Shutdown Server ...")

// 	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
// 	defer cancel()
// 	if err := srv.Shutdown(ctx); err != nil {
// 		log.Fatal("Server Shutdown err:", err)
// 	}
// 	// catching ctx.Done(). timeout of 5 seconds.
// 	select {
// 	case <-ctx.Done():
// 		log.Println("timeout of 5 seconds.")
// 	}
// 	log.Println("Server exiting")
// }

package main

import (
	"back/infra/postgres"
	"back/vitrina/config"
	"back/vitrina/internal/delivery"
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

	r := mux.NewRouter().PathPrefix("/api").Subrouter()

	r.HandleFunc("/ok", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("query to path: " + r.URL.String())
		w.WriteHeader(http.StatusOK)
	})
	//r.Use() //тут мидлвары вставить надо аргументами

	info := r.PathPrefix("/info").Subrouter()
	{
		info.HandleFunc("", infoHandler.GetInfo).Methods(http.MethodGet, http.MethodOptions)
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
