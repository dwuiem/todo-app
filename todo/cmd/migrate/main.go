package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"todo/internal/config"
)

func main() {
	var configPath, migrationsPath, migrationsTable string
	flag.StringVar(&configPath, "config", "config/local.yaml", "Path to config file")
	flag.StringVar(&migrationsPath, "migrations-path", "migrations", "Set migrations path")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "Set migrations table")
	flag.Parse()

	cfg := config.MustLoad()

	if migrationsPath == "" {
		panic("Must specify migrations path")
	}

	postgresURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable&x-migrations-table=%s",
		cfg.PostgresDB.Username, cfg.PostgresDB.Password,
		cfg.PostgresDB.Host, cfg.PostgresDB.Port,
		cfg.PostgresDB.DBName, migrationsTable,
	)
	sourceURL := fmt.Sprintf(
		"file://%s",
		migrationsPath,
	)

	m, err := migrate.New(
		sourceURL,
		postgresURL,
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("Migrations table is already up to date")
			return
		}
		panic(err)
	}
	fmt.Println("Migrations applied successfully")
}
