package model

type App struct {
	ID        int    `db:"id"`
	Name      string `db:"name"`
	SecretKey string `db:"secret_key"`
}
