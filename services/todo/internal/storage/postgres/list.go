package postgres

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"todo/internal/domain/model"
	"todo/internal/storage"
)

type ListStorage struct {
	*Storage
}

func NewListStorage(s *Storage) *ListStorage {
	return &ListStorage{s}
}

func (s *ListStorage) Create(c *gin.Context, list model.List) (int64, error) {
	const op = "storage.postgres.ListStorage.Create"

	const query = `INSERT INTO lists (title, user_id) VALUES ($1, $2) RETURNING id`

	row := s.conn.QueryRow(c.Request.Context(), query, list.Title, list.UserID)
	var id int64
	if err := row.Scan(&id); err != nil {
		return -1, fmt.Errorf("%s %s", op, err)
	}
	return id, nil
}

func (s *ListStorage) Update(c *gin.Context, list model.List) error {
	const op = "storage.postgres.ListStorage.Update"

	const query = `UPDATE lists SET title = $1 WHERE id = $2 AND user_id = $3`

	tag, err := s.conn.Exec(c.Request.Context(), query, list.Title, list.ID, list.UserID)
	if err != nil {
		return fmt.Errorf("%s %s", op, err)
	}
	if tag.RowsAffected() == 0 {
		return storage.ErrListNotFound
	}
	return nil
}

func (s *ListStorage) GetAllByUserID(c *gin.Context, userID int64) ([]model.List, error) {
	const op = "storage.postgres.ListStorage.GetAllByUserID"

	const query = `SELECT id, title, user_id FROM lists WHERE user_id = $1`

	rows, err := s.conn.Query(c.Request.Context(), query, userID)
	if err != nil {
		return nil, fmt.Errorf("%s %s", op, err)
	}
	defer rows.Close()

	var lists []model.List
	for rows.Next() {
		var list model.List
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

func (s *ListStorage) GetByID(c *gin.Context, userID, listID int64) (model.List, error) {
	const op = "storage.postgres.ListStorage.GetByID"

	const query = `SELECT id, title, user_id FROM lists WHERE id = $1 AND user_id = $2`

	row := s.conn.QueryRow(c.Request.Context(), query, listID, userID)

	var list model.List
	if err := row.Scan(&list.ID, &list.Title, &list.UserID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.List{}, storage.ErrListNotFound
		}
		return model.List{}, fmt.Errorf("%s %s", op, err)
	}
	return list, nil
}

func (s *ListStorage) DeleteByID(c *gin.Context, listID int64) error {
	const op = "storage.postgres.ListStorage.DeleteByID"

	const query = `DELETE FROM lists WHERE id = $1`

	tag, err := s.conn.Exec(c.Request.Context(), query, listID)
	if err != nil {
		return fmt.Errorf("%s %s", op, err)
	}
	if tag.RowsAffected() == 0 {
		return storage.ErrListNotFound
	}
	return nil
}
