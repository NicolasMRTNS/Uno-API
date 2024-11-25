package services

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/NicolasMRTNS/Uno-API/enums"
	"github.com/NicolasMRTNS/Uno-API/models"
	"github.com/NicolasMRTNS/Uno-API/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Game struct {
	Id            string          `json:"id"`
	Players       []models.Player `json:"players"`
	GameDeck      models.Card     `json:"gameDeck"`
	DrawPile      models.Deck     `json:"drawPile"`
	State         enums.GameState `json:"state"`
	ActivePlayer  string          `json:"activePlayer"`
	PlayerSockets map[string]*websocket.Conn
	Mutex         sync.Mutex
}

var (
	fullDeck    = utils.GenerateFullDeck()
	gameManager = GetGameManager()
)

func (g *Game) AddPlayer(player models.Player) error {
	for _, p := range g.Players {
		if p.Id == player.Id {
			return fmt.Errorf("player with ID %s already exists", player.Id)
		}
	}

	if g.State != enums.WaitingForPlayers {
		return fmt.Errorf("cannot add player as the game has already started")
	}

	g.Players = append(g.Players, player)
	return nil
}

// AddPlayerSocket adds a player's WebSocket connection to the game
func (g *Game) AddPlayerSocket(playerID string, conn *websocket.Conn) error {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()

	// Check if the player is part of the game
	playerFound := false
	for _, player := range g.Players {
		if player.Id == playerID {
			playerFound = true
			break
		}
	}
	if !playerFound {
		return fmt.Errorf("player %s not found in the game", playerID)
	}

	// Store the WebSocket connection
	g.PlayerSockets[playerID] = conn
	return nil
}

func (g *Game) Broadcast(message []byte) {
	g.Mutex.Lock()
	defer g.Mutex.Unlock()

	for playerId, conn := range g.PlayerSockets {
		if conn != nil {
			err := conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				fmt.Printf("Failed to send message to player %s: %v\n", playerId, err)
				conn.Close()
				delete(g.PlayerSockets, playerId) // Cleanup broken connections
			}
		}
	}
}

func (g *Game) SetGameStatus(status enums.GameState) *Game {
	g.State = status
	return g
}

func CreateNewGame(c *gin.Context) {
	shuffledFullDeck := utils.ShuffleDeck(fullDeck)
	// Get 21 cards: 1 card for the main game deck and 20 for the draw pile
	startingDeckAndDrawPile := shuffledFullDeck[:21]

	startingDrawPile := models.Deck{
		Cards:        startingDeckAndDrawPile[1:],
		IsPlayerDeck: false,
	}

	playerName := c.Param("playerName")
	println(c.Request.URL.String())

	currentPlayer := CreatePlayer(shuffledFullDeck, playerName)

	game := &Game{
		Id:           uuid.NewString(),
		Players:      []models.Player{currentPlayer},
		GameDeck:     startingDeckAndDrawPile[0],
		DrawPile:     startingDrawPile,
		State:        enums.WaitingForPlayers,
		ActivePlayer: currentPlayer.Id,
	}

	gameManager.AddGameToGameManager(game)
	c.JSON(http.StatusOK, game)
}
