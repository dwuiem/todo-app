package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"sso/internal/domain/model"
	"sso/internal/storage"
)

const (
	ErrUniqueViolation = "23505"
)

func (s Storage) SaveUser(
	ctx context.Context,
	username string,
	passwordHash string,
) (uuid.UUID, error) {
	const op = "storage.postgres.SaveUser"

	const query = "INSERT INTO users (username, password_hash) VALUES ($1, $2) RETURNING id"
	row := s.conn.QueryRow(ctx, query, username, passwordHash)

	var id uuid.UUID
	err := row.Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == ErrUniqueViolation {
			return uuid.Nil, storage.ErrUserExists
		}
		return uuid.Nil, fmt.Errorf("%s %s", op, err)
	}

	return id, nil
}

func (s Storage) GetUser(ctx context.Context, username string) (model.User, error) {
	const op = "storage.postgres.GetUser"

	const query = "SELECT * FROM users WHERE username = $1"
	row := s.conn.QueryRow(ctx, query, username)

	var user model.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.IsAdmin,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, storage.ErrUserNotFound
		}
		return model.User{}, fmt.Errorf("%s %s", op, err)
	}

	return user, nil
}

func (s Storage) GetApp(ctx context.Context, appID int) (model.App, error) {
	const op = "storage.postgres.GetApp"

	const query = "SELECT * FROM apps WHERE id = $1"
	row := s.conn.QueryRow(ctx, query, appID)

	var app model.App
	err := row.Scan(&app.ID, &app.Name, &app.SecretKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.App{}, storage.ErrAppNotFound
		}
		return model.App{}, fmt.Errorf("%s %s", op, err)
	}

	return app, nil
}

func (s Storage) IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	const op = "storage.postgres.IsAdmin"

	const query = "SELECT is_admin FROM users WHERE id = $1"
	row := s.conn.QueryRow(ctx, query, userID)

	var isAdmin bool
	err := row.Scan(&isAdmin)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, storage.ErrUserNotFound
		}
		return false, fmt.Errorf("%s %s", op, err)
	}

	return isAdmin, nil
}
