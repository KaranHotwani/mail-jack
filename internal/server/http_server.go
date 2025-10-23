package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/kk/mail-jack/internal/middleware"
	"github.com/kk/mail-jack/internal/models"
	"github.com/kk/mail-jack/internal/service"
)

func StartHTTPServer(svc *service.EmailService, port string) error {
	// Create the email handler function
	emailHandler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
			return
		}

		var req models.EmailRequest
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&req)
		if err != nil {
			log.Printf("Failed to decode request: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request"})
			return
		}

		providerName := os.Getenv("EMAIL_PROVIDER")
		if providerName == "" {
			log.Println("EMAIL_PROVIDER not set")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "EMAIL_PROVIDER not set"})
			return
		}

		resp, err := svc.SendEmail(req, providerName)
		if err != nil {
			log.Printf("Failed to send email: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		json.NewEncoder(w).Encode(resp)
	}

	// Wrap the handler with API key authentication middleware
	http.HandleFunc("/send-email", middleware.APIKeyAuth(emailHandler))

	return http.ListenAndServe(":"+port, nil)
}
