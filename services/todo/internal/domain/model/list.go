package model

type List struct {
	ID     int64
	Title  string
	UserID int64
}

type ListUpdate struct {
	Title string
}
