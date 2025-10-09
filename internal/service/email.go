package service

import (
	"fmt"

	models "github.com/kk/mail-jack/internal/models"
	provider "github.com/kk/mail-jack/internal/provider"
	repository "github.com/kk/mail-jack/internal/repository"
)

type EmailService struct {
	Providers map[string]provider.SendEmailProvider // supports multiple providers
	LogRepo   *repository.EmailLogRepository
}

// SendEmail sends email using the provider selected
func (s *EmailService) SendEmail(req models.EmailRequest, providerName string) (models.EmailResponse, error) {
	selectedProvider, ok := s.Providers[providerName]
	if !ok {
		return models.EmailResponse{Status: "failed"}, fmt.Errorf("provider %s not found", providerName)
	}
	resp, err := selectedProvider.Send(req)
    if err != nil {
		return models.EmailResponse{Status: "failed"}, fmt.Errorf("failed to send email: %w", err)
	}
	fmt.Printf("selectedProvider.Send response: %v\n", resp)
	log := models.EmailLog{
		From:     req.From,
		To:       req.To,
		Subject:  req.Subject,
		Body:     req.Body,
		HTML:     req.HTML,
		CCEmails: req.CCEmails,
	}
    if err := s.LogRepo.InsertEmailLog(log, resp); err != nil {
        // Non-fatal: email sent but logging failed
        fmt.Printf("Warning: failed to log email: %v\n", err)
    }
    return resp, nil
}
