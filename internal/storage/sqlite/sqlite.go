package sqlite

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"
	db, err := sql.Open("sqlite3", "./url-shortener.db")
	if err != nil {
		return nil, fmt.Errorf(" %s: %w ", op, err)
	}
}
