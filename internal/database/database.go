package database

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// connecting postgresql driver to project to communicate with db
func Connect(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
