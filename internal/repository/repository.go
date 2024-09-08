package repository

import (
	"github.com/jmoiron/sqlx"
	"todo-app/internal/model"
)

type Authorization interface {
	CreateUser(user model.User) (int, error)
	GetUser(username, password string) (model.User, error)
}

type List interface {
	Create(userID int, list model.List) (int, error)
	GetAll(userID int) ([]model.List, error)
	GetByID(userID, listID int) (model.List, error)
	Delete(userID, listID int) error
	Update(userID, listID int, input model.UpdateListInput) error
}

type Item interface {
	Create(listID int, item model.Item) (int, error)
	GetAll(userID, listID int) ([]model.Item, error)
	GetByID(userID, itemID int) (model.Item, error)
	Delete(userID, itemID int) error
	Update(userID, itemID int, input model.UpdateItemInput) error
}

type Repository struct {
	Authorization
	List
	Item
}

func New(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthRepository(db),
		List:          NewListRepository(db),
		Item:          NewItemRepository(db),
	}
}
