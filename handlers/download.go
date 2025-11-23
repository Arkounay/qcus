package handlers

import (
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"go-quick-cli-upload-server/storage"
)

// DownloadHandler handles file download requests
type DownloadHandler struct {
	Store *storage.FileStore
}

// NewDownloadHandler creates a new DownloadHandler
func NewDownloadHandler(store *storage.FileStore) *DownloadHandler {
	return &DownloadHandler{Store: store}
}

// ServeHTTP implements http.Handler
func (h *DownloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fileID := r.URL.Path[len("/download/"):]
	if fileID == "" {
		http.Error(w, "File ID required", http.StatusBadRequest)
		return
	}

	sf, exists := h.Store.Get(fileID)
	if !exists {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	f, err := os.Open(sf.Path)
	if err != nil {
		http.Error(w, "Failed to open file", http.StatusInternalServerError)
		log.Printf("Error opening file %s: %v", sf.Path, err)
		return
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Printf("Error closing downloaded file: %v", err)
		}
	}()

	h.setDownloadHeaders(w, sf.OriginalName, fileID)

	if _, err := io.Copy(w, f); err != nil {
		log.Printf("Error streaming file: %v", err)
		return
	}

	log.Printf("File downloaded and deleted: %s", fileID)
	h.Store.Delete(fileID)
}

// setDownloadHeaders sets appropriate headers for file download
func (h *DownloadHandler) setDownloadHeaders(w http.ResponseWriter, originalName, fileID string) {
	name := originalName
	if name == "" {
		name = filepath.Base(fileID)
	}

	name = sanitizeFilename(name)

	w.Header().Set("Content-Disposition", mime.FormatMediaType("attachment", map[string]string{"filename": name}))
	w.Header().Set("Content-Type", "application/octet-stream")
}

// sanitizeFilename removes control characters and potential header injection attempts
func sanitizeFilename(name string) string {
	return strings.Map(func(r rune) rune {
		if r < 32 || r == 127 {
			return -1 // Remove control characters
		}
		return r
	}, name)
}
