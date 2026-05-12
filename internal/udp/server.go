package udp
import (
	"encoding/json"
	"log"
	"net"
	"sync"
	"time"
	"strconv"
)
type Notification struct {
	Type      string `json:"type"`
	MangaID   string `json:"manga_id"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}
type Server struct {
	port    string
	clients map[string]*net.UDPAddr
	mu      sync.RWMutex
}
func NewServer(port int) *Server {
	return &Server{
		port:    strconv.Itoa(port),
		clients: make(map[string]*net.UDPAddr),
	}
}
func (s *Server) Start() error {
	addr, err := net.ResolveUDPAddr("udp", ":"+s.port)
	if err != nil {
		return err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return err
	}
	defer conn.Close()
	log.Printf("[UDP] Server is listening on port %s", s.port)
	buffer := make([]byte, 1024)
	for {
		
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Printf("[UDP] Data error: %v", err)
			continue
		}
		go s.handleMessage(clientAddr, buffer[:n])
	}
}

func (s *Server) handleMessage(addr *net.UDPAddr, data []byte) {
	var msg struct {
		Type   string `json:"type"`
		UserID string `json:"user_id"`
	}
	if err := json.Unmarshal(data, &msg); err != nil {
		log.Printf("[UDP] JSON error from %s: %v", addr.String(), err)
		return
	}
	
	if msg.Type == "register" {
		s.mu.Lock()
		s.clients[addr.String()] = addr
		s.mu.Unlock()
		log.Printf("[UDP] Client %s registered for notification", addr.String())
	}
}

func (s *Server) Broadcast(mangaID, message string) {
	notification := Notification{
		Type:      "chapter_update",
		MangaID:   mangaID,
		Message:   message,
		Timestamp: time.Now().Unix(),
	}
	payload, _ := json.Marshal(notification)
	s.mu.RLock()
	defer s.mu.RUnlock()
	sendConn, err := net.ListenPacket("udp", ":0")
	if err != nil {
		log.Printf("[UDP] Failed to create send socket: %v", err)
		return
	}
	defer sendConn.Close()
	for _, clientAddr := range s.clients {
		_, err := sendConn.WriteTo(payload, clientAddr)
		if err != nil {
			log.Printf("[UDP] Error sending to %s: %v", clientAddr.String(), err)
		}
	}
	log.Printf("[UDP] Broadcasted notification to %d clients", len(s.clients))
}