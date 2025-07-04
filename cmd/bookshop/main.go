package main

import (
	"context"
	"log"
	"net/http"
	"os"
	osignal "os/signal"
	"syscall"
	"time"

	chimid "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	redisv9 "github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"

	httpdelivery "github.com/yourorg/bookshop/internal/delivery/http"
	"github.com/yourorg/bookshop/internal/integration"
	"github.com/yourorg/bookshop/internal/repository"
	"github.com/yourorg/bookshop/internal/service"
)

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./configs/config.yaml"
	}
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("Bookshop starting...")

	// --- Postgres ---
	pgURL := "postgres://" + viper.GetString("postgres.user") + ":" + viper.GetString("postgres.password") + "@" + viper.GetString("postgres.host") + ":" + viper.GetString("postgres.port") + "/" + viper.GetString("postgres.dbname") + "?sslmode=" + viper.GetString("postgres.sslmode")
	dbpool, err := pgxpool.New(context.Background(), pgURL)
	if err != nil {
		logger.Error("failed to connect to postgres", "err", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	// --- Redis ---
	rdb := redisv9.NewClient(&redisv9.Options{
		Addr: viper.GetString("redis.addr"),
		DB:   viper.GetInt("redis.db"),
	})
	defer rdb.Close()

	// --- Kafka ---
	kafkaProducer := integration.NewKafkaProducer() // TODO: передать конфиг

	// --- Keycloak ---
	keycloak := integration.NewKeycloakClient() // TODO: передать конфиг

	// --- Интеграции ---
	redisCache := integration.NewRedisCache(rdb)

	// --- Репозитории ---
	bookRepo := repository.NewBookPostgres(dbpool)
	categoryRepo := repository.NewCategoryPostgres(dbpool)
	cartRepo := repository.NewCartPostgres(dbpool)
	orderRepo := repository.NewOrderPostgres(dbpool)

	// --- Сервисы ---
	bookService := service.NewBookService(bookRepo, categoryRepo)
	categoryService := service.NewCategoryService(categoryRepo, bookRepo)
	cartService := service.NewCartService(cartRepo, bookRepo, redisCache)
	orderService := service.NewOrderService(orderRepo, cartRepo, bookRepo, kafkaProducer)

	// --- Delivery ---
	handler := httpdelivery.NewHandler(bookService, categoryService, cartService, orderService, logger)
	auth := httpdelivery.NewAuthMiddleware(keycloak, logger)
	router := handler.Router(auth)

	// --- HTTP server ---
	srv := &http.Server{
		Addr:    viper.GetString("http.addr"),
		Handler: chimid.RequestID(chimid.Logger(router)),
	}
	go func() {
		logger.Info("HTTP server started", "addr", viper.GetString("http.addr"))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", "err", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	osignal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("shutdown error", "err", err)
	}
	logger.Info("Server exited")
}
