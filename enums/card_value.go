package enums

import (
	"encoding/json"
	"errors"
)

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

// MarshalJSON converts CardValue to a JSON string.
func (v CardValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

// UnmarshalJSON parses a JSON string into a CardValue.
func (v *CardValue) UnmarshalJSON(data []byte) error {
	var valueStr string
	if err := json.Unmarshal(data, &valueStr); err != nil {
		return err
	}

	// Map string to CardValue
	switch valueStr {
	case "Zero":
		*v = Zero
	case "One":
		*v = One
	case "Two":
		*v = Two
	case "Three":
		*v = Three
	case "Four":
		*v = Four
	case "Five":
		*v = Five
	case "Six":
		*v = Six
	case "Seven":
		*v = Seven
	case "Eight":
		*v = Eight
	case "Nine":
		*v = Nine
	case "Skip":
		*v = Skip
	case "Reverse":
		*v = Reverse
	case "Draw Two":
		*v = DrawTwo
	case "Wild Card":
		*v = WildCard
	case "Wild Draw Four":
		*v = WildDrawFour
	default:
		return errors.New("invalid card value")
	}

	return nil
}
