package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/viper"
	"golang.org/x/exp/slog"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./configs/config.yaml"
	}
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		logger.Error("failed to read config", "err", err)
		os.Exit(1)
	}
	pgURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		viper.GetString("postgres.user"),
		viper.GetString("postgres.password"),
		viper.GetString("postgres.host"),
		viper.GetString("postgres.port"),
		viper.GetString("postgres.dbname"),
		viper.GetString("postgres.sslmode"),
	)
	db, err := sql.Open("pgx", pgURL)
	if err != nil {
		logger.Error("failed to connect to db", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	files, err := ioutil.ReadDir("./migrations")
	if err != nil {
		logger.Error("failed to read migrations", "err", err)
		os.Exit(1)
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
			logger.Error("failed to read migration", "file", fname, "err", err)
			os.Exit(1)
		}
		if _, err := db.ExecContext(context.Background(), string(data)); err != nil {
			logger.Error("failed to exec migration", "file", fname, "err", err)
			os.Exit(1)
		}
		logger.Info("applied migration", "file", fname)
	}
	logger.Info("all migrations applied")
}
