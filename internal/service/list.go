package service

import (
	"todo-app/internal/model"
	"todo-app/internal/repository"
)

type ListService struct {
	repo repository.List
}

func NewListService(repo repository.List) *ListService {
	return &ListService{repo: repo}
}

func (s *ListService) Create(userID int, list model.List) (int, error) {
	return s.repo.Create(userID, list)
}
