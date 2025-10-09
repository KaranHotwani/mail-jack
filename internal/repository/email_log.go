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
		status TEXT,
    	message_id TEXT,
		created_at TIMESTAMP NOT NULL
	);
	`
	db.Exec(query)
	return &EmailLogRepository{DB: db}
}

func ( r *EmailLogRepository) InsertEmailLog(log models.EmailLog, resp models.EmailResponse) error{
	ccJSON, err := json.Marshal(log.CCEmails)
	if err != nil {
		return fmt.Errorf("failed to marshal ccEmails: %w", err)
	}

	recipientMap := make(map[string]models.RecipientResult)
	for _, result := range resp.Results {
		recipientMap[result.Email] = result
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
				result := recipientMap[to]
				query := `
					INSERT INTO email_logs (from_email, to_email, subject, body, html, cc_emails, created_at, status, message_id)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
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
					string(result.Status),
					result.MessageID,
				)
				
				if err!= nil{
					fmt.Printf("failed to insert log for %s: %v\n", to, err)
				}
			}
		}()
	}

	
	wg.Wait()
	fmt.Printf("Done Insert sql logs to db")
	return nil
}