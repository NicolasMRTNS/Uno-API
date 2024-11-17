package models

type Deck struct {
	Cards        []Card `json:"cards"`
	IsPlayerDeck bool   `json:"isPlayerDeck"`
}
