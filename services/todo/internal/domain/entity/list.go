package entity

import "github.com/google/uuid"

type List struct {
	ID     uuid.UUID
	UserID uuid.UUID
	Title  string
}

type CreateListIn struct {
	Title string `json:"title"`
}

type UpdateListIn struct {
	Title string `json:"title"`
}

type ListOut struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
}

type ListsOut struct {
	lists []ListOut
}
