package card

const (
	spyCount        = 2
	guardCount      = 6
	priestCount     = 2
	baronCount      = 2
	handmaidCount   = 2
	princeCount     = 2
	chancellorCount = 2
	kingCount       = 1
	countessCount   = 1
	princessCount   = 1
)

type Card interface {
	Name() string
	Value() int
	ToString() string
	ShortString() string
	ActionText() string
}

type CardSettings struct {
	description string
	value       int
	name        string
}

//func NewCard() Card {
//	return &card{}
//}

//func (c card) GetCards() []card {
//	//TODO implement me
//	panic("implement me")
//}
