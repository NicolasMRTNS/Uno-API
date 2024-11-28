package enums

type GameActionType string

const (
	ActionPlayCard GameActionType = "play_card"
	ActionDrawCard GameActionType = "draw_card"
	ActionEndTurn  GameActionType = "end_turn"
)
