package utils

import (
	"math/rand"

	"github.com/NicolasMRTNS/Uno-API/enums"
	"github.com/NicolasMRTNS/Uno-API/models"
)

// Get all combinations of colors and values
func GenerateFullDeck() []models.Card {
	fullDeck := []models.Card{}
	for _, color := range []enums.CardColor{enums.Red, enums.Blue, enums.Green, enums.Yellow} {
		for _, value := range []enums.CardValue{
			enums.Zero, enums.One, enums.Two, enums.Three, enums.Four,
			enums.Five, enums.Six, enums.Seven, enums.Eight, enums.Nine,
			enums.Skip, enums.Reverse, enums.DrawTwo} {
			fullDeck = append(fullDeck, models.Card{Color: color, Value: value})
		}
	}

	// Add Wild Cards (no specific color)
	for _, value := range []enums.CardValue{enums.WildCard, enums.WildDrawFour} {
		fullDeck = append(fullDeck, models.Card{Color: enums.Wild, Value: value})
	}

	return fullDeck
}

func ShuffleDeck(deck []models.Card) []models.Card {
	shuffled := make([]models.Card, len(deck))
	perm := rand.Perm(len(deck))
	for i, v := range perm {
		shuffled[i] = deck[v]
	}
	return shuffled
}
