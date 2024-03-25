package card

import "fmt"

const (
	baronName        = "Baron"
	baronDescription = "Privately compare hand with opponent. Player with lowest are eliminated"
	baronValue       = 3
)

type Baron interface {
	Card
}

type baron struct {
	CardSettings
}

func (c baron) ToString() string {
	return fmt.Sprintf("[%d] %s", c.value, c.name)
}

func (c baron) ShortString() string {
	return fmt.Sprintf("[%d]", c.value)
}

func (c baron) Name() string {
	return c.name
}

func (c baron) Value() int {
	return c.value
}

func (c baron) ActionText() string {
	return "Write the name of the player you would like to compare hands with"
}
func NewBaron() Baron {
	return &baron{CardSettings{
		description: baronDescription,
		value:       baronValue,
		name:        baronName,
	}}
}

func GetBaronRule() string {
	return fmt.Sprintf("\u001B[32;1m[%d]\u001B[0m%s(x%d) - %s\n\n", baronValue, baronName, baronCount, baronDescription)
}
