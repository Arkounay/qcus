// Package main is the entry point for the go-quick-cli-upload-server.
// This server provides temporary file uploads with password protection,
// automatic expiry, and real-time download notifications via WebSocket.
package main

import (
	"log"
	"net/http"
	"os"

	"go-quick-cli-upload-server/config"
	"go-quick-cli-upload-server/handlers"
	"go-quick-cli-upload-server/storage"
)

func main() {
	// Load configuration from environment variables
	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Log configuration
	log.Println("Server configuration:")
	for _, line := range cfg.LogSummary() {
		log.Printf("  %s", line)
	}

	// Create uploads directory if it doesn't exist
	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		log.Fatalf("Failed to create uploads directory: %v", err)
	}

	// Clean up old files on startup (handles crash recovery)
	if err := storage.CleanupOldFiles(cfg.UploadDir, cfg.FileExpiryMinutes); err != nil {
		log.Printf("Warning: cleanup failed: %v", err)
	}

	// Initialize file store
	store := storage.NewFileStore()

	// Create HTTP handlers
	uploadHandler := handlers.NewUploadHandler(store, cfg)
	downloadHandler := handlers.NewDownloadHandler(store)
	loginHandler := handlers.NewLoginHandler(cfg)
	wsHandler := handlers.NewWebSocketHandler(store)
	configHandler := handlers.NewConfigHandler(cfg)

	// Serve static files from public directory (Svelte build output)
	fileServer := http.FileServer(http.Dir("./public"))

	// Register routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Handle POST/PUT as file uploads
		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			uploadHandler.ServeHTTP(w, r)
			return
		}
		// Serve static files for GET requests
		fileServer.ServeHTTP(w, r)
	})

	http.Handle("/login", loginHandler)
	http.Handle("/config", configHandler)
	http.Handle("/download/", downloadHandler)
	http.Handle("/ws/", wsHandler)

	// Start server
	log.Printf("Server starting on http://localhost%s", cfg.Port)
	if err := http.ListenAndServe(cfg.Port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
