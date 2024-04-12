package card

import "fmt"

const (
	priestName        = "Priest"
	priestDescription = "You may privately see an opponents hand"
	priestValue       = 2
)

type Priest interface {
	Card
}

type priest struct {
	CardSettings
}

func (c priest) SetIndex(index int) {
	c.index = index
}

func (c priest) Index() int {
	return c.index
}

func (c priest) ToString() string {
	return fmt.Sprintf("[%d]%s", c.value, c.name)
}

func (c priest) ShortString() string {
	return fmt.Sprintf("[%d]", c.value)
}

func (c priest) Name() string {
	return c.name
}

func (c priest) Value() int {
	return c.value
}

func (c priest) ActionText() string {
	return "Write the name of the player which card you want to see"
}

func NewPriest() Priest {
	return &priest{CardSettings{
		description: priestDescription,
		value:       priestValue,
		name:        priestName,
	}}
}

func GetPriestRule() string {
	return fmt.Sprintf("\u001B[32;1m[%d]\u001B[0m%s(x%d) - %s\n\n", priestValue, priestName, priestCount, priestDescription)
}
