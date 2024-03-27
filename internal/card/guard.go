package card

import (
	"fmt"
)

const (
	guardName        = "Guard"
	guardDescription = "Guess an oppenent hand (other than a Guard)"
	guardValue       = 1
)

type Guard interface {
	Card
}

type guard struct {
	CardSettings
}

func (c guard) ToString() string {
	return fmt.Sprintf("[%d]%s", c.value, c.name)
}

func (c guard) ShortString() string {
	return fmt.Sprintf("[%d]", c.value)
}

func (c guard) Name() string {
	return c.name
}

func (c guard) Value() int {
	return c.value
}

func (c guard) ActionText() string {
	return "Write <player name-card number>, eg. 'player1-5' for guessing that Player1 has [5]Prince"
}

func NewGuard() Guard {
	return &guard{CardSettings{
		description: guardDescription,
		value:       guardValue,
		name:        guardName,
	}}
}

func GetGuardRule() string {
	return fmt.Sprintf("\u001B[32;1m[%d]\u001B[0m%s(x%d) - %s\n\n", guardValue, guardName, guardCount, guardDescription)
}
