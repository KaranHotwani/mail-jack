package provider

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
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
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	// getting called before fn exists to free up resources
	defer cancel()

	resultCh := make(chan models.RecipientResult, len(email.To))
	wg := sync.WaitGroup{}

	for _, addr := range email.To {
		recipient := addr
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Printf("[mail-jack] sending to %s", recipient)
			input := &ses.SendEmailInput{
				Destination: &types.Destination{
					ToAddresses: []string{recipient},
				},
				Message: &types.Message{
					Subject: &types.Content{Data: aws.String(email.Subject)},
					Body: &types.Body{
						Text: &types.Content{Data: aws.String(email.Body)},
						Html: &types.Content{Data: aws.String(email.HTML)},
					},
				},
				Source: aws.String(email.From),
			}

			res, err := s.Client.SendEmail(ctx, input)
			if err != nil {
				log.Printf("[mail-jack] failed to send to %s: %v", recipient, err)
				resultCh <- models.RecipientResult{Email: recipient, Status: models.StatusFailed, MessageID: "", Error: err.Error()}
				return
			}
			log.Printf("[mail-jack] sent to %s, messageId=%s", recipient, *res.MessageId)
			resultCh <- models.RecipientResult{Email: recipient, Status: models.StatusSuccess, MessageID: *res.MessageId}
		}()
	}
	wg.Wait()
	close(resultCh)

	results := make([]models.RecipientResult, 0, len(email.To))
	successes := 0
	for i := 0; i < len(email.To); i++ {
		r := <-resultCh
		if r.Status == models.StatusSuccess {
			successes++
		}
		results = append(results, r)
	}

	overall := models.StatusSuccess
	if successes == 0 {
		overall = models.StatusFailed
	} else if successes < len(email.To) {
		overall = models.StatusPartialSuccess
	}

	log.Printf("[mail-jack] summary: total=%d success=%d failed=%d status=%s", len(email.To), successes, len(email.To)-successes, overall)

	return models.EmailResponse{Status: overall, Results: results}, nil
}
