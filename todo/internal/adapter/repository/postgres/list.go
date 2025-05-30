package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"todo/internal/adapter/repository"
	"todo/internal/domain/entity"
)

type List struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewList(db *pgxpool.Pool, logger *slog.Logger) *List {
	return &List{
		db:  db,
		log: logger,
	}
}

func (s *List) Create(ctx context.Context, list entity.List) (uuid.UUID, error) {
	const op = "storage.postgres.List.Create"
	log := s.log.With(slog.String("operation", op), slog.String("list id", list.ID.String()))

	const query = `INSERT INTO lists (title, user_id) VALUES ($1, $2) RETURNING id`

	var id uuid.UUID
	err := s.db.QueryRow(ctx, query, list.Title, list.UserID).Scan(&id)
	if err != nil {
		log.Error("Failed to insert new list")
		return uuid.Nil, fmt.Errorf("%s %s", op, err)
	}

	log.Debug("Inserted new list")
	return id, nil
}

func (s *List) Update(ctx context.Context, list entity.List) error {
	const op = "storage.postgres.List.Update"
	log := s.log.With(
		slog.String("operation", op),
		slog.String("list id", list.ID.String()),
		slog.String("user id", list.UserID.String()),
	)

	const query = `UPDATE lists SET title = $1 WHERE id = $2 AND user_id = $3`

	tag, err := s.db.Exec(ctx, query, list.Title, list.ID, list.UserID)
	if err != nil {
		log.Error("Failed to update list", slog.Any("error", err))
		return fmt.Errorf("%s %s", op, err)
	}
	if tag.RowsAffected() == 0 {
		log.Debug("List not found")
		return repository.ErrListNotFound
	}

	log.Debug("Updated list")
	return nil
}

func (s *List) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]entity.List, error) {
	const op = "storage.postgres.List.GetAllByUserID"
	log := s.log.With(slog.String("operation", op), slog.String("user id", userID.String()))

	const query = `SELECT id, title, user_id FROM lists WHERE user_id = $1`

	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		log.Error("Failed to get lists", slog.Any("error", err))
		return nil, fmt.Errorf("%s %s", op, err)
	}
	defer rows.Close()

	var lists []entity.List
	for rows.Next() {
		var list entity.List
		if err := rows.Scan(&list.ID, &list.Title, &list.UserID); err != nil {
			log.Error("Failed to scan row", slog.Any("error", err))
			return nil, fmt.Errorf("%s %s", op, err)
		}
		lists = append(lists, list)
	}
	if err := rows.Err(); err != nil {
		log.Error("Rows iteration error", slog.Any("error", err))
		return nil, fmt.Errorf("%s %s", op, err)
	}

	log.Debug("Retrieved lists", slog.Int("count", len(lists)))
	return lists, nil
}

func (s *List) GetByID(ctx context.Context, listID uuid.UUID) (entity.List, error) {
	const op = "storage.postgres.List.GetByID"
	log := s.log.With(slog.String("operation", op), slog.String("list id", listID.String()))

	const query = `SELECT id, title, user_id FROM lists WHERE id = $1`

	var list entity.List
	err := s.db.QueryRow(ctx, query, listID).Scan(&list.ID, &list.Title, &list.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug("List not found")
			return entity.List{}, repository.ErrListNotFound
		}
		log.Error("Failed to get list", slog.Any("error", err))
		return entity.List{}, fmt.Errorf("%s %s", op, err)
	}

	log.Debug("Retrieved list")
	return list, nil
}

func (s *List) ExistsByUserID(ctx context.Context, userID, listID uuid.UUID) error {
	const op = "storage.postgres.List.ExistsByUserID"
	log := s.log.With(
		slog.String("operation", op),
		slog.String("list id", listID.String()),
		slog.String("user id", userID.String()),
	)

	const query = `SELECT 1 FROM lists WHERE id = $1 AND user_id = $2`

	var exists int
	err := s.db.QueryRow(ctx, query, listID, userID).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug("List not found")
			return repository.ErrListNotFound
		}
		log.Error("Failed to check list existence", slog.Any("error", err))
		return fmt.Errorf("%s %s", op, err)
	}

	log.Debug("List exists")
	return nil
}

func (s *List) DeleteByID(ctx context.Context, listID uuid.UUID) error {
	const op = "storage.postgres.List.DeleteByID"
	log := s.log.With(slog.String("operation", op), slog.String("list id", listID.String()))

	const query = `DELETE FROM lists WHERE id = $1`

	tag, err := s.db.Exec(ctx, query, listID)
	if err != nil {
		log.Error("Failed to delete list", slog.Any("error", err))
		return fmt.Errorf("%s %s", op, err)
	}
	if tag.RowsAffected() == 0 {
		log.Debug("List not found")
		return repository.ErrListNotFound
	}

	log.Debug("Deleted list")
	return nil
}
