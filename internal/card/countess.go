package card

import "fmt"

const (
	countessName        = "Countess"
	countessDescription = "Must be played if you have King or Prince"
	countessValue       = 8
)

type Countess interface {
	Card
}

type countess struct {
	CardSettings
}

func (c countess) SetIndex(index int) {
	c.index = index
}

func (c countess) Index() int {
	return c.index
}

func (c countess) ToString() string {
	return fmt.Sprintf("[%d]%s", c.value, c.name)
}

func (c countess) ShortString() string {
	return fmt.Sprintf("[%d]", c.value)
}

func (c countess) Name() string {
	return c.name
}

func (c countess) Value() int {
	return c.value
}

func (c countess) ActionText() string {
	return ""
}

func NewCountess() Countess {
	return &countess{CardSettings{
		description: countessDescription,
		value:       countessValue,
		name:        countessName,
	}}
}

func GetCountessRule() string {
	return fmt.Sprintf("\u001B[32;1m[%d]\u001B[0m%s(x%d) - %s\n\n", countessValue, countessName, countessCount, countessDescription)
}
