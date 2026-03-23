package auth

import "github.com/jmoiron/sqlx"

func NewAuth(db *sqlx.DB) *Repository {
	return &Repository{Db: db}
}
