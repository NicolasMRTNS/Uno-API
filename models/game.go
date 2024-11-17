package models

import "github.com/NicolasMRTNS/Uno-API/enums"

type Game struct {
	Id             string          `json:"id"`
	PlayersIds     []string        `json:"players"`
	GameDeck       Card            `json:"gameDeck"`
	DrawPile       Deck            `json:"drawPile"`
	State          enums.GameState `json:"state"`
	ActivePlayerId string          `json:"activePlayer"`
}
