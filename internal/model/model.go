package model

import "errors"

type List struct {
	ID    int    `json:"id" db:"id"`
	Title string `json:"title" db:"title" binding:"required"`
}

type User struct {
	ID       int    `json:"-" db:"id"`
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Item struct {
	ID          int    `json:"id" db:"id"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

type ListUserRelation struct {
	ListID int `json:"list_id" db:"list_id"`
	UserID int `json:"user_id" db:"user_id"`
}

type ListItemRelation struct {
	ListID int `json:"list_id" db:"list_id"`
	ItemID int `json:"item_id" db:"item_id"`
}

type UpdateListInput struct {
	Title *string `json:"title"`
}

func (i UpdateListInput) Validate() error {
	if i.Title == nil {
		return errors.New("Title is required")
	}
	return nil
}

type UpdateItemInput struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Completed   *bool   `json:"completed"`
}

func (i UpdateItemInput) Validate() error {
	if i.Title == nil && i.Description == nil && i.Completed != nil {
		return errors.New("There is nothing to update")
	}
	return nil
}
