package postgres

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"todo/internal/domain/model"
	"todo/internal/storage"
)

type TaskStorage struct {
	*Storage
}

func NewTaskStorage(s *Storage) *TaskStorage {
	return &TaskStorage{s}
}

func (s *TaskStorage) GetAllByListID(c *gin.Context, listID int64) ([]model.Task, error) {
	const op = "storage.postgres.TaskStorage.GetAllByListID"

	const query = `
		SELECT id, title, description, completed, list_id FROM tasks WHERE list_id = $1
	`

	var tasks []model.Task
	rows, err := s.conn.Query(c.Request.Context(), query, listID)
	if err != nil {
		return nil, fmt.Errorf("%s %s", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var task model.Task
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Completed,
			&task.ListID,
		)
		if err != nil {
			return nil, fmt.Errorf("%s %s", op, err)
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s %s", op, err)
	}
	return tasks, nil
}

func (s *TaskStorage) Create(c *gin.Context, task model.Task) (int64, error) {
	const op = "storage.postgres.TaskStorage.Create"

	const query = `
		INSERT INTO tasks (title, description, completed) VALUES ($1, $2, $3) RETURNING id
	`
	row := s.conn.QueryRow(c.Request.Context(), query, task.Title, task.Description, task.Completed)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("%s %s", op, err)
	}
	return id, nil
}

func (s *TaskStorage) GetByID(c *gin.Context, userID, taskID int64) (model.Task, error) {
	const op = "storage.postgres.TaskStorage.GetByID"

	const query = `
		SELECT t.id, t.title, t.description, t.completed FROM tasks AS t
		INNER JOIN lists AS l on l.task_id = t.id WHERE l.user_id = $1 AND t.id = $2
	`
	row := s.conn.QueryRow(c.Request.Context(), query, userID, taskID)

	var task model.Task
	err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Completed,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Task{}, storage.ErrTaskNotFound
		}
		return model.Task{}, fmt.Errorf("%s: %w", op, err)
	}
	return task, nil
}

func (s *TaskStorage) Update(c *gin.Context, task model.Task) (int64, error) {
	const op = "storage.postgres.TaskStorage.Update"

	const query = `
		UPDATE tasks SET title = $1, description = $2, completed = $3 WHERE id = $4
	`

	tag, err := s.conn.Exec(c.Request.Context(), query, task.Title, task.Description, task.Completed, task.ID)
	if err != nil {
		return -1, fmt.Errorf("%s %s", op, err)
	}

	if tag.RowsAffected() == 0 {
		return -1, storage.ErrTaskNotFound
	}

	return task.ID, nil
}

func (s *TaskStorage) DeleteByID(c *gin.Context, taskID int64) error {
	const op = "storage.postgres.TaskStorage.DeleteByID"

	const query = `DELETE FROM tasks WHERE id = $1`

	tag, err := s.conn.Exec(c.Request.Context(), query, taskID)
	if err != nil {
		return fmt.Errorf("%s %s", op, err)
	}

	if tag.RowsAffected() == 0 {
		return storage.ErrTaskNotFound
	}

	return nil
}
