// @title           Bookshop API
// @version         1.0
// @description     API for Bookshop service
// @BasePath        /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	osignal "os/signal"
	"sort"
	"strings"
	"syscall"
	"time"

	chimid "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	redisv9 "github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"

	httpdelivery "github.com/yourorg/bookshop/internal/delivery/http"
	"github.com/yourorg/bookshop/internal/integration"
	"github.com/yourorg/bookshop/internal/repository"
	"github.com/yourorg/bookshop/internal/service"
)

func runMigrations(ctx context.Context, dbpool *pgxpool.Pool, logger *slog.Logger) error {
	logger.Info("Running database migrations...")

	// Используем стандартный sql.DB для миграций
	db, err := sql.Open("pgx", dbpool.Config().ConnString())
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	files, err := ioutil.ReadDir("./migrations")
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var sqlFiles []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".sql") {
			sqlFiles = append(sqlFiles, f.Name())
		}
	}
	sort.Strings(sqlFiles)

	for _, fname := range sqlFiles {
		data, err := ioutil.ReadFile("./migrations/" + fname)
		if err != nil {
			return fmt.Errorf("read migration file %s: %w", fname, err)
		}
		if _, err := db.ExecContext(ctx, string(data)); err != nil {
			return fmt.Errorf("exec migration %s: %w", fname, err)
		}
		logger.Info("Applied migration", "file", fname)
	}

	logger.Info("All migrations completed")
	return nil
}

func main() {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./configs/config.yaml"
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		logger.Error("failed to read config", "err", err)
		os.Exit(1)
	}
	logger.Info("Bookshop starting...")

	// --- Postgres ---
	pgURL := "postgres://" + viper.GetString("postgres.user") + ":" + viper.GetString("postgres.password") + "@" + viper.GetString("postgres.host") + ":" + viper.GetString("postgres.port") + "/" + viper.GetString("postgres.dbname") + "?sslmode=" + viper.GetString("postgres.sslmode")
	dbpool, err := pgxpool.New(context.Background(), pgURL)
	if err != nil {
		logger.Error("failed to connect to postgres", "err", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	// --- Run migrations ---
	if err := runMigrations(context.Background(), dbpool, logger); err != nil {
		logger.Error("failed to run migrations", "err", err)
		os.Exit(1)
	}

	// --- Redis ---
	rdb := redisv9.NewClient(&redisv9.Options{
		Addr: viper.GetString("redis.addr"),
		DB:   viper.GetInt("redis.db"),
	})
	defer rdb.Close()

	// --- Kafka ---
	kafkaProducer := integration.NewKafkaProducer()

	// --- Keycloak ---
	keycloak := integration.NewKeycloakClient()

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
	cartService := service.NewCartService(cartRepo, bookRepo, redisCache, logger)
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
