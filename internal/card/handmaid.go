package card

import "fmt"

const (
	handmaidName        = "Handmaid"
	handmaidDescription = "You cannot be affected by any other player's cards until your next turn"
	handmaidValue       = 4
)

type Handmaid interface {
	Card
}

type handmaid struct {
	CardSettings
}

func (c handmaid) SetIndex(index int) {
	c.index = index
}

func (c handmaid) Index() int {
	return c.index
}

func (c handmaid) ToString() string {
	return fmt.Sprintf("[%d]%s", c.value, c.name)
}

func (c handmaid) ShortString() string {
	return fmt.Sprintf("[%d]", c.value)
}

func (c handmaid) Name() string {
	return c.name
}

func (c handmaid) Value() int {
	return c.value
}

func (c handmaid) ActionText() string {
	return ""
}

func NewHandmaid() Handmaid {
	return &handmaid{CardSettings{
		description: handmaidDescription,
		value:       handmaidValue,
		name:        handmaidName,
	}}
}

func GetHandmaidRule() string {
	return fmt.Sprintf("\u001B[32;1m[%d]\u001B[0m%s(x%d) - %s\n\n", handmaidValue, handmaidName, handmaidCount, handmaidDescription)
}
