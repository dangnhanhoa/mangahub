package websocket

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

func (h *Handler) ServeWS(c *gin.Context) {
	mangaID := c.Param("mangaId")
	
	userID := c.MustGet("user_id").(string)
	username := "User_" + userID[:5] 

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("[WS] Error upgrading connection: %v", err)
		return 
	}

	client := &Client{
		Conn:     conn,	
		UserID:   userID,
		Username: username,
	}

	h.hub.Register(mangaID, client)

	
	defer func() {
		h.hub.Unregister(mangaID, client)
	}()

	for {
		var msg ChatMessage
		
		err := conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WS] Error reading message from %s: %v", username, err)
			}
			break 
		}

		msg.UserID = client.UserID
		msg.Username = client.Username
		msg.Timestamp = time.Now().Unix()

		h.hub.Broadcast(mangaID, msg)
	}
}		