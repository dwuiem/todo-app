package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"todo-app/internal/config"
)

func New(cfg config.Config) *sqlx.DB {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.PostgresDB.Host, cfg.PostgresDB.Port, cfg.PostgresDB.Username, cfg.PostgresDB.Password, cfg.PostgresDB.DBName))
	if err != nil {
		log.Fatal("Can not connect to postgres", err.Error())
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Can not connect to postgres")
	}
	return db
}
