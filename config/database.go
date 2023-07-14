package config

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgres://username:password@localhost:5432/database_name?sslmode=disable")
	if err != nil {
		return nil, err
	}

	return db, nil
}
