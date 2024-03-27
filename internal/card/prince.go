package card

import "fmt"

const (
	princeName        = "Prince"
	princeDescription = "Choose any player (including yourself) to discard their hand and draw a new one"
	princeValue       = 5
)

type Prince interface {
	Card
}

type prince struct {
	CardSettings
}

func (c prince) ToString() string {
	return fmt.Sprintf("[%d]%s", c.value, c.name)
}

func (c prince) ShortString() string {
	return fmt.Sprintf("[%d]", c.value)
}

func (c prince) Name() string {
	return c.name
}

func (c prince) Value() int {
	return c.value
}

func (c prince) ActionText() string {
	return "Write the player name that should discard his card"
}

func NewPrince() Prince {
	return &prince{CardSettings{
		description: princeDescription,
		value:       princeValue,
		name:        princeName,
	}}
}

func GetPrinceRule() string {
	return fmt.Sprintf("\u001B[32;1m[%d]\u001B[0m%s(x%d) - %s\n\n", princeValue, princeName, princeCount, princeDescription)
}
