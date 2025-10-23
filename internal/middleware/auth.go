package middleware

import (
	"encoding/json"
	"net/http"
	"os"
)

// APIKeyAuth middleware checks for MAIL_JACK_API_KEY environment variable
// against X-API-KEY request header
func APIKeyAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get the API key from environment variable
		expectedAPIKey := os.Getenv("MAIL_JACK_API_KEY")
		if expectedAPIKey == "" {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "API key not configured"})
			return
		}

		// Get the API key from request header
		providedAPIKey := r.Header.Get("X-API-KEY")
		if providedAPIKey == "" {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "X-API-KEY header is required"})
			return
		}

		// Compare the API keys
		if providedAPIKey != expectedAPIKey {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid API key"})
			return
		}

		// API key is valid, proceed to the next handler
		next.ServeHTTP(w, r)
	}
}
