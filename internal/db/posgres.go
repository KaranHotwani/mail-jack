package db

import (
	"fmt"
	"os"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitPosgres() (*sql.DB, error) {
	host:= os.Getenv("DB_HOST")
	port:= os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")


	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Failed to open postgres connection:", err)
		return nil, err
	}

	if err:= db.Ping(); err!= nil {
		log.Fatal("Failed to ping postgres:", err)
		return nil, err
	}
	return db, nil

}