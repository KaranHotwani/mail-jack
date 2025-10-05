package models

type EmailRequest struct {
	From    string
	To      []string
	Subject string
	Body    string
	HTML    string
}

type EmailResponse struct {
	Status    string
	MessageID string
}
