package dbrepository

import (
	"database/sql"
	"fmt"
	"time"
)

func New(databaseURL string) *sql.DB {

	db, err := sql.Open("postgres", databaseURL)

	if err != nil {
		panic(fmt.Errorf("can't open db, %w", err))
	}

	if err := db.Ping(); err != nil {
		panic(fmt.Errorf("can't ping db, %w", err))
	}

	db.SetConnMaxLifetime(3 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}