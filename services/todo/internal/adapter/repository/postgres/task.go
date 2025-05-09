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

type Task struct {
	db *pgxpool.Pool
}

func NewTask(db *pgxpool.Pool) *Task {
	return &Task{db: db}
}

func (s *Task) GetAllByListID(ctx context.Context, listID uuid.UUID) ([]entity.Task, error) {
	const op = "adapter.repository.postgres.Task.GetAllByListID"

	const query = `
		SELECT id, description, completed, list_id, created_at, deadline FROM tasks WHERE list_id = $1
	`

	var tasks []entity.Task
	rows, err := s.db.Query(ctx, query, listID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var task entity.Task
		err := rows.Scan(
			&task.ID,
			&task.Description,
			&task.Completed,
			&task.ListID,
			&task.CreatedAt,
			&task.Deadline,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return tasks, nil
}

func (s *Task) Create(ctx context.Context, task entity.Task) (uuid.UUID, error) {
	const op = "adapter.repository.postgres.Task.Create"

	const query = `
		INSERT INTO tasks (description, completed, list_id, created_at, deadline)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`

	var id uuid.UUID
	err := s.db.QueryRow(ctx, query, task.Description, task.Completed, task.ListID, task.CreatedAt, task.Deadline).Scan(&id)
	if err != nil {
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (s *Task) GetByID(ctx context.Context, taskID uuid.UUID) (entity.Task, error) {
	const op = "adapter.repository.postgres.Task.GetByID"

	const query = `
		SELECT id, description, completed, list_id, created_at, deadline FROM tasks WHERE id = $1
	`
	var task entity.Task
	err := s.db.QueryRow(ctx, query, taskID).Scan(
		&task.ID,
		&task.Description,
		&task.Completed,
		&task.ListID,
		&task.CreatedAt,
		&task.Deadline,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Task{}, repository.ErrTaskNotFound
		}
		return entity.Task{}, fmt.Errorf("%s: %w", op, err)
	}
	return task, nil
}

func (s *Task) ExistsByUserID(ctx context.Context, userID, taskID uuid.UUID) error {
	const op = "adapter.repository.postgres.Task.ExistsByUserID"

	const query = `
		SELECT 1 FROM tasks t
		INNER JOIN lists l ON l.id = t.list_id
		WHERE t.id = $1 AND l.user_id = $2
	`
	var exists int
	err := s.db.QueryRow(ctx, query, taskID, userID).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrTaskNotFound
		}
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (s *Task) Update(ctx context.Context, task entity.Task) error {
	const op = "adapter.repository.postgres.Task.Update"

	const query = `
		UPDATE tasks SET description = $1, completed = $2, deadline = $3 WHERE id = $4
	`

	tag, err := s.db.Exec(ctx, query, task.Description, task.Completed, task.Deadline, task.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if tag.RowsAffected() == 0 {
		return repository.ErrTaskNotFound
	}
	return nil
}

func (s *Task) DeleteByID(ctx context.Context, taskID uuid.UUID) error {
	const op = "adapter.repository.postgres.Task.DeleteByID"

	const query = `DELETE FROM tasks WHERE id = $1`

	tag, err := s.db.Exec(ctx, query, taskID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if tag.RowsAffected() == 0 {
		return repository.ErrTaskNotFound
	}
	return nil
}
