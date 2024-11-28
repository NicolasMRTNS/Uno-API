package models

import (
	"github.com/NicolasMRTNS/Uno-API/enums"
)

type GameAction struct {
	Type     enums.GameActionType `json:"type"`
	PlayerId string               `json:"player_id"`
	Card     Card                 `json:"card"`
}
