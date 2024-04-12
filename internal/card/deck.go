// -------------------------------------------------------------------
// Copyright (c) Axis Communications AB, SWEDEN. All rights reserved.
// -------------------------------------------------------------------

package card

import (
	"math/rand"
	"time"
)

type Deck interface {
	Shuffle(nrOfPlayers int)
	Cards() []Card
	FaceDownCard() Card
	FaceUpCards() []Card
	PickCard() *Card
	InsertCards(cardsToInsert []Card)
}

type deck struct {
	cards        []Card
	faceDownCard Card
	faceUpCards  []Card
}

func (d *deck) Cards() []Card {
	return d.cards
}

func (d *deck) FaceDownCard() Card {
	return d.faceDownCard
}

func (d *deck) FaceUpCards() []Card {
	return d.faceUpCards
}

func (d *deck) PickCard() *Card {

	if len(d.cards) == 0 {
		return nil
	}

	card := d.cards[0]
	d.cards = d.cards[1:]
	return &card
}

func (d *deck) InsertCards(cardsToInsert []Card) {
	d.cards = append(d.cards, cardsToInsert...)
}

func NewDeck(nrOfPlayers int) Deck {
	d := deck{}
	//d.fakeInit()
	d.init()
	d.Shuffle(nrOfPlayers)
	return &d
}

func (d *deck) Shuffle(nrOfPlayers int) {
	rand.NewSource(time.Now().UnixNano())
	rand.Shuffle(len(d.cards), func(i, j int) { d.cards[i], d.cards[j] = d.cards[j], d.cards[i] })
	d.faceDownCard = d.cards[0]

	if nrOfPlayers == 2 {
		d.faceUpCards = d.cards[3:]
	}

	d.cards = d.cards[1:]

	for i, c := range d.cards {
		c.SetIndex(i)
	}
}

func (d *deck) fakeInit() {
	var cards []Card

	cards = append(cards, NewChancellor())

	cards = append(cards, NewGuard())
	cards = append(cards, NewPrince())
	//cards = append(cards, NewPrincess())
	//cards = append(cards, NewCountess())
	//cards = append(cards, NewCountess())
	cards = append(cards, NewChancellor())
	d.faceDownCard = NewChancellor()
	cards = append(cards, NewKing())
	//cards = append(cards, NewHandmaid())
	cards = append(cards, NewPriest())
	cards = append(cards, NewSpy())
	cards = append(cards, NewCountess())
	//cards = append(cards, NewPrincess())
	//cards = append(cards, NewSpy())

	//cards = append(cards, NewHandmaid())
	//for i := 0; i < princeCount; i++ {
	//	cards = append(cards, NewPrince())
	//}
	d.cards = cards

	for i, _ := range d.cards {
		d.cards[i].SetIndex(i)
	}
	return

	for i := 0; i < baronCount-1; i++ {
		cards = append(cards, NewBaron())
	}

	cards = append(cards, NewSpy())

	for i := 0; i < handmaidCount; i++ {
		cards = append(cards, NewHandmaid())
	}

	for i := 0; i < guardCount; i++ {
		cards = append(cards, NewGuard())
	}
	for i := 0; i < priestCount; i++ {
		cards = append(cards, NewPriest())
	}

	for i := 0; i < spyCount; i++ {
		cards = append(cards, NewSpy())
	}

	for i := 0; i < chancellorCount; i++ {
		cards = append(cards, NewChancellor())
	}

	for i := 0; i < kingCount; i++ {
		cards = append(cards, NewKing())
	}

	for i := 0; i < countessCount; i++ {
		cards = append(cards, NewCountess())
	}

	for i := 0; i < princessCount; i++ {
		cards = append(cards, NewPrincess())
	}

	d.cards = cards
}

func (d *deck) init() {
	var cards []Card

	for i := 0; i < spyCount; i++ {
		cards = append(cards, NewSpy())
	}

	for i := 0; i < guardCount; i++ {
		cards = append(cards, NewGuard())
	}

	for i := 0; i < priestCount; i++ {
		cards = append(cards, NewPriest())
	}

	for i := 0; i < baronCount; i++ {
		cards = append(cards, NewBaron())
	}

	for i := 0; i < handmaidCount; i++ {
		cards = append(cards, NewHandmaid())
	}

	for i := 0; i < princeCount; i++ {
		cards = append(cards, NewPrince())
	}

	for i := 0; i < chancellorCount; i++ {
		cards = append(cards, NewChancellor())
	}

	for i := 0; i < kingCount; i++ {
		cards = append(cards, NewKing())
	}

	for i := 0; i < countessCount; i++ {
		cards = append(cards, NewCountess())
	}

	for i := 0; i < princessCount; i++ {
		cards = append(cards, NewPrincess())
	}

	d.cards = cards
}
