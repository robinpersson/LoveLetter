package card

import "fmt"

const (
	kingName        = "King"
	kingDescription = "Trade hands with another player"
	kingValue       = 7
)

type King interface {
	Card
}

type king struct {
	CardSettings
}

func (c king) ToString() string {
	return fmt.Sprintf("[%d]%s", c.value, c.name)
}

func (c king) ShortString() string {
	return fmt.Sprintf("[%d]", c.value)
}

func (c king) Name() string {
	return c.name
}

func (c king) Value() int {
	return c.value
}

func (c king) ActionText() string {
	return "Write the name of the player which you would like to trade hands with"
}

func NewKing() King {
	return &king{CardSettings{
		description: kingDescription,
		value:       kingValue,
		name:        kingName,
	}}
}

func GetkingRule() string {
	return fmt.Sprintf("\u001B[32;1m[%d]\u001B[0m%s(x%d) - %s\n\n", kingValue, kingName, kingCount, kingDescription)
}
