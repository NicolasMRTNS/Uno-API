package services

import (
	"fmt"
	"sync"
	"time"

	"github.com/NicolasMRTNS/Uno-API/enums"
	"github.com/NicolasMRTNS/Uno-API/models"
)

type GameManager struct {
	games sync.Map // Thread-safe map to store active games
}

var (
	gameManagerInstance *GameManager
	once                sync.Once
)

func (gm *GameManager) StartGame(game *Game) {
	actionChan := make(chan models.GameAction)
	gameId := game.Id

	gm.games.Store(gameId, actionChan)

	go func() {
		defer gm.games.Delete(gameId) // Cleanup on game end
		gameLoop(game, actionChan)
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

func (gm *GameManager) AddPlayerToGame(gameId, playerName string) error {
	game, error := gm.GetGame(gameId)
	if error != nil {
		fmt.Errorf("game not found")
	}

	newPlayer := CreatePlayer(shuffledFullDeck, playerName)

	if err := game.AddPlayer(newPlayer); err != nil {
		return err
	}

	gm.games.Store(game.Id, game)
	return nil
}

// Function to get the GameManager instance and create one if needed (Singleton)
func GetGameManager() *GameManager {
	once.Do(func() {
		gameManagerInstance = &GameManager{}
	})
	return gameManagerInstance
}

func gameLoop(game *Game, actions chan models.GameAction) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case action, ok := <-actions:
			if !ok {
				fmt.Printf("Game %s ended\n", game.Id)
				return
			}

			handleAction(game, action)

		case <-ticker.C:
			// Check for inactivity
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
