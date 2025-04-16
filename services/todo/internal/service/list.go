package service

import (
	"todo/internal/storage/postgres"
)

type ListService struct {
	storage *postgres.ListStorage
}

func NewListService(storage *postgres.ListStorage) *ListService {
	return &ListService{storage: storage}
}
