package model

type User struct {
	ID           int64  `db:"id"`
	Username     string `db:"username"`
	PasswordHash string `db:"password_hash"`
	IsAdmin      bool   `db:"is_admin"`
}
