package service

import (
	"todo-app/internal/model"
	"todo-app/internal/repository"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (int, error)
}

type List interface {
	Create(userID int, list model.List) (int, error)
	GetAll(userID int) ([]model.List, error)
	GetByID(userID, listID int) (model.List, error)
	Delete(userID, listID int) error
	Update(userID int, listID int, input model.UpdateListInput) error
}

type Item interface {
}

type Service struct {
	Authorization
	List
	Item
}

func New(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		List:          NewListService(repos.List),
	}
}
