package main

import (
	"back/infra/click"
	"back/infra/minio"
	"back/infra/postgres"
	"back/lk/config"
	"back/lk/internal/delivery"
	"back/lk/internal/middleware"
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

	minioClient, err := minio.NewMinioClient(
		os.Getenv("MINIO_ENDPOINT"),
		os.Getenv("MINIO_ROOT_USER"),
		os.Getenv("MINIO_ROOT_PASSWORD"),
		os.Getenv("MINIO_BUCKET_NAME"),
	)
	if err != nil {
		log.Fatalf("Error initializing MinIO: %v", err)
	}

	clickClient, err := click.NewClickHouseClient(
		os.Getenv("HOST"),
		os.Getenv("CLICKHOUSE_PORT"),
		os.Getenv("CLICKHOUSE_DB"),
		os.Getenv("CLICKHOUSE_USER"),
		os.Getenv("CLICKHOUSE_PASSWORD"),
	)
	if err != nil {
		log.Fatalf("Failed to connect clickhouse: %v", err)
	}

	minio := repo.NewMinio(minioClient)
	clickhouse := repo.NewClickHouse(clickClient)

	menuRepo := repo.NewMenu(db)
	menuUsecase := usecase.NewMenu(menuRepo, minio)
	menuHandler := delivery.NewMenuHandler(menuUsecase)

	infoRepo := repo.NewRest(db)
	infoUsecase := usecase.NewRest(infoRepo, minio)
	infoHandler := delivery.NewRestHandler(infoUsecase)

	orderRepo := repo.NewOrder(db)
	orderUsecase := usecase.NewOrder(orderRepo, clickhouse)
	orderHandler := delivery.NewOrder(orderUsecase)

	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.Use(middleware.CorsMiddleware)

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("query to path: " + r.URL.String())
		w.WriteHeader(http.StatusNotFound)
	})
	menu := r.PathPrefix("/menu").Subrouter()
	{
		menu.HandleFunc("", menuHandler.GetMenu).Methods(http.MethodGet, http.MethodOptions)
		menu.HandleFunc("/category-list", menuHandler.GetCategoryList).Methods(http.MethodGet, http.MethodOptions)
		menu.HandleFunc("/{category_id}", menuHandler.GetFoodByStatus).Methods(http.MethodGet, http.MethodOptions)
		menu.HandleFunc("/food/add", menuHandler.AddFood).Methods(http.MethodPost, http.MethodOptions)
		menu.HandleFunc("/food/{id}", menuHandler.DeleteFood).Methods(http.MethodDelete, http.MethodOptions)
		menu.HandleFunc("/food/{id}/change", menuHandler.EditFood).Methods(http.MethodPut, http.MethodOptions)
		menu.HandleFunc("/food/{id}/status", menuHandler.ChangeStatus).Methods(http.MethodPut, http.MethodOptions)
		menu.HandleFunc("/food/{id}/img", menuHandler.UploadFoodImage).Methods(http.MethodPut, http.MethodOptions)
		menu.HandleFunc("/category/add", menuHandler.AddCategory).Methods(http.MethodPost, http.MethodOptions)
		menu.HandleFunc("/category/{id}", menuHandler.DeleteCategory).Methods(http.MethodDelete, http.MethodOptions)
	}

	info := r.PathPrefix("/info").Subrouter()
	{
		info.HandleFunc("", infoHandler.GetInfo).Methods(http.MethodGet, http.MethodOptions)
		info.HandleFunc("/base", infoHandler.UploadBaseInfo).Methods(http.MethodPost, http.MethodOptions)
		info.HandleFunc("/upload-logo", infoHandler.UploadImage).Methods(http.MethodPost, http.MethodOptions)
		info.HandleFunc("/site-content", infoHandler.UploadDescriptionsAndImages).Methods(http.MethodPost, http.MethodOptions)
	}

	order := r.PathPrefix("/order").Subrouter()
	{
		order.HandleFunc("/mini-list", orderHandler.GetMiniOrders).Methods(http.MethodGet, http.MethodOptions)
		order.HandleFunc("/{id}", orderHandler.GetOrderById).Methods(http.MethodGet, http.MethodOptions)
		order.HandleFunc("/{id}", orderHandler.ChangeStatus).Methods(http.MethodPut, http.MethodOptions)
	}

	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%s", os.Getenv("HOST"), cfg.Server.Port),
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
