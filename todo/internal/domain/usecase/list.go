package usecase

import (
	"context"
	"log/slog"
	"strings"
	"todo/internal/domain/entity"

	"github.com/google/uuid"
)

type ListRepository interface {
	Create(ctx context.Context, list entity.List) (uuid.UUID, error)
	GetAllByUserID(ctx context.Context, userID uuid.UUID) ([]entity.List, error)
	GetByID(ctx context.Context, listID uuid.UUID) (entity.List, error)
	ExistsByUserID(ctx context.Context, userID, listID uuid.UUID) error
	Update(ctx context.Context, list entity.List) error
	DeleteByID(ctx context.Context, listID uuid.UUID) error
}

type List struct {
	repo ListRepository
	log  *slog.Logger
}

func NewList(repo ListRepository, logger *slog.Logger) *List {
	return &List{
		repo: repo,
		log:  logger,
	}
}

func (l *List) ExistsByUserID(ctx context.Context, userID, listID uuid.UUID) error {
	return l.repo.ExistsByUserID(ctx, userID, listID)
}

func (l *List) Create(ctx context.Context, userID uuid.UUID, in entity.CreateListIn) (uuid.UUID, error) {
	log := l.log.With(slog.String("operation", "usecase.CreateList"), slog.Any("user_id", userID))

	if len(strings.TrimSpace(in.Title)) == 0 {
		log.Debug("Title is not valid")
		return uuid.Nil, ErrListTitleNotValid
	}

	return l.repo.Create(ctx, entity.List{
		Title:  in.Title,
		UserID: userID,
	})
}

func (l *List) GetAll(ctx context.Context, userID uuid.UUID) ([]entity.List, error) {
	return l.repo.GetAllByUserID(ctx, userID)
}

func (l *List) Get(ctx context.Context, userID uuid.UUID, listID uuid.UUID) (entity.List, error) {
	log := l.log.With(slog.String("operation", "usecase.GetAllLists"), slog.Any("user_id", userID))

	if err := l.repo.ExistsByUserID(ctx, userID, listID); err != nil {
		log.Debug("List not found")
		return entity.List{}, err
	}

	return l.repo.GetByID(ctx, listID)
}

func (l *List) Update(ctx context.Context, userID uuid.UUID, in entity.UpdateListIn) error {
	log := l.log.With(slog.String("operation", "usecase.UpdateList"), slog.Any("user_id", userID))

	if len(strings.TrimSpace(in.Title)) == 0 {
		log.Debug("Title is not valid")
		return ErrListTitleNotValid
	}

	return l.repo.Update(ctx, entity.List{
		Title:  in.Title,
		UserID: userID,
	})
}

func (l *List) Delete(ctx context.Context, userID, listID uuid.UUID) error {
	log := l.log.With(slog.String("operation", "usecase.UpdateList"), slog.Any("user_id", userID))

	if err := l.repo.ExistsByUserID(ctx, userID, listID); err != nil {
		log.Debug("User has no list")
		return err
	}

	return l.repo.DeleteByID(ctx, listID)
}
