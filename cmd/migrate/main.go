package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spf13/viper"
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
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	files, err := ioutil.ReadDir("./migrations")
	if err != nil {
		log.Fatalf("failed to read migrations: %v", err)
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
			log.Fatalf("failed to read migration %s: %v", fname, err)
		}
		if _, err := db.ExecContext(context.Background(), string(data)); err != nil {
			log.Fatalf("failed to exec migration %s: %v", fname, err)
		}
		log.Printf("applied migration: %s", fname)
	}
	log.Println("all migrations applied")
}
