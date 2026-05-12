package udp

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

func (c *Client) TriggerBroadcast(mangaID, message string) {
	conn, err := net.Dial("udp", c.addr)
	if err != nil {
		log.Printf("[UDP Client] Failed to dial %s: %v", c.addr, err)
		return
	}
	defer conn.Close()

	notif := Notification{
		Type:    "BROADCAST",
		MangaID: mangaID,
		Message: message,
	}

	payload, err := json.Marshal(notif)
	if err != nil {
		log.Printf("[UDP Client] Failed to marshal payload: %v", err)
		return
	}

	_, err = conn.Write(payload)
	if err != nil {
		log.Printf("[UDP Client] Failed to send broadcast: %v", err)
	}
}
