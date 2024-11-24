package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/NicolasMRTNS/Uno-API/enums"
	"github.com/NicolasMRTNS/Uno-API/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type GameManager struct {
	games   sync.Map // Thread-safe map of gameId -> *Game
	sockets sync.Map // Thread-safe map of gameId -> []*websocket.Conn
}

var (
	gameManagerInstance *GameManager
	once                sync.Once
)

func (gm *GameManager) StartGame(gameId string) {
	game, _ := gameManager.GetGame(gameId)
	game.State = enums.InProgress
	gm.games.Store(game.Id, game)

	actionChan := make(chan models.GameAction)

	gm.games.Store(game.Id+"_action", actionChan)

	go func() {
		defer gm.games.Delete(game.Id) // Cleanup on game end
		defer gm.games.Delete(game.Id + "_action")
		defer gm.sockets.Delete(game.Id)
		gameLoop(game, actionChan, gm)
	}()
}

func (gm *GameManager) SendAction(gameId string, action models.GameAction) error {
	value, ok := gm.games.Load(gameId)

	// TODO: handle error
	if !ok {
		return fmt.Errorf("game not found")
	}

	actionChan := value.(chan models.GameAction)
	actionChan <- action
	return nil
}

func (gm *GameManager) StopGame(gameId string) {
	value, ok := gm.games.Load(gameId)
	if ok {
		actionChan := value.(chan models.GameAction)
		close(actionChan)
	}
	gm.games.Delete(gameId)
}

func (gm *GameManager) AddGameToGameManager(game *Game) error {
	if _, loaded := gm.games.LoadOrStore(game.Id, game); loaded {
		return fmt.Errorf("game with ID %s already exists", game.Id)
	}
	return nil
}

func (gm *GameManager) GameExists(gameId string) bool {
	_, exists := gm.games.Load(gameId)
	return exists
}

func (gm *GameManager) GetGame(gameId string) (*Game, error) {
	value, exists := gm.games.Load(gameId)
	if !exists {
		return nil, fmt.Errorf("game with Id %s not found", gameId)
	}
	return value.(*Game), nil
}

func (gm *GameManager) GetSocketConnection(gameId string) (*sync.Map, error) {
	value, exists := gm.sockets.Load(gameId)
	if !exists {
		return nil, fmt.Errorf("socket connection for game ID %s not found", gameId)
	}
	return value.(*sync.Map), nil
}

func AddPlayerToGame(c *gin.Context) {
	gameId := c.Param("game_id")
	playerName := c.Param("player_name")
	gm := gameManagerInstance

	game, _ := gm.GetGame(gameId)

	newPlayer := CreatePlayer(fullDeck, playerName)

	if err := game.AddPlayer(newPlayer); err != nil {
		fmt.Print(err)
	}

	gm.games.Store(game.Id, game)
	c.JSON(http.StatusCreated, game)
}

func (gm *GameManager) AddPlayerSocket(gameId string, conn *websocket.Conn) {
	value, _ := gm.sockets.LoadOrStore(gameId, &sync.Map{})
	connections := value.(*sync.Map)
	connections.Store(conn, true)
}

func (gm *GameManager) BroadcastToGame(gameId string, message []byte) {
	connections, _ := gm.GetSocketConnection(gameId)
	connections.Range(func(key, _ interface{}) bool {
		conn := key.(*websocket.Conn)
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			conn.Close()
			connections.Delete(conn)
		}
		return true
	})
}

// Function to get the GameManager instance and create one if needed (Singleton)
func GetGameManager() *GameManager {
	once.Do(func() {
		gameManagerInstance = &GameManager{}
	})
	return gameManagerInstance
}

func gameLoop(game *Game, actionChan chan models.GameAction, gm *GameManager) {
	for {
		select {
		case action := <-actionChan:
			// Process the action
			handleAction(game, action)

			// Broadcast updated game state to all players
			gameStateJson, _ := json.Marshal(game)
			gm.BroadcastToGame(game.Id, gameStateJson)
		default:
			time.Sleep(1 * time.Second)
		}

		if game.State == enums.Completed || game.State == enums.Cancelled {
			break
		}
	}
}

func handleAction(game *Game, action models.GameAction) {
	switch action.Type {
	case enums.ActionPlayCard:
		// Handle play card logic
	case enums.ActionDrawCard:
		// Handle card draw logic
	case enums.ActionEndTurn:
		// Handle end turn logic
	}
}
