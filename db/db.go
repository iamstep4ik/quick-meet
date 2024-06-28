package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
)

func InitDB() (*sqlx.DB, error) {
	dbConnStr := os.Getenv("DATABASE_URL")
	db, err := sqlx.Open("postgres", dbConnStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping the database: %v", err)
	}

	return db, nil
}
