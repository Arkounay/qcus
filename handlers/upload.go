// Package handlers contains HTTP request handlers for the upload server.
package handlers

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"go-quick-cli-upload-server/config"
	"go-quick-cli-upload-server/storage"

	"github.com/mdp/qrterminal/v3"
)

// UploadHandler handles file upload requests via POST or PUT
type UploadHandler struct {
	Store  *storage.FileStore
	Config *config.Config
}

// NewUploadHandler creates a new UploadHandler
func NewUploadHandler(store *storage.FileStore, cfg *config.Config) *UploadHandler {
	return &UploadHandler{
		Store:  store,
		Config: cfg,
	}
}

// ServeHTTP implements http.Handler
func (h *UploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !h.validatePassword(r) {
		http.Error(w, "Unauthorized: Invalid or missing password", http.StatusUnauthorized)
		log.Printf("Upload attempt with invalid password from %s", r.RemoteAddr)
		return
	}

	maxBytes := h.Config.MaxFileBytes()
	ct := r.Header.Get("Content-Type")
	isMultipart := h.isMultipartRequest(ct)

	if !isMultipart && r.ContentLength > maxBytes {
		http.Error(w, fmt.Sprintf("File too large (max: %d MB)", h.Config.MaxFileSizeMB), http.StatusRequestEntityTooLarge)
		log.Printf("Rejected upload: Content-Length %d exceeds max %d bytes", r.ContentLength, maxBytes)
		return
	}

	fileID, err := storage.GenerateID()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Printf("Failed to generate file ID: %v", err)
		return
	}

	originalName, src, closeSrc, err := h.parseUpload(r, isMultipart, maxBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer closeSrc()

	fileSize, err := h.saveFile(fileID, src, maxBytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filePath := filepath.Join(h.Config.UploadDir, fileID)
	h.Store.Add(fileID, filePath, originalName, h.Config.FileExpiryMinutes)

	h.sendSuccessResponse(w, r, fileID, originalName, fileSize)
	log.Printf("File uploaded: %s (original: %s, size: %d bytes)", fileID, originalName, fileSize)
}

// validatePassword checks if the provided password matches the configured password
func (h *UploadHandler) validatePassword(r *http.Request) bool {
	providedPassword := r.Header.Get("X-Upload-Password")
	if providedPassword == "" {
		providedPassword = r.FormValue("password")
	}
	return providedPassword == h.Config.UploadPassword
}

// isMultipartRequest checks if the Content-Type indicates multipart form data
func (h *UploadHandler) isMultipartRequest(contentType string) bool {
	return contentType != "" && len(contentType) >= len("multipart/form-data") &&
		contentType[:len("multipart/form-data")] == "multipart/form-data"
}

// parseUpload extracts the file from either multipart or raw upload
func (h *UploadHandler) parseUpload(r *http.Request, isMultipart bool, maxBytes int64) (
	originalName string, src io.ReadCloser, closeSrc func(), err error) {

	closeSrc = func() {}

	if isMultipart {
		return h.parseMultipartUpload(r, maxBytes)
	}
	return h.parseRawUpload(r)
}

// parseMultipartUpload handles multipart/form-data uploads
func (h *UploadHandler) parseMultipartUpload(r *http.Request, maxBytes int64) (
	originalName string, src io.ReadCloser, closeSrc func(), err error) {

	maxMemory := h.Config.MaxFileBytes()
	if parseErr := r.ParseMultipartForm(maxMemory); parseErr != nil {
		err = fmt.Errorf("failed to parse multipart form (max size: %d MB)", h.Config.MaxFileSizeMB)
		return
	}

	var fh *multipart.FileHeader
	if r.MultipartForm != nil {
		for _, fhs := range r.MultipartForm.File {
			if len(fhs) > 0 {
				fh = fhs[0]
				break
			}
		}
	}

	if fh == nil {
		err = fmt.Errorf("no file part in multipart form")
		return
	}

	if fh.Size > maxBytes {
		err = fmt.Errorf("file too large (max: %d MB)", h.Config.MaxFileSizeMB)
		log.Printf("Rejected upload: File %s size %d exceeds max %d bytes", fh.Filename, fh.Size, maxBytes)
		return
	}

	originalName = fh.Filename
	file, openErr := fh.Open()
	if openErr != nil {
		err = fmt.Errorf("failed to open uploaded file: %w", openErr)
		return
	}

	src = file
	closeSrc = func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Printf("Error closing uploaded file: %v", closeErr)
		}
	}
	return
}

