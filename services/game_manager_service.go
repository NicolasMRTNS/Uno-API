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

func (gm *GameManager) StartGame(game *models.Game) {
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

func NewGameManager() *GameManager {
	return &GameManager{}
}

func gameLoop(game *models.Game, actions chan models.GameAction) {
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

func handleAction(game *models.Game, action models.GameAction) {
	switch action.Type {
	case enums.ActionPlayCard:
		// Handle play card logic
	case enums.ActionDrawCard:
		// Handle card draw logic
	case enums.ActionEndTurn:
		// Handle end turn logic
	}
}
