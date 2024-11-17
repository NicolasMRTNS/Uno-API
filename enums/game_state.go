package enums

type GameState string

const (
	WaitingForPlayers GameState = "waiting_for_players"
	InProgress        GameState = "in_progress"
	Completed         GameState = "completed"
	Cancelled         GameState = "cancelled"
)
