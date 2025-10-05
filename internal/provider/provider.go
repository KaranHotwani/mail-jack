package provider

import models "github.com/kk/mail-jack/internal/models"

type SendEmailProvider interface {
	Send(email models.EmailRequest) (models.EmailResponse, error)
}
