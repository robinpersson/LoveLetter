package card

import "fmt"

const (
	princessName        = "Princess"
	princessDescription = "If you discard this card you are eliminated"
	princessValue       = 9
)

type Princess interface {
	Card
}

type princess struct {
	CardSettings
}

func (c princess) ToString() string {
	return fmt.Sprintf("[%d]%s", c.value, c.name)
}

func (c princess) ShortString() string {
	return fmt.Sprintf("[%d]", c.value)
}

func (c princess) Name() string {
	return c.name
}

func (c princess) Value() int {
	return c.value
}

func (c princess) ActionText() string {
	return ""
}

func NewPrincess() Princess {
	return &princess{CardSettings{
		description: princessDescription,
		value:       princessValue,
		name:        princessName,
	}}
}

func GetPrincessRule() string {
	return fmt.Sprintf("\u001B[32;1m[%d]\u001B[0m%s(x%d) - %s\n\n", princessValue, princessName, princessCount, princessDescription)
}

func GetPrincessKey() string {
	return fmt.Sprintf("F9 Princess")
}
