package tcp

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"strconv"
	"sync"
)

type Message struct {
	Type string `json:"type"`
	UserID string `json:"user_id,omitempty"`
	Token string `json:"token,omitempty"`
	Data json.RawMessage `json:"data,omitempty`
}

type Server struct {
	port string
	clients map[net.Conn] string
	mu sync.RWMutex
}

func NewServer(port int) *Server {
	return &Server{
		port:strconv.Itoa(port),
		clients:make(map[net.Conn]string),
	}
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", ":" + s.port)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("[TCP] Server is Starting at port %s", s.port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[TCP] Error: %v", err)
			continue
		}

		s.mu.Lock()
		s.clients[conn] = "unknown"
		s.mu.Unlock()
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn){
	addr := conn.RemoteAddr().String()
	log.Printf("[TCP] Client new connected from: %s", addr)

	defer func(){
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
		conn.Close()
		//slog.Printf("[TCP] Client disconnected: %s\n", addr)
	}()

	
	scanner := bufio.NewScanner(conn)
	for scanner.Scan(){
		line := scanner.Text()

		var msg Message
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			log.Printf("[TCP] Error: %v", err)
			continue
		}

		switch msg.Type{
		case "connect":
			s.mu.Lock()
			s.clients[conn] = msg.UserID
			s.mu.Unlock()
			log.Printf("[TCP] User '%s' requested to connect with token: '%s'\n", msg.UserID, msg.Token)
		case "ping":
			pong, _:= json.Marshal(Message{Type: "pong"})
			conn.Write(append(pong, '\n'))
		default:
			log.Printf("[TCP] Received message type '%s' from %s\n", msg.Type, addr)
		}

		if err := scanner.Err(); err != nil {
			log.Printf("[TCP] Error: %v", err)
			break
		}
	}
	log.Printf("[TCP] Client %s has disconnected\n",addr)

}

func (s *Server) Broadcast(data interface{}){
	payload, err := json.Marshal(data)
	if err  != nil {
		log.Printf("[TCP] Broadcast encode error: %v", err)
		return 
	}

	msg := Message{
		Type: "update",
		Data: payload,
	}

	msgJSON, _ := json.Marshal(msg)
	msgJSON = append(msgJSON, '\n')

	s.mu.RLock()
	defer s.mu.RUnlock()

	for conn, userID := range s.clients{
		_, err := conn.Write(msgJSON)
		if err != nil {
			log.Printf("[TCP] Fail to send broadcast to %s: %v",userID,err)
		}
	}

	log.Printf("[TCP] Broadcasted %d clients", len(s.clients))

	
}