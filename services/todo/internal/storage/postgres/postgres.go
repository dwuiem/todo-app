package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"todo/internal/config"
)

type Storage struct {
	conn *pgx.Conn
}

func New(cfg *config.Config) (*Storage, error) {
	const op = "storage.postgres.New"

	conn, err := pgx.Connect(context.Background(),
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.PostgresDB.Username,
			cfg.PostgresDB.Password,
			cfg.PostgresDB.Host,
			cfg.PostgresDB.Port,
			cfg.PostgresDB.DBName,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s %s", op, err)
	}

	return &Storage{conn: conn}, nil
}
