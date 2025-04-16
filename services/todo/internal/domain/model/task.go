package model

type Task struct {
	ID          int64
	Title       string
	Description string
	Completed   bool
	ListID      int64
}
