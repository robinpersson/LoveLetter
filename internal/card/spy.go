package card

import "fmt"

const (
	spyName        = "Spy"
	spyDescription = "If you are the only one still in play who has played or discarded this card by the end of the round, you gain an favor token"
	spyValue       = 0
)

type Spy interface {
	Card
}

type spy struct {
	CardSettings
}

func (c spy) ToString() string {
	return fmt.Sprintf("[%d] %s", c.value, c.name)
}

func (c spy) ShortString() string {
	return fmt.Sprintf("[%d]", c.value)
}

func (c spy) Name() string {
	return c.name
}

func (c spy) Value() int {
	return c.value
}

func (c spy) ActionText() string {
	return ""
}

func NewSpy() Spy {
	return &spy{CardSettings{
		description: spyDescription,
		value:       spyValue,
		name:        spyName,
	}}
}

func GetSpyRule() string {
	return fmt.Sprintf("\u001B[32;1m[%d]\u001B[0m%s(x%d) - %s\n\n", spyValue, spyName, spyCount, spyDescription)
}
