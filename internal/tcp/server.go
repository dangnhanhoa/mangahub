package tcp

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"strconv"
)

type Message struct {
	Type string `json:"type"`
	UserID string `json:"user_id"`
	Token string `json:"token"`
}

type Server struct {
	port string
}

func NewServer(port int) *Server {
	return &Server{port:strconv.Itoa(port)}
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
		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn){
	defer conn.Close()

	addr := conn.RemoteAddr().String()
	log.Printf("[TCP] Client new connected from: %s", addr)
	
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
			log.Printf("[TCP] User '%s' requested to connect with token: '%s'\n", msg.UserID, msg.Token)
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