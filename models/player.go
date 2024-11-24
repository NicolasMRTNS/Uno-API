package models

type Player struct {
	Id         string `json:"id"`
	Name       string `json:"name"`
	PlayerDeck Deck   `json:"deck"`
	IsInGame   bool   `json:"isInGame"`
}
