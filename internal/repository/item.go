package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"todo-app/internal/model"
	"todo-app/internal/repository/postgres"
)

type ItemRepository struct {
	db *sqlx.DB
}

func NewItemRepository(db *sqlx.DB) *ItemRepository {
	return &ItemRepository{db}
}

func (r *ItemRepository) Create(listID int, item model.Item) (int, error) {
	// Begin Transaction

	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var itemID int

	// Put item into "items" table

	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $1) RETURNING id", postgres.ItemsTable)
	row := tx.QueryRow(createItemQuery, item.Title, item.Description)
	err = row.Scan(&itemID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Put item into "lists_items" table, creating relation

	createListItemQuery := fmt.Sprintf("INSERT INTO %s (list_id, item_id) VALUES ($1, $2)", postgres.ListsItemsTable)
	_, err = tx.Exec(createListItemQuery, listID, itemID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemID, tx.Commit()
}

func (r *ItemRepository) GetAll(userID int, listID int) ([]model.Item, error) {
	var items []model.Item
	query := fmt.Sprintf("SELECT (ti.id, ti.title, ti.description, ti.completed) FROM %s ti INNER JOIN %s li on li.item_id = ti.id INNER JOIN %s ul on ul.list_id = li.id",
		postgres.ItemsTable, postgres.ListsItemsTable, postgres.ListsItemsTable)
	if err := r.db.Select(&items, query, listID, userID); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ItemRepository) GetByID(userID int, itemID int) (model.Item, error) {
	var item model.Item
	query := fmt.Sprintf("SELECT (ti.id, ti.title, ti.description, ti,completed) FROM %s ti " +
		"INNER JOIN %s li on li.item_id = ti.id INNER JOIN %s ul on ul.list_id = li.id WHERE ")
}

func (r *ItemRepository) Delete(listID int, itemID int) error {
	//TODO implement me
	panic("implement me")
}

func (r *ItemRepository) Update(listID int, input model.UpdateItemInput) error {
	//TODO implement me
	panic("implement me")
}
