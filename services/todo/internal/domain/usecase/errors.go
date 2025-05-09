package usecase

import "errors"

var (
	ErrListTitleNotValid       = errors.New("list title is not valid")
	ErrTaskDescriptionNotValid = errors.New("task description is not valid")
)
