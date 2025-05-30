package usecase

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"strings"
	"time"
	"todo/internal/domain/entity"
)

type TaskRepository interface {
	Create(ctx context.Context, task entity.Task) (uuid.UUID, error)
	Update(ctx context.Context, task entity.Task) error
	GetAllByListID(ctx context.Context, listID uuid.UUID) ([]entity.Task, error)
	GetByID(ctx context.Context, taskID uuid.UUID) (entity.Task, error)
	ExistsByUserID(ctx context.Context, userID, taskID uuid.UUID) error
	DeleteByID(ctx context.Context, taskID uuid.UUID) error
}

type listChecker interface {
	ExistsByUserID(ctx context.Context, userID, listID uuid.UUID) error
}

type Task struct {
	repo TaskRepository
	list listChecker
	log  *slog.Logger
}

func NewTask(repo TaskRepository, list listChecker, logger *slog.Logger) *Task {
	return &Task{
		repo: repo,
		list: list,
		log:  logger,
	}
}

func (t *Task) Create(ctx context.Context, userID, listID uuid.UUID, in entity.CreateTaskIn) (uuid.UUID, error) {
	log := t.log.With(slog.String("operation", "usecase.CreateTask"), slog.With("user_id", userID))

	if err := t.list.ExistsByUserID(ctx, userID, listID); err != nil {
		log.Debug("Task does not exist")
		return uuid.Nil, err
	}

	if len(strings.TrimSpace(in.Description)) == 0 {
		log.Debug("Task description is empty")
		return uuid.Nil, ErrTaskDescriptionNotValid
	}

	if in.Deadline != nil && in.Deadline.Before(time.Now()) {
		log.Debug("Task deadline is before now")
		return uuid.Nil, ErrDeadlineNotValid
	}

	return t.repo.Create(ctx, entity.Task{
		ListID:      listID,
		CreatedAt:   time.Now(),
		Description: in.Description,
		Completed:   in.Completed,
		Deadline:    in.Deadline,
	})
}

// Update TODO: partial updating
func (t *Task) Update(ctx context.Context, userID, taskID uuid.UUID, in entity.UpdateTaskIn) error {
	if err := t.repo.ExistsByUserID(ctx, userID, taskID); err != nil {
		return err
	}
	if len(strings.TrimSpace(in.Description)) == 0 {
		return ErrTaskDescriptionNotValid
	}
	completed := false
	if in.Completed != nil {
		completed = *in.Completed
	}
	return t.repo.Update(ctx, entity.Task{
		Description: in.Description,
		Completed:   completed,
	})
}

func (t *Task) GetByListID(ctx context.Context, userID, listID uuid.UUID) ([]entity.Task, error) {
	log := t.log.With(slog.String("operation", "usecase.GetTaskByListID"), slog.With("user_id", userID))

	if err := t.list.ExistsByUserID(ctx, userID, listID); err != nil {
		log.Debug("Task does not exist")
		return nil, err
	}
	return t.repo.GetAllByListID(ctx, listID)
}

func (t *Task) Get(ctx context.Context, userID, taskID uuid.UUID) (entity.Task, error) {
	log := t.log.With(slog.String("operation", "usecase.GetTask"), slog.With("user_id", userID))

	if err := t.repo.ExistsByUserID(ctx, userID, taskID); err != nil {
		log.Debug("Task does not exist")
		return entity.Task{}, err
	}

	return t.repo.GetByID(ctx, taskID)
}

func (t *Task) Delete(ctx context.Context, userID, taskID uuid.UUID) error {
	log := t.log.With(slog.String("operation", "usecase.GetTask"), slog.With("user_id", userID))

	if err := t.repo.ExistsByUserID(ctx, userID, taskID); err != nil {
		log.Debug("Task does not exist")
		return err
	}

	return t.repo.DeleteByID(ctx, taskID)
}
