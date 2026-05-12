package tcp

import (
	"encoding/json"
	"log"
	"net"
)

type Client struct {
	addr string
}

func NewClient(addr string) *Client {
	return &Client{addr: addr}
}

func (c *Client) Broadcast(data interface{}) {
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		log.Printf("[TCP Client] Failed to dial %s: %v", c.addr, err)
		return
	}
	defer conn.Close()

	payloadData, err := json.Marshal(data)
	if err != nil {
		log.Printf("[TCP Client] Failed to marshal payload: %v", err)
		return
	}

	msg := Message{
		Type: "broadcast",
		Data: payloadData,
	}

	msgJSON, err := json.Marshal(msg)
	if err != nil {
		return
	}

	msgJSON = append(msgJSON, '\n')
	_, err = conn.Write(msgJSON)
	if err != nil {
		log.Printf("[TCP Client] Failed to send broadcast: %v", err)
	}
}
