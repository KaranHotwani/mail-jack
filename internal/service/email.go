package service

import (
	"fmt"
	models "github.com/kk/mail-jack/internal/models"
	provider "github.com/kk/mail-jack/internal/provider"
)

type EmailService struct {
    Providers map[string]provider.SendEmailProvider // supports multiple providers
}

// SendEmail sends email using the provider selected
func (s *EmailService) SendEmail(req models.EmailRequest, providerName string) (models.EmailResponse, error) {
    selectedProvider, ok := s.Providers[providerName]
    if !ok {
        return models.EmailResponse{Status: "failed"}, fmt.Errorf("provider %s not found", providerName)
    }
    return selectedProvider.Send(req)
}
