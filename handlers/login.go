package handlers

import (
	"log"
	"net/http"

	"go-quick-cli-upload-server/config"
)

// LoginHandler validates upload password
type LoginHandler struct {
	Config *config.Config
}

// NewLoginHandler creates a new LoginHandler
func NewLoginHandler(cfg *config.Config) *LoginHandler {
	return &LoginHandler{Config: cfg}
}

// ServeHTTP implements http.Handler
func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	providedPassword := r.Header.Get("X-Upload-Password")
	if providedPassword == "" {
		providedPassword = r.FormValue("password")
	}

	if providedPassword == h.Config.UploadPassword {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte(`{"success": true}`)); err != nil {
			log.Printf("Error writing login response: %v", err)
		}
	} else {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		log.Printf("Failed login attempt from %s", r.RemoteAddr)
	}
}
