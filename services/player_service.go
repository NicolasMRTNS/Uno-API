package services

import (
	"github.com/NicolasMRTNS/Uno-API/models"
	"github.com/google/uuid"
)

func CreatePlayer(fulldeck []models.Card, playerName string) models.Player {
	initialDeck := models.Deck{
		Cards:        fulldeck[:5],
		IsPlayerDeck: true,
	}
	return models.Player{
		Id:         uuid.NewString(),
		Name:       playerName,
		PlayerDeck: initialDeck,
		IsInGame:   true,
	}
}
