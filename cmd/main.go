package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	// models "github.com/kk/mail-jack/internal/models"
	provider "github.com/kk/mail-jack/internal/provider"
	"github.com/kk/mail-jack/internal/service"
	"github.com/kk/mail-jack/internal/server"
	"github.com/kk/mail-jack/internal/db"
	repository "github.com/kk/mail-jack/internal/repository"
)

func main() {
	// Load .env
	godotenv.Load()

	providerName := os.Getenv("EMAIL_PROVIDER")
	if providerName == "" {
		log.Fatal("EMAIL_PROVIDER is not set in .env or environment")
	}

	var selectedProvider provider.SendEmailProvider
	var err error

	switch providerName {
	case "SES":
		selectedProvider, err = provider.NewSesProvider()
		if err != nil {
			log.Fatal("Failed to init SES:", err)
		}
	// case "sendgrid":
	//     selectedProvider, err = provider.NewSendGridProvider()
	//     if err != nil { log.Fatal("Failed to init SendGrid:", err) }
	default:
		log.Fatalf("Unsupported provider: %s", providerName)
	}

	dbConn, err := db.InitPosgres()
	if(err!= nil) {
		log.Fatal("Main failed to open postgres connection:", err)
	}
	
	repo:= repository.InitEmailLogRepo(dbConn)

	emailService := &service.EmailService{
		Providers: map[string]provider.SendEmailProvider{
			"SES": selectedProvider,
		},
		LogRepo: repo,
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	

	fmt.Printf("🚀 Email HTTP server running on port %s using provider %s\n", port, providerName)
	log.Fatal(server.StartHTTPServer(emailService, port))
	
}
