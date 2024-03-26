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

func GetCardByValue(value int) Card {
	switch value {
	case 0:
		return NewSpy()
	case 1:
		return NewGuard()
	case 2:
		return NewPriest()
	case 3:
		return NewBaron()
	case 4:
		return NewHandmaid()
	case 5:
		return NewPrince()
	case 6:
		return NewChancellor()
	case 7:
		return NewKing()
	case 8:
		return NewCountess()
	case 9:
		return NewPrincess()
	default:
		return nil
	}
}

//func NewCard() Card {
//	return &card{}
//}

//func (c card) GetCards() []card {
//	//TODO implement me
//	panic("implement me")
//}
