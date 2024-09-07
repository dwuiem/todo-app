package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"todo-app/internal/model"
	"todo-app/internal/repository/postgres"
)

type ListRepository struct {
	db *sqlx.DB
}

func NewListRepository(db *sqlx.DB) *ListRepository {
	return &ListRepository{db: db}
}

func (r *ListRepository) Create(userID int, list model.List) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}

	var id int
	createListQuery := fmt.Sprintf("INSERT INTO %s (title) VALUES ($1) RETURNING id", postgres.ListsTable)
	row := tx.QueryRow(createListQuery, list.Title)
	err = row.Scan(&id)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	createUserListQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", postgres.UsersListsTable)
	_, err = tx.Exec(createUserListQuery, userID, id)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	return id, tx.Commit()
}

func (r *ListRepository) GetAll(userID int) ([]model.List, error) {
	var lists []model.List
	query := fmt.Sprintf("SELECT tl.id, tl.title FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1", postgres.ListsTable, postgres.UsersListsTable)
	err := r.db.Select(&lists, query, userID)
	return lists, err
}

func (r *ListRepository) GetByID(userID int, listID int) (model.List, error) {
	var list model.List
	query := fmt.Sprintf("SELECT tl.id, tl.title FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1 AND ul.list_id = $2", postgres.ListsTable, postgres.UsersListsTable)
	err := r.db.Get(&list, query, userID, listID)
	return list, err
}
