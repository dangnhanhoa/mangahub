package websocket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)


type ChatMessage struct {
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

type Client struct {
	Conn     *websocket.Conn
	UserID   string
	Username string
}

type Hub struct {
	rooms map[string]map[*Client]bool
	mu    sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		rooms: make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Register(mangaID string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.rooms[mangaID] == nil {
		h.rooms[mangaID] = make(map[*Client]bool)
	}

	h.rooms[mangaID][client] = true
	log.Printf("[WS] User %s joined room %s", client.Username, mangaID)
}

func (h *Hub) Unregister(mangaID string, client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if room, ok := h.rooms[mangaID]; ok {
		if _, exists := room[client]; exists {
			delete(room, client)
			client.Conn.Close()
			log.Printf("[WS] User %s left room %s", client.Username, mangaID)
			
			if len(room) == 0 {
				delete(h.rooms, mangaID)
			}
		}
	}
}

func (h *Hub) Broadcast(mangaID string, msg ChatMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	room, exists := h.rooms[mangaID]
	if !exists {
		return
	}

	for client := range room {
		err := client.Conn.WriteJSON(msg)
		if err != nil {
			log.Printf("[WS] Error sending message to %s: %v", client.Username, err)
		}
	}
}
