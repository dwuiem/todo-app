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
	//TODO implement me
	panic("implement me")
}

func (s *ItemService) GetByID(userID, listID int, itemID int) (model.Item, error) {
	//TODO implement me
	panic("implement me")
}

func (s *ItemService) Delete(userID, listID int, itemID int) error {
	//TODO implement me
	panic("implement me")
}

func (s *ItemService) Update(userID, listID int, input model.UpdateItemInput) error {
	//TODO implement me
	panic("implement me")
}
