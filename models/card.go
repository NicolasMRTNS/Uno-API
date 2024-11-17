package models

import "github.com/NicolasMRTNS/Uno-API/enums"

type Card struct {
	Value enums.CardValue `json:"value"`
	Color enums.CardColor `json:"color"`
}
