package postgres

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
	"todo-app/internal/config"
)

const (
	UsersTable      = "users"
	ListsTable      = "lists"
	UsersListsTable = "users_lists"
	ItemsTable      = "items"
	ListsItemsTable = "lists_items"
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

	createTablesQuery := `
	CREATE TABLE IF NOT EXISTS ` + UsersTable + ` (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		username TEXT NOT NULL UNIQUE,
		password_hash TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS ` + ListsTable + ` (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS ` + UsersListsTable + ` (
		user_id INT REFERENCES ` + UsersTable + `(id) ON DELETE CASCADE,
		list_id INT REFERENCES ` + ListsTable + `(id) ON DELETE CASCADE,
		PRIMARY KEY (user_id, list_id)
	);
	CREATE TABLE IF NOT EXISTS ` + ItemsTable + ` (
		id SERIAL PRIMARY KEY,
		description TEXT NOT NULL,
		completed BOOLEAN NOT NULL DEFAULT FALSE
	);
	CREATE TABLE IF NOT EXISTS ` + ListsItemsTable + ` (
		list_id INT REFERENCES ` + ListsTable + `(id) ON DELETE CASCADE,
		item_id INT REFERENCES ` + ItemsTable + `(id) ON DELETE CASCADE,
		PRIMARY KEY (list_id, item_id)
	);
	`
	_, err = db.Exec(createTablesQuery)
	if err != nil {
		log.Fatal("Can not create table", err.Error())
	}
	log.Println("Tables were created/loaded")
	return db
}
