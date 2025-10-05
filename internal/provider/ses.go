package provider

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	models "github.com/kk/mail-jack/internal/models"
)

type SESProvider struct {
	Client *ses.Client
}

func NewSesProvider() (*SESProvider, error) {
	region := os.Getenv("AWS_REGION")
	if region == "" {
		return nil, fmt.Errorf("AWS_REGION environment variable is required")
	}

	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %v", err)
	}

	client := ses.NewFromConfig(cfg)

	// Return SESProvider instance
	return &SESProvider{Client: client}, nil
}

func (s *SESProvider) Send(email models.EmailRequest) (models.EmailResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	input := &ses.SendEmailInput{
		Destination: &types.Destination{
			ToAddresses: email.To,
		},
		Message: &types.Message{
			Subject: &types.Content{
				Data: aws.String(email.Subject),
			},
			Body: &types.Body{
				Text: &types.Content{
					Data: aws.String(email.Body),
				},
				Html: &types.Content{
					Data: aws.String(email.HTML),
				},
			},
		},
		Source: aws.String(email.From),
	}
	result, err := s.Client.SendEmail(ctx, input)
	if err != nil {
		return models.EmailResponse{}, err
	}
	return models.EmailResponse{
		Status:    "success",
		MessageID: *result.MessageId,
	}, nil
}
