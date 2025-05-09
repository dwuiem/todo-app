package entity

import (
	"github.com/google/uuid"
	"time"
)

type Task struct {
	ID          uuid.UUID
	ListID      uuid.UUID
	CreatedAt   time.Time
	Description string
	Completed   bool
	Deadline    *time.Time
}

type CreateTaskIn struct {
	Description string     `json:"description" binding:"required"`
	Deadline    *time.Time `json:"deadline"`
	Completed   bool       `json:"completed"`
}

type UpdateTaskIn struct {
	Description string     `json:"description"`
	Deadline    *time.Time `json:"deadline"`
	Completed   *bool      `json:"completed"`
}
