package repository

import "errors"

var (
	ErrTaskNotFound = errors.New("task not found")
	ErrListNotFound = errors.New("list not found")
)
