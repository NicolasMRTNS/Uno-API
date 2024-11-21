package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// WebSocket upgrader to upgrade HTTP connections to WebSocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for simplicity for now
	},
}

var connections = make(map[*websocket.Conn]bool)

func main() {
	router := gin.Default()

	router.GET("/ws", HandlerWebSocket)

	router.Run(":8080")
}

// Handles incoming WebSocket connections
func HandlerWebSocket(c *gin.Context) {
	// Upgrade the HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error: ", err)
		return
	}
	defer conn.Close()

	// Register the connection
	connections[conn] = true
	defer delete(connections, conn)

	for {
		// Read message from the client
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("WebSocket read error: ", err)
			break
		}
		fmt.Printf("Received message: %s\n", message)

		// Broadcast the message to all connected clients
		broadcastMessage(message)
	}
}

func broadcastMessage(message []byte) {
	for conn := range connections {
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Println("WebSocket write error: ", err)
			conn.Close()
			delete(connections, conn)
		}
	}
}
