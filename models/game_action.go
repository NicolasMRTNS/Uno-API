package models

import "github.com/NicolasMRTNS/Uno-API/enums"

type GameAction struct {
	Type    enums.GameActionType
	Payload interface{}
}
