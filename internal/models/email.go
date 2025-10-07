package models

import "time"

type EmailRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	HTML    string   `json:"html"`
	CCEmails []string `json:"ccEmails"`
}

// EmailStatus is a typed enum for overall and per‑recipient status
type EmailStatus string

const (
	StatusSuccess        EmailStatus = "SUCCESS"
	StatusPartialSuccess EmailStatus = "PARTIAL_SUCCESS"
	StatusFailed         EmailStatus = "FAILED"
)

type EmailResponse struct {
	Status  EmailStatus       `json:"status"`
	Results []RecipientResult `json:"results"`
}

type RecipientResult struct {
	Email     string      `json:"email"`
	Status    EmailStatus `json:"status"`
	MessageID string      `json:"messageId"`
	Error     string      `json:"error"`
}

type  EmailLog struct { 
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Body    string   `json:"body"`
	HTML    string   `json:"html"`
	CCEmails []string `json:"ccEmails"`
	Created_at time.Time `json:"created_at"`
}