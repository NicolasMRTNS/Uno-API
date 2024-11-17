package enums

type CardValue int

const (
	Zero CardValue = iota
	One
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Skip
	Reverse
	DrawTwo
	WildCard
	WildDrawFour
)

func (v CardValue) String() string {
	return [...]string{"Zero", "One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine",
		"Skip", "Reverse", "Draw Two", "Wild Card", "Wild Draw Four"}[v]
}
