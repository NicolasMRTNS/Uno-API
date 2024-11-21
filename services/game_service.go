package services

import (
	"fmt"

	"github.com/NicolasMRTNS/Uno-API/enums"
	"github.com/NicolasMRTNS/Uno-API/models"
	"github.com/NicolasMRTNS/Uno-API/utils"
	"github.com/google/uuid"
)

type Game struct {
	Id           string          `json:"id"`
	Players      []models.Player `json:"players"`
	GameDeck     models.Card     `json:"gameDeck"`
	DrawPile     models.Deck     `json:"drawPile"`
	State        enums.GameState `json:"state"`
	ActivePlayer models.Player   `json:"activePlayer"`
}

var (
	shuffledFullDeck = utils.ShuffleDeck(utils.GenerateFullDeck())
	gameManager      = GetGameManager()
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

func CreateNewGame(playerName string) {
	// Get 21 cards: 1 card for the main game deck and 20 for the draw pile
	startingDeckAndDrawPile := shuffledFullDeck[:21]

	startingDrawPile := models.Deck{
		Cards:        startingDeckAndDrawPile[1:],
		IsPlayerDeck: false,
	}

	currentPlayer := CreatePlayer(shuffledFullDeck, playerName)

	game := &Game{
		Id:           uuid.NewString(),
		Players:      []models.Player{currentPlayer},
		GameDeck:     startingDeckAndDrawPile[0],
		DrawPile:     startingDrawPile,
		State:        enums.WaitingForPlayers,
		ActivePlayer: currentPlayer,
	}

	gameManager.AddGameToGameManager(game)
}
