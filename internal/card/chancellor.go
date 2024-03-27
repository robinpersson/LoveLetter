package card

import "fmt"

const (
	chancellorName        = "Chancellor"
	chancellorDescription = "Draw two cards from the deck, then choose one of the three cards to keep, and place the other two at the bottom of the deck in any order"
	chancellorValue       = 6
)

type Chancellor interface {
	Card
}

type chancellor struct {
	CardSettings
}

func (c chancellor) ToString() string {
	return fmt.Sprintf("[%d]%s", c.value, c.name)
}

func (c chancellor) ShortString() string {
	return fmt.Sprintf("[%d]", c.value)
}

func (c chancellor) Name() string {
	return c.name
}

func (c chancellor) Value() int {
	return c.value
}

func (c chancellor) ActionText() string {
	return "hard one"
}

func NewChancellor() Chancellor {
	return &chancellor{CardSettings{
		description: chancellorDescription,
		value:       chancellorValue,
		name:        chancellorName,
	}}
}

func GetChancellorRule() string {
	return fmt.Sprintf("\u001B[32;1m[%d]\u001B[0m%s(x%d) - %s\n\n", chancellorValue, chancellorName, chancellorCount, chancellorDescription)
}
