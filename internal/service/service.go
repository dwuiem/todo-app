package service

import "todo-app/internal/repository"

type Authorization interface {
}

type TodoList interface {
}

type TodoItem interface {
}

type Service struct {
	Authorization
	TodoList
	TodoItem
	Repository repository.Repository
}

func New(repos *repository.Repository) *Service {
	return &Service{
		Repository: *repos,
	}
}
