package enums

import (
	"encoding/json"
	"errors"
)

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

// MarshalJSON converts the CardColor to a JSON string.
func (c CardColor) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

// UnmarshalJSON parses the JSON string into a CardColor.
func (c *CardColor) UnmarshalJSON(data []byte) error {
	var colorStr string
	if err := json.Unmarshal(data, &colorStr); err != nil {
		return err
	}

	// Map string to CardColor
	switch colorStr {
	case "Red":
		*c = Red
	case "Blue":
		*c = Blue
	case "Green":
		*c = Green
	case "Yellow":
		*c = Yellow
	case "Wild":
		*c = Wild
	default:
		return errors.New("invalid card color")
	}

	return nil
}
