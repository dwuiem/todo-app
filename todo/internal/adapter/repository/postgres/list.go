package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"todo/internal/adapter/repository"
	"todo/internal/domain/entity"
)

type List struct {
	db *pgxpool.Pool
}

func NewList(db *pgxpool.Pool) *List {
	return &List{db: db}
}

func (s *List) Create(ctx context.Context, list entity.List) (uuid.UUID, error) {
	const op = "storage.postgres.List.Create"

	const query = `INSERT INTO lists (title, user_id) VALUES ($1, $2) RETURNING id`

	var id uuid.UUID
	err := s.db.QueryRow(ctx, query, list.Title, list.UserID).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s %s", op, err)
	}
	return id, nil
}

func (s *List) Update(ctx context.Context, list entity.List) error {
	const op = "storage.postgres.List.Update"

	const query = `UPDATE lists SET title = $1 WHERE id = $2 AND user_id = $3`

	tag, err := s.db.Exec(ctx, query, list.Title, list.ID, list.UserID)
	if err != nil {
		return fmt.Errorf("%s %s", op, err)
	}
	if tag.RowsAffected() == 0 {
		return repository.ErrListNotFound
	}
	return nil
}

func (s *List) GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]entity.List, error) {
	const op = "storage.postgres.List.GetAllByUserID"

	const query = `SELECT id, title, user_id FROM lists WHERE user_id = $1`

	rows, err := s.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s %s", op, err)
	}
	defer rows.Close()

	var lists []entity.List
	for rows.Next() {
		var list entity.List
		if err := rows.Scan(&list.ID, &list.Title, &list.UserID); err != nil {
			return nil, fmt.Errorf("%s %s", op, err)
		}
		lists = append(lists, list)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s %s", op, err)
	}
	return lists, nil
}

func (s *List) GetByID(ctx context.Context, listID uuid.UUID) (entity.List, error) {
	const op = "storage.postgres.List.GetByID"

	const query = `SELECT id, title, user_id FROM lists WHERE id = $1`

	var list entity.List
	err := s.db.QueryRow(ctx, query, listID).Scan(&list.ID, &list.Title, &list.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.List{}, repository.ErrListNotFound
		}
		return entity.List{}, fmt.Errorf("%s %s", op, err)
	}
	return list, nil
}

func (s *List) ExistsByUserID(ctx context.Context, userID, listID uuid.UUID) error {
	const op = "storage.postgres.List.ExistsByUserID"

	const query = `SELECT 1 FROM lists WHERE id = $1 AND user_id = $2`

	var exists int
	err := s.db.QueryRow(ctx, query, listID, userID).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrListNotFound
		}
		return fmt.Errorf("%s %s", op, err)
	}
	return nil
}

func (s *List) DeleteByID(ctx context.Context, listID uuid.UUID) error {
	const op = "storage.postgres.List.DeleteByID"

	const query = `DELETE FROM lists WHERE id = $1`

	tag, err := s.db.Exec(ctx, query, listID)
	if err != nil {
		return fmt.Errorf("%s %s", op, err)
	}
	if tag.RowsAffected() == 0 {
		return repository.ErrListNotFound
	}
	return nil
}
