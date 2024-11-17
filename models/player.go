package models

type Player struct {
	Id         string `json:"id"`
	Name       string `json:"players"`
	PlayerDeck Deck   `json:"deck"`
	IsInGame   bool   `json:"isInGame"`
}
