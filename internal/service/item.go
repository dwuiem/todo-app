package service

import (
	"todo-app/internal/model"
	"todo-app/internal/repository"
)

type ItemService struct {
	itemRepo repository.Item
	listRepo repository.List
}

func NewItemService(itemRepo repository.Item, listRepo repository.List) *ItemService {
	return &ItemService{itemRepo: itemRepo, listRepo: listRepo}
}

func (s *ItemService) Create(userID, listID int, item model.Item) (int, error) {
	_, err := s.listRepo.GetByID(userID, listID)
	if err != nil {
		return 0, err
	}
	return s.itemRepo.Create(listID, item)
}

func (s *ItemService) GetAll(userID, listID int) ([]model.Item, error) {
	return s.itemRepo.GetAll(userID, listID)
}

func (s *ItemService) GetByID(userID, itemID int) (model.Item, error) {
	return s.itemRepo.GetByID(userID, itemID)
}

func (s *ItemService) Delete(userID, itemID int) error {
	return s.itemRepo.Delete(userID, itemID)
}

func (s *ItemService) Update(userID, itemID int, input model.UpdateItemInput) error {
	return s.itemRepo.Update(userID, itemID, input)
}
