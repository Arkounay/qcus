package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"go-quick-cli-upload-server/config"
)

// ConfigHandler provides public configuration information
type ConfigHandler struct {
	Config *config.Config
}

// NewConfigHandler creates a new ConfigHandler
func NewConfigHandler(cfg *config.Config) *ConfigHandler {
	return &ConfigHandler{Config: cfg}
}

// PublicConfig contains non-sensitive configuration exposed to the frontend
type PublicConfig struct {
	IsDefaultPassword bool `json:"isDefaultPassword"`
	FileExpiryMinutes int  `json:"fileExpiryMinutes"`
	MaxFileSizeMB     int  `json:"maxFileSizeMB"`
}

// ServeHTTP implements http.Handler
func (h *ConfigHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	publicConfig := PublicConfig{
		IsDefaultPassword: h.Config.IsDefaultPassword,
		FileExpiryMinutes: h.Config.FileExpiryMinutes,
		MaxFileSizeMB:     h.Config.MaxFileSizeMB,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(publicConfig); err != nil {
		log.Printf("Error encoding config response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
