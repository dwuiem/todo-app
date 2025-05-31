package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"sso/internal/adapter/repository"
	"sso/internal/app/config"
	"sso/internal/domain/entity"
)

const (
	ErrUniqueViolation = "23505"
)

type Auth struct {
	db *pgxpool.Pool
}

func New(cfg *config.Config) (*Auth, error) {
	const op = "repository.postgres.New"

	pool, err := pgxpool.New(context.Background(),
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.PostgresDB.Username,
			cfg.PostgresDB.Password,
			cfg.PostgresDB.Host,
			cfg.PostgresDB.Port,
			cfg.PostgresDB.DBName),
	)

	if err != nil {
		return nil, fmt.Errorf("%s %s", op, err)
	}

	return &Auth{db: pool}, nil
}

func (s *Auth) SaveUser(
	ctx context.Context,
	username string,
	passwordHash string,
) (uuid.UUID, error) {
	const op = "repository.postgres.SaveUser"

	const query = "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id"
	row := s.db.QueryRow(ctx, query, username, passwordHash)

	var id uuid.UUID
	err := row.Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == ErrUniqueViolation {
			return uuid.Nil, repository.ErrUserExists
		}
		return uuid.Nil, fmt.Errorf("%s %s", op, err)
	}

	return id, nil
}

func (s *Auth) GetUser(ctx context.Context, username string) (entity.User, error) {
	const op = "repository.postgres.GetUser"

	const query = "SELECT * FROM users WHERE username = $1"
	row := s.db.QueryRow(ctx, query, username)

	var user entity.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.IsAdmin,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.User{}, repository.ErrUserNotFound
		}
		return entity.User{}, fmt.Errorf("%s %s", op, err)
	}

	return user, nil
}

func (s *Auth) GetApp(ctx context.Context, appID int) (entity.App, error) {
	const op = "repository.postgres.GetApp"

	const query = "SELECT * FROM apps WHERE id = $1"
	row := s.db.QueryRow(ctx, query, appID)

	var app entity.App
	err := row.Scan(&app.ID, &app.Name, &app.SecretKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.App{}, repository.ErrAppNotFound
		}
		return entity.App{}, fmt.Errorf("%s %s", op, err)
	}

	return app, nil
}

func (s *Auth) IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	const op = "repository.postgres.IsAdmin"

	const query = "SELECT is_admin FROM users WHERE id = $1"
	row := s.db.QueryRow(ctx, query, userID)

	var isAdmin bool
	err := row.Scan(&isAdmin)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, repository.ErrUserNotFound
		}
		return false, fmt.Errorf("%s %s", op, err)
	}

	return isAdmin, nil
}
