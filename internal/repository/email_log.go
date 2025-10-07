package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/kk/mail-jack/internal/models"
)

// create a new repo by db connection
// insert log method

type EmailLogRepository struct {
	DB *sql.DB
}

func InitEmailLogRepo(db *sql.DB) *EmailLogRepository {
	query := `
	CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

	CREATE TABLE IF NOT EXISTS email_logs (
		id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
		from_email TEXT NOT NULL,
		to_email TEXT NOT NULL,
		subject TEXT,
		body TEXT,
		html TEXT,
		cc_emails JSONB,
		created_at TIMESTAMP NOT NULL
	);
	`
	db.Exec(query)
	return &EmailLogRepository{DB: db}
}

func ( r *EmailLogRepository) InsertEmailLog(log models.EmailLog) error{
	ccJSON, err := json.Marshal(log.CCEmails)
	if err != nil {
		return fmt.Errorf("failed to marshal ccEmails: %w", err)
	}

	const maxWorkers= 10
	tasks:= make(chan string, len(log.To))
	
	for _, recipient := range log.To {
		tasks <- recipient
	}
	close(tasks)

	var wg sync.WaitGroup

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			for to := range tasks {
				query := `
					INSERT INTO email_logs (from_email, to_email, subject, body, html, cc_emails, created_at)
					VALUES ($1, $2, $3, $4, $5, $6, $7)
				`
				_, err := r.DB.Exec(
					query,
					log.From,
					to,
					log.Subject,
					log.Body,
					log.HTML,
					ccJSON,
					time.Now(),
				)
				
				if err!= nil{
					fmt.Printf("failed to insert log for %s: %v\n", to, err)
				}
				fmt.Printf("Done Insert sql loggs")
			}
		}()
	}

	
	wg.Wait()
	return nil
}