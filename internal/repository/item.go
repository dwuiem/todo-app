package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"todo-app/internal/model"
	"todo-app/internal/repository/postgres"
)

// Implements Item interface

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

	createItemQuery := fmt.Sprintf(`
		INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id
		`, postgres.ItemsTable)

	row := tx.QueryRow(createItemQuery, item.Title, item.Description)
	err = row.Scan(&itemID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Put item into "lists_items" table, creating relation

	createListItemQuery := fmt.Sprintf(`
		INSERT INTO %s (list_id, item_id) VALUES ($1, $2)
		`, postgres.ListsItemsTable)

	_, err = tx.Exec(createListItemQuery, listID, itemID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	return itemID, tx.Commit()
}

func (r *ItemRepository) GetAll(userID int, listID int) ([]model.Item, error) {
	var items []model.Item

	// Get all items from "items" table

	query := fmt.Sprintf(`
		SELECT ti.id, ti.title, ti.description, ti.completed FROM %s ti 
		INNER JOIN %s li on li.item_id = ti.id 
		INNER JOIN %s ul on ul.list_id = li.list_id 
		WHERE li.list_id = $1 AND ul.user_id = $2
		`, postgres.ItemsTable, postgres.ListsItemsTable, postgres.UsersListsTable)

	if err := r.db.Select(&items, query, listID, userID); err != nil {
		return nil, err
	}

	return items, nil
}

func (r *ItemRepository) GetByID(userID int, itemID int) (model.Item, error) {
	var item model.Item

	query := fmt.Sprintf(`
		SELECT ti.id, ti.title, ti.description, ti.completed FROM %s ti 
        INNER JOIN %s li on li.item_id = ti.id 
        INNER JOIN %s ul on ul.list_id = li.list_id 
        WHERE ti.id = $1 AND ul.user_id = $2
        `, postgres.ItemsTable, postgres.ListsItemsTable, postgres.UsersListsTable)

	if err := r.db.Get(&item, query, itemID, userID); err != nil {
		return item, err
	}

	return item, nil
}

func (r *ItemRepository) Delete(userID int, itemID int) error {
	query := fmt.Sprintf(`
		DELETE FROM %s ti USING %s li, %s ul 
		WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $1 AND ti.id = $2
		`, postgres.ItemsTable, postgres.ListsItemsTable, postgres.UsersListsTable)
	_, err := r.db.Exec(query, userID, itemID)
	return err
}

func (r *ItemRepository) Update(userID, itemID int, input model.UpdateItemInput) error {
	setValues := make([]string, 0) // set key=value into table
	args := make([]interface{}, 0)
	argID := 1

	// Fill slices

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argID))
		args = append(args, *input.Title)
		argID++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argID))
		args = append(args, *input.Description)
		argID++
	}

	if input.Completed != nil {
		setValues = append(setValues, fmt.Sprintf("completed=$%d", argID))
		args = append(args, *input.Completed)
		argID++
	}

	setQuery := strings.Join(setValues, ", ")
	args = append(args, userID, itemID)

	query := fmt.Sprintf(`
		UPDATE %s ti SET %s FROM %s li, %s ul
		WHERE ti.id = li.item_id AND li.list_id = ul.list_id AND ul.user_id = $%d AND ti.id = $%d`,
		postgres.ItemsTable, setQuery, postgres.ListsItemsTable, postgres.UsersListsTable, argID, argID+1)

	_, err := r.db.Exec(query, args...)
	return err
}
