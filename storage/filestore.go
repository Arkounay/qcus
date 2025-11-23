// Package storage provides file storage and management functionality.
package storage

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// FileStore manages uploaded files and their associated WebSocket connections
type FileStore struct {
	mu        sync.RWMutex
	files     map[string]StoredFile
	wsClients map[string][]*websocket.Conn
}

// StoredFile contains metadata about an uploaded file
type StoredFile struct {
	Path         string
	OriginalName string
	UploadTime   time.Time
}

// NewFileStore creates a new FileStore instance
func NewFileStore() *FileStore {
	return &FileStore{
		files:     make(map[string]StoredFile),
		wsClients: make(map[string][]*websocket.Conn),
	}
}

// Add stores a new file in the FileStore and schedules its automatic deletion
func (fs *FileStore) Add(id, path, originalName string, expiryMinutes int) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	fs.files[id] = StoredFile{
		Path:         path,
		OriginalName: originalName,
		UploadTime:   time.Now(),
	}

	// Schedule automatic deletion after expiry duration
	expiryDuration := time.Duration(expiryMinutes) * time.Minute
	time.AfterFunc(expiryDuration, func() {
		fs.Delete(id)
	})
}

// Get retrieves file metadata by ID
func (fs *FileStore) Get(id string) (StoredFile, bool) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	f, exists := fs.files[id]
	return f, exists
}

// Delete removes a file from storage and notifies all connected WebSocket clients
func (fs *FileStore) Delete(id string) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	f, exists := fs.files[id]
	if !exists {
		return
	}

	if err := os.Remove(f.Path); err != nil && !os.IsNotExist(err) {
		log.Printf("Error removing file %s: %v", f.Path, err)
	}

	delete(fs.files, id)

	fs.notifyClients(id)

	delete(fs.wsClients, id)
}

// AddWSClient registers a new WebSocket client for download notifications
func (fs *FileStore) AddWSClient(fileID string, conn *websocket.Conn) {
	fs.mu.Lock()
	defer fs.mu.Unlock()
	fs.wsClients[fileID] = append(fs.wsClients[fileID], conn)
}

// RemoveWSClient unregisters a WebSocket client
func (fs *FileStore) RemoveWSClient(fileID string, conn *websocket.Conn) {
	fs.mu.Lock()
	defer fs.mu.Unlock()

	clients := fs.wsClients[fileID]
	for i, c := range clients {
		if c == conn {
			// Remove this client from the slice
			fs.wsClients[fileID] = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

// notifyClients sends download notification to all WebSocket clients for a file
// Must be called with fs.mu lock held
func (fs *FileStore) notifyClients(fileID string) {
	for _, conn := range fs.wsClients[fileID] {
		err := conn.WriteJSON(map[string]interface{}{
			"downloaded": true,
			"fileID":     fileID,
		})
		if err != nil {
			log.Printf("Error sending WebSocket message: %v", err)
		}
		conn.Close()
	}
}

// BroadcastMessage sends a custom message to all WebSocket clients for a file
func (fs *FileStore) BroadcastMessage(fileID string, message interface{}) {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	for _, conn := range fs.wsClients[fileID] {
		if err := conn.WriteJSON(message); err != nil {
			log.Printf("Error broadcasting message to WebSocket client: %v", err)
		}
	}
}

// GenerateID creates a cryptographically secure random file ID
func GenerateID() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random ID: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// CleanupOldFiles removes files older than the specified duration from the upload directory
func CleanupOldFiles(uploadDir string, expiryMinutes int) error {
	log.Printf("Starting cleanup of files older than %d minutes...", expiryMinutes)

	files, err := os.ReadDir(uploadDir)
	if err != nil {
		return fmt.Errorf("error reading uploads directory: %w", err)
	}

	expiryDuration := time.Duration(expiryMinutes) * time.Minute
	now := time.Now()
	cleanedCount := 0

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(uploadDir, file.Name())
		info, err := file.Info()
		if err != nil {
			log.Printf("Error getting file info for %s: %v", file.Name(), err)
			continue
		}

		fileAge := now.Sub(info.ModTime())
		if fileAge > expiryDuration {
			if err := os.Remove(filePath); err != nil {
				log.Printf("Error removing old file %s: %v", file.Name(), err)
			} else {
				log.Printf("Removed old file: %s (age: %v)", file.Name(), fileAge)
				cleanedCount++
			}
		}
	}

	log.Printf("Cleanup completed. Removed %d old files.", cleanedCount)
	return nil
}
