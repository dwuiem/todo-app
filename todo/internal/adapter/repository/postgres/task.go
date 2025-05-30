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

type Task struct {
	db  *pgxpool.Pool
	log *slog.Logger
}

func NewTask(db *pgxpool.Pool, logger *slog.Logger) *Task {
	return &Task{db: db, log: logger}
}

func (s *Task) GetAllByListID(ctx context.Context, listID uuid.UUID) ([]entity.Task, error) {
	const op = "adapter.repository.postgres.Task.GetAllByListID"
	log := s.log.With(slog.String("operation", op), slog.String("listID", listID.String()))

	const query = `
		SELECT id, description, completed, list_id, created_at, deadline FROM tasks WHERE list_id = $1
	`

	rows, err := s.db.Query(ctx, query, listID)
	if err != nil {
		log.Error("Failed to get all tasks", slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var tasks []entity.Task
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
			log.Error("Failed to scan task row", slog.Any("error", err))
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		tasks = append(tasks, task)
	}
	if err := rows.Err(); err != nil {
		log.Error("Row iteration error", slog.Any("error", err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Debug("Fetched all tasks for list")
	return tasks, nil
}

func (s *Task) Create(ctx context.Context, task entity.Task) (uuid.UUID, error) {
	const op = "adapter.repository.postgres.Task.Create"
	log := s.log.With(slog.String("operation", op), slog.String("listID", task.ListID.String()))

	const query = `
		INSERT INTO tasks (description, completed, list_id, created_at, deadline)
		VALUES ($1, $2, $3, $4, $5) RETURNING id
	`

	var id uuid.UUID
	err := s.db.QueryRow(ctx, query, task.Description, task.Completed, task.ListID, task.CreatedAt, task.Deadline).Scan(&id)
	if err != nil {
		log.Error("Failed to insert new task", slog.Any("error", err))
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Debug("Inserted new task", slog.String("taskID", id.String()))
	return id, nil
}

func (s *Task) GetByID(ctx context.Context, taskID uuid.UUID) (entity.Task, error) {
	const op = "adapter.repository.postgres.Task.GetByID"
	log := s.log.With(slog.String("operation", op), slog.String("taskID", taskID.String()))

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
			log.Debug("Task not found")
			return entity.Task{}, repository.ErrTaskNotFound
		}
		log.Error("Failed to get task", slog.Any("error", err))
		return entity.Task{}, fmt.Errorf("%s: %w", op, err)
	}
	log.Debug("Fetched task")
	return task, nil
}

func (s *Task) ExistsByUserID(ctx context.Context, userID, taskID uuid.UUID) error {
	const op = "adapter.repository.postgres.Task.ExistsByUserID"
	log := s.log.With(slog.String("operation", op), slog.String("taskID", taskID.String()), slog.String("userID", userID.String()))

	const query = `
		SELECT 1 FROM tasks t
		INNER JOIN lists l ON l.id = t.list_id
		WHERE t.id = $1 AND l.user_id = $2
	`
	var exists int
	err := s.db.QueryRow(ctx, query, taskID, userID).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug("Task not found for user")
			return repository.ErrTaskNotFound
		}
		log.Error("Failed to check task existence", slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}
	log.Debug("Task exists for user")
	return nil
}

func (s *Task) Update(ctx context.Context, task entity.Task) error {
	const op = "adapter.repository.postgres.Task.Update"
	log := s.log.With(slog.String("operation", op), slog.String("taskID", task.ID.String()))

	const query = `
		UPDATE tasks SET description = $1, completed = $2, deadline = $3 WHERE id = $4
	`

	tag, err := s.db.Exec(ctx, query, task.Description, task.Completed, task.Deadline, task.ID)
	if err != nil {
		log.Error("Failed to update task", slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}
	if tag.RowsAffected() == 0 {
		log.Debug("Task to update not found")
		return repository.ErrTaskNotFound
	}
	log.Debug("Task updated")
	return nil
}

func (s *Task) DeleteByID(ctx context.Context, taskID uuid.UUID) error {
	const op = "adapter.repository.postgres.Task.DeleteByID"
	log := s.log.With(slog.String("operation", op), slog.String("taskID", taskID.String()))

	const query = `DELETE FROM tasks WHERE id = $1`

	tag, err := s.db.Exec(ctx, query, taskID)
	if err != nil {
		log.Error("Failed to delete task", slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}
	if tag.RowsAffected() == 0 {
		log.Debug("Task to delete not found")
		return repository.ErrTaskNotFound
	}
	log.Debug("Task deleted")
	return nil
}
