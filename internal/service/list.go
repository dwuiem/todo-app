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

func (s *ListService) GetAll(userID int) ([]model.List, error) {
	return s.repo.GetAll(userID)
}

func (s *ListService) GetByID(userID int, listID int) (model.List, error) {
	return s.repo.GetByID(userID, listID)
}

func (s *ListService) Delete(userID int, listID int) error {
	return s.repo.Delete(userID, listID)
}

func (s *ListService) Update(userID int, listID int, input model.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userID, listID, input)
}
