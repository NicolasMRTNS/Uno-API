package services

import (
	"github.com/NicolasMRTNS/Uno-API/enums"
	"github.com/NicolasMRTNS/Uno-API/models"
	"github.com/NicolasMRTNS/Uno-API/utils"
	"github.com/google/uuid"
)

func CreateNewGame(playerName string) *models.Game {
	// Shuffle the entire deck
	shuffledDeck := utils.ShuffleDeck(utils.GenerateFullDeck())

	// Get 21 cards: 1 card for the main game deck and 20 for the draw pile
	startingDeckAndDrawPile := shuffledDeck[:21]

	startingDrawPile := models.Deck{
		Cards:        startingDeckAndDrawPile[1:],
		IsPlayerDeck: false,
	}

	return &models.Game{
		Id:             uuid.NewString(),
		PlayersIds:     []string{playerName},
		GameDeck:       startingDeckAndDrawPile[0],
		DrawPile:       startingDrawPile,
		State:          enums.WaitingForPlayers,
		ActivePlayerId: playerName,
	}
}