// parseRawUpload handles raw PUT/POST uploads
func (h *UploadHandler) parseRawUpload(r *http.Request) (
	originalName string, src io.ReadCloser, closeSrc func(), err error) {

	// Extract filename from URL path
	if p := path.Base(r.URL.Path); p != "" && p != "/" {
		originalName = p
	}

	src = r.Body
	closeSrc = func() {
		if closeErr := r.Body.Close(); closeErr != nil {
			log.Printf("Error closing request body: %v", closeErr)
		}
	}
	return
}

// saveFile writes the uploaded file to disk with size limits and returns the file size
func (h *UploadHandler) saveFile(fileID string, src io.Reader, maxBytes int64) (int64, error) {
	pathToSave := filepath.Join(h.Config.UploadDir, fileID)

	f, err := os.Create(pathToSave)
	if err != nil {
		return 0, fmt.Errorf("failed to create file: %w", err)
	}

	limitedReader := io.LimitReader(src, maxBytes+1)
	written, err := io.Copy(f, limitedReader)

	if closeErr := f.Close(); closeErr != nil {
		log.Printf("Error closing file: %v", closeErr)
	}

	if err != nil {
		os.Remove(pathToSave)
		return 0, fmt.Errorf("failed to save file: %w", err)
	}

	// Check if file exceeded size limit
	if written > maxBytes {
		os.Remove(pathToSave)
		return 0, fmt.Errorf("file too large (max: %d MB)", h.Config.MaxFileSizeMB)
	}

	return written, nil
}

// sendSuccessResponse sends the upload success response with download URL and cURL command
func (h *UploadHandler) sendSuccessResponse(w http.ResponseWriter, r *http.Request, fileID, originalName string, fileSize int64) {
	if originalName == "" {
		originalName = fileID
	}

	// Determine the scheme (http or https)
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	// Check X-Forwarded-Proto header for reverse proxy setups
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}

	downloadURL := fmt.Sprintf("%s://%s/download/%s", scheme, r.Host, fileID)
	curlCommand := fmt.Sprintf("curl -o \"%s\" %s", originalName, downloadURL)
	fileSizeStr := formatFileSize(fileSize)

	w.Header().Set("Content-Type", "text/plain")

	var response string

	if isTerminalRequest(r) {
		qrCode := generateTerminalQRCode(downloadURL)
		if qrCode != "" {
			response = qrCode + "\n"
		}
	}

	response += fmt.Sprintf("File uploaded successfully!\nOriginal name: %s\nFile size: %s\nDownload URL: %s\ncURL command: %s\n",
		originalName, fileSizeStr, downloadURL, curlCommand)

	if _, err := w.Write([]byte(response)); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

// formatFileSize formats bytes into human-readable format
func formatFileSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// isTerminalRequest detects if the request came from a terminal/CLI tool
func isTerminalRequest(r *http.Request) bool {
	userAgent := strings.ToLower(r.Header.Get("User-Agent"))

	// Check for common terminal clients
	terminalClients := []string{"curl", "wget", "httpie", "powershell", "python-requests"}
	for _, client := range terminalClients {
		if strings.Contains(userAgent, client) {
			return true
		}
	}

	// Check if Accept header doesn't include HTML (browsers typically accept text/html)
	accept := strings.ToLower(r.Header.Get("Accept"))
	if accept != "" && !strings.Contains(accept, "text/html") && !strings.Contains(accept, "*/*") {
		return true
	}

	return false
}

// generateTerminalQRCode generates an ASCII QR code for terminal display
func generateTerminalQRCode(url string) string {
	var buf bytes.Buffer

	conf := qrterminal.Config{
		Level:     qrterminal.L,
		Writer:    &buf,
		BlackChar: qrterminal.BLACK,
		WhiteChar: qrterminal.WHITE,
		QuietZone: 1,
	}

	qrterminal.GenerateWithConfig(url, conf)

	return buf.String()
}
