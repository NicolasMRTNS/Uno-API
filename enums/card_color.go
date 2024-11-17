package enums

type CardColor int

const (
	Red CardColor = iota
	Blue
	Green
	Yellow
	Wild
)

func (c CardColor) String() string {
	return [...]string{"Red", "Blue", "Green", "Yellow", "Wild"}[c]
}
