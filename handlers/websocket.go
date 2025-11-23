package handlers

import (
	"log"
	"net/http"

	"go-quick-cli-upload-server/storage"

	"github.com/gorilla/websocket"
)

// WebSocketHandler handles WebSocket connections for download notifications
type WebSocketHandler struct {
	Store    *storage.FileStore
	Upgrader websocket.Upgrader
}

// NewWebSocketHandler creates a new WebSocketHandler
func NewWebSocketHandler(store *storage.FileStore) *WebSocketHandler {
	return &WebSocketHandler{
		Store: store,
		Upgrader: websocket.Upgrader{
			CheckOrigin: checkOrigin,
		},
	}
}

// ServeHTTP implements http.Handler
func (h *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fileID := r.URL.Path[len("/ws/"):]
	if fileID == "" {
		http.Error(w, "File ID required", http.StatusBadRequest)
		return
	}

	conn, err := h.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	_, exists := h.Store.Get(fileID)
	if !exists {
		// File already downloaded or doesn't exist - notify and close
		h.sendDownloadedNotification(conn, fileID)
		return
	}

	h.Store.AddWSClient(fileID, conn)
	log.Printf("WebSocket client connected for file: %s", fileID)

	h.handleConnection(fileID, conn)
}

// sendDownloadedNotification sends a "downloaded" message and closes the connection
func (h *WebSocketHandler) sendDownloadedNotification(conn *websocket.Conn, fileID string) {
	if err := conn.WriteJSON(map[string]interface{}{
		"downloaded": true,
		"fileID":     fileID,
	}); err != nil {
		log.Printf("Error writing WebSocket message: %v", err)
	}
	if err := conn.Close(); err != nil {
		log.Printf("Error closing WebSocket connection: %v", err)
	}
}

// handleConnection manages the WebSocket connection lifecycle
func (h *WebSocketHandler) handleConnection(fileID string, conn *websocket.Conn) {
	go func() {
		defer func() {
			h.Store.RemoveWSClient(fileID, conn)
			if err := conn.Close(); err != nil && !isClosedError(err) {
				log.Printf("Error closing WebSocket: %v", err)
			}
		}()

		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				if !websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) && !isClosedError(err) {
					log.Printf("WebSocket error: %v", err)
				}
				break
			}
		}
	}()
}

// isClosedError checks if error is from an already-closed connection
func isClosedError(err error) bool {
	return err != nil && (err.Error() == "use of closed network connection" ||
		websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway))
}

// checkOrigin validates WebSocket connection origin
func checkOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")

	// Allow empty origin for curl/testing tools
	if origin == "" {
		return true
	}

	// Allow connections from the same host
	return origin == "http://"+r.Host || origin == "https://"+r.Host
}
