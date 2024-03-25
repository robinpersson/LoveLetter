package game

import (
	"github.com/robinpersson/LoveLetter/internal/card"
)

type PlayerChoices interface {
	PickCard() *card.Card
	DealCard()
	LayCard(card card.Card)
}

type Player struct {
	CurrentCard *card.Card
	LaidCards   []card.Card
	PickedCard  *card.Card
	Game        *Game
}

func NewPlayer(game Game) *Player {
	return &Player{
		CurrentCard: nil,
		LaidCards:   nil,
		PickedCard:  nil,
		Game:        &game,
	}
}

func (p *Player) PickCard() *card.Card {
	c := p.Game.PickCard()
	p.PickedCard = c
	return p.PickedCard
}

func (p *Player) DealCard() {
	p.CurrentCard = p.Game.PickCard()
}

func (p *Player) LayCard(c card.Card) {
	p.PickedCard = nil
	var cards []card.Card
	cards = append(cards, c)
	p.LaidCards = append(cards, p.LaidCards...)
}
