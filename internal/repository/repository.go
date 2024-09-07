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
	GetByID(userID int, listID int) (model.List, error)
}

type Item interface {
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
	}
}
