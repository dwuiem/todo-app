package service

import (
	"todo/internal/storage/postgres"
)

type TaskService struct {
	storage *postgres.TaskStorage
}

func NewTaskService(storage *postgres.TaskStorage) *TaskService {
	return &TaskService{
		storage: storage,
	}
}
