package websocket

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"mangahub/pkg/models"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// HandleConnection handles a new WebSocket connection from Gin
func HandleConnection(c *gin.Context, hub *Hub, room string) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// Create chat client
	client := &models.ChatClient{
		UserID:   c.GetString("user_id"),
		Username: c.GetString("username"),
		RoomID:   room,
	}

	// Register with hub
	hub.Register <- client

	// Handle connection
	hub.HandleConnection(conn, client)
}
