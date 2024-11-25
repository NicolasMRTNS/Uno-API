package main

import (
	"fmt"
	"net/http"

	"github.com/NicolasMRTNS/Uno-API/services"
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

	router.POST("/create-game/:playerName", services.CreateNewGame)
	router.POST("/add-player/:gameId/:playerName", services.AddPlayerToGame)
	router.POST("start-game/:gameId", services.StartGame)

	router.GET("/ws", handlerWebSocket)

	router.Run(":8080")
}

// Handles incoming WebSocket connections
func handlerWebSocket(c *gin.Context) {
	playerId := c.Param("playerId")
	gameId := c.Param("gameId")

	// Upgrade the HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("WebSocket upgrade error: ", err)
		return
	}
	defer conn.Close()
	gameManager := services.GameManagerInstance

	if err := gameManager.AddPlayerSocket(gameId, playerId, conn); err != nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Error: "+err.Error()))
		return
	}

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
