package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

func InitPosgres() (*sql.DB, error) {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
		return nil, fmt.Errorf("DATABASE_URL environment variable is required")
	}

	fmt.Printf("Using DATABASE_URL for connection\n")

	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		log.Fatal("Failed to open postgres connection:", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping postgres:", err)
		return nil, err
	}
	return db, nil

}
