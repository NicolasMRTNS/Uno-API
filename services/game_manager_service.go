package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/NicolasMRTNS/Uno-API/enums"
	"github.com/NicolasMRTNS/Uno-API/models"
	"github.com/NicolasMRTNS/Uno-API/utils"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type GameManager struct {
	games   sync.Map // Thread-safe map of gameId -> *Game
	sockets sync.Map // Thread-safe map of gameId -> []*websocket.Conn
}

var (
	GameManagerInstance *GameManager
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
	gameId := c.Param("gameId")
	playerName := c.Param("playerName")
	gm := GetGameManager()

	game, _ := gm.GetGame(gameId)

	newPlayer := CreatePlayer(fullDeck, playerName)

	if err := game.AddPlayer(newPlayer); err != nil {
		fmt.Print(err)
	}

	gm.games.Store(game.Id, game)
	c.JSON(http.StatusCreated, game)
}

func StartGame(c *gin.Context) {
	gameId := c.Param("gameId")
	game, _ := GameManagerInstance.GetGame(gameId)
	if len(game.Players) >= 2 {
		game = game.SetGameStatus(enums.InProgress)
		GameManagerInstance.games.Store(game.Id, game)
		c.JSON(http.StatusOK, game)
	} else {
		c.JSON(http.StatusBadRequest, game)
	}
}

func (gm *GameManager) AddPlayerSocket(gameID, playerID string, conn *websocket.Conn) error {
	// Retrieve the game from the manager
	game, _ := gm.GetGame(gameID)

	// Add the player's WebSocket to the game's PlayerSockets
	game.Mutex.Lock()
	defer game.Mutex.Unlock()

	// Check if the player is part of the game
	playerFound := false
	for _, player := range game.Players {
		if player.Id == playerID {
			playerFound = true
			break
		}
	}

	if !playerFound {
		return fmt.Errorf("player %s is not in the game", playerID)
	}

	// Add the WebSocket connection
	if game.PlayerSockets == nil {
		game.PlayerSockets = make(map[string]*websocket.Conn)
	}

	game.PlayerSockets[playerID] = conn
	return nil
}

// Function to get the GameManager instance and create one if needed (Singleton)
func GetGameManager() *GameManager {
	once.Do(func() {
		GameManagerInstance = &GameManager{}
	})
	return GameManagerInstance
}

func gameLoop(game *Game, actionChan chan models.GameAction, gm *GameManager) {
	for {
		select {
		case action := <-actionChan:
			// Process the action
			handleAction(game, action)

			// Broadcast updated game state to all players
			gameStateJson, _ := json.Marshal(game)
			game.Broadcast(gameStateJson)
		default:
			time.Sleep(1 * time.Second)
		}

		if game.State == enums.Completed || game.State == enums.Cancelled {
			break
		}
	}
}

func handleAction(game *Game, action models.GameAction) {
	player := findPlayerById(game, action.PlayerId)
	if player == nil {
		fmt.Printf("Player %s not found \n", action.PlayerId)
		return
	}

	switch action.Type {
	case enums.ActionPlayCard:
		// Validate the card exists in the player's hand
		card, err := findCardInHand(player, &action.Card)
		if err != nil {
			fmt.Println("Invalid card: ", card)
			return
		}

		// Validate if the card can be played
		if !canPlayCard(game, card) {
			fmt.Println("Card cannot be played: ", card)
			return
		}

		// Play the card
		removeCardFromHand(player, card)
		game.GameDeck = *card

		// Handle special card effects
		handleSpecialCard(game, card)

		fmt.Printf("Player %s played card %s\n", action.PlayerId, card)

	case enums.ActionDrawCard:
		// Get the card drawn
		cardDrawn := game.DrawPile.Cards[0]

		// Remove the card from the pile
		game.DrawPile.Cards = game.DrawPile.Cards[1:]

		// Add the card to the player's deck
		player.PlayerDeck.Cards = append(player.PlayerDeck.Cards, cardDrawn)

		// Refill the pile
		game.DrawPile.Cards = append(game.DrawPile.Cards, utils.ShuffleDeck(fullDeck)[0])

		fmt.Printf("Player %s drew a card: %s\n", action.PlayerId, cardDrawn)
	case enums.ActionEndTurn:
		if len(player.PlayerDeck.Cards) == 0 {
			game.State = enums.Completed
			fmt.Printf("Player %s has won the game!\n", action.PlayerId)
			return
		}

		selectNextPlayer(game)
		fmt.Printf("Turn ended. Next player is: %s\n", game.ActivePlayer)
	}
}

// TODO: Reorganize these next functions
func findPlayerById(game *Game, playerId string) *models.Player {
	for _, player := range game.Players {
		if player.Id == playerId {
			return &player
		}
	}
	return nil
}

func findCardInHand(player *models.Player, cardToPlay *models.Card) (*models.Card, error) {
	for _, card := range player.PlayerDeck.Cards {
		if card.Value == cardToPlay.Value && card.Color == cardToPlay.Color {
			return &card, nil
		}
	}
	return nil, fmt.Errorf("card not found in hand")
}

func removeCardFromHand(player *models.Player, card *models.Card) {
	for i, c := range player.PlayerDeck.Cards {
		if c.Color == card.Color && c.Value == card.Value {
			player.PlayerDeck.Cards = append(player.PlayerDeck.Cards[:i], player.PlayerDeck.Cards[i+1:]...)
			return
		}
	}
}

func canPlayCard(game *Game, card *models.Card) bool {
	return card.Color == game.GameDeck.Color && card.Value == game.GameDeck.Value
}

func handleSpecialCard(game *Game, card *models.Card) {
	switch card.Value {
	case enums.Skip:
		selectNextPlayer(game)
	case enums.Reverse:
		reversePlayDirection(game)
	case enums.DrawTwo:
		// drawCardsForNextPlayer
	case enums.CardValue(enums.Wild):
		// Let the player choose a color
	case enums.WildDrawFour:
		// drawCardsForNextPlayer
		// Let the player choose a color
	}
}

func selectNextPlayer(game *Game) {
	nextPlayer := getNextPlayerWithDirection(game, 1)
	if nextPlayer != nil {
		game.ActivePlayer = *nextPlayer
		fmt.Printf("Player skipped! Next player is: %s\n", game.ActivePlayer)
	}
}

// The step parameter is in case we need to skip more than one player
func getNextPlayerWithDirection(game *Game, step int) *models.Player {
	// Find the current player's index
	currentIndex := -1
	for i, player := range game.Players {
		if player.Id == game.ActivePlayer.Id {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		fmt.Println("Error: ActivePlayer not found in Players list")
		return nil
	}

	// Calculate the next index with the given step
	var nextIndex int
	if game.Reverse {
		nextIndex = (currentIndex - step + len(game.Players)) % len(game.Players)
	} else {
		nextIndex = (currentIndex + step) % len(game.Players)
	}

	// Return the next player
	return &game.Players[nextIndex]
}

func reversePlayDirection(game *Game) {
	game.Reverse = !game.Reverse
}
