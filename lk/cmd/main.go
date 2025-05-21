package main

import (
	"back/infra/click"
	"back/infra/minio"
	"back/infra/postgres"
	"back/lk/config"
	"back/lk/internal/delivery"
	"back/lk/internal/entity"
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

	token, err := entity.NewJWT(os.Getenv("JWT_SECRET"), os.Getenv("JWT_DURATION"))
	if err != nil {
		log.Fatal("jwt troubles")
	}

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
	clickhouse := repo.NewAnalytics(clickClient)

	menuRepo := repo.NewMenu(db)
	menuUsecase := usecase.NewMenu(menuRepo, minio)
	menuHandler := delivery.NewMenuHandler(menuUsecase)

	infoRepo := repo.NewRest(db)
	infoUsecase := usecase.NewRest(infoRepo, minio)
	infoHandler := delivery.NewRestHandler(infoUsecase)

	orderRepo := repo.NewOrder(db)
	orderUsecase := usecase.NewOrder(orderRepo, clickhouse)
	orderHandler := delivery.NewOrder(orderUsecase)

	analitcsUsecase := usecase.NewAnalytics(clickhouse, menuRepo)
	analitcsHandler := delivery.NewAnalytics(analitcsUsecase)

	authRepo := repo.NewUser(db)
	authUsecase := usecase.NewUser(authRepo)
	authHandler := delivery.NewAuthHandler(authUsecase, token)

	r := mux.NewRouter().PathPrefix("/api").Subrouter()
	r.Use(middleware.CorsMiddleware)

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("query to path: " + r.URL.String())
		w.WriteHeader(http.StatusNotFound)
	})

	r.HandleFunc("/signin", authHandler.SignIn).Methods(http.MethodPost, http.MethodOptions)
	r.HandleFunc("/signup", authHandler.SignUp).Methods(http.MethodPost, http.MethodOptions)

	menu := r.PathPrefix("/menu").Subrouter()
	{
		menu.HandleFunc("", middleware.JWTMiddleware(menuHandler.GetMenu)).Methods(http.MethodGet, http.MethodOptions)
		menu.HandleFunc("/category-list", middleware.JWTMiddleware(menuHandler.GetCategoryList)).Methods(http.MethodGet, http.MethodOptions)
		menu.HandleFunc("/{category_id}", middleware.JWTMiddleware(menuHandler.GetFoodByStatus)).Methods(http.MethodGet, http.MethodOptions)
		menu.HandleFunc("/food/add", middleware.JWTMiddleware(menuHandler.AddFood)).Methods(http.MethodPost, http.MethodOptions)
		menu.HandleFunc("/food/{id}", middleware.JWTMiddleware(menuHandler.DeleteFood)).Methods(http.MethodDelete, http.MethodOptions)
		menu.HandleFunc("/food/{id}/change", middleware.JWTMiddleware(menuHandler.EditFood)).Methods(http.MethodPut, http.MethodOptions)
		menu.HandleFunc("/food/{id}/status", middleware.JWTMiddleware(menuHandler.ChangeStatus)).Methods(http.MethodPut, http.MethodOptions)
		menu.HandleFunc("/food/{id}/img", menuHandler.UploadFoodImage).Methods(http.MethodPut, http.MethodOptions)
		menu.HandleFunc("/category/add", middleware.JWTMiddleware(menuHandler.AddCategory)).Methods(http.MethodPost, http.MethodOptions)
		menu.HandleFunc("/category/{id}", menuHandler.DeleteCategory).Methods(http.MethodDelete, http.MethodOptions)
	}

	info := r.PathPrefix("/info").Subrouter()
	{
		info.HandleFunc("", middleware.JWTMiddleware(infoHandler.GetInfo)).Methods(http.MethodGet, http.MethodOptions)
		info.HandleFunc("/base", middleware.JWTMiddleware(infoHandler.UploadBaseInfo)).Methods(http.MethodPost, http.MethodOptions)
		info.HandleFunc("/upload-logo", middleware.JWTMiddleware(infoHandler.UploadImage)).Methods(http.MethodPost, http.MethodOptions)
		info.HandleFunc("/site-content", middleware.JWTMiddleware(infoHandler.UploadDescriptionsAndImages)).Methods(http.MethodPost, http.MethodOptions)
	}

	order := r.PathPrefix("/order").Subrouter()
	{
		order.HandleFunc("/mini-list", middleware.JWTMiddleware(orderHandler.GetMiniOrders)).Methods(http.MethodGet, http.MethodOptions)
		order.HandleFunc("/{id}", middleware.JWTMiddleware(orderHandler.GetOrderById)).Methods(http.MethodGet, http.MethodOptions)
		order.HandleFunc("/{id}", middleware.JWTMiddleware(orderHandler.ChangeStatus)).Methods(http.MethodPut, http.MethodOptions)
	}

	r.HandleFunc("/analytics", middleware.JWTMiddleware(analitcsHandler.GetAnalytics)).Methods(http.MethodGet, http.MethodOptions)

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
