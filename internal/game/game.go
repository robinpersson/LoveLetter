package game

import (
	"github.com/robinpersson/LoveLetter/internal/card"
	"math/rand"
)

//type Game interface {
//	StartNewGame(player player.Player)
//	JoinGame(player player.Player, id string)
//	Start()
//	DealCards()
//}

type Game struct {
	Id int
	//Players []*chat.User
	Deck        card.Deck
	Started     bool
	FavorTokens int
	Round       int
}

func NewGame() *Game {
	return &Game{
		Id: rand.Intn(100),
		//Players: nil,
		Deck:    nil,
		Started: false,
	}
}

func (g *Game) StartNewGame(nrOfUsers int) {
	deck := card.NewDeck(nrOfUsers)
	g.Deck = deck
	g.Started = true
	g.FavorTokens = getFavorTokens(nrOfUsers)

	//g.Players = append(g.Players, player)
}

func getFavorTokens(nrOfUsers int) int {
	if nrOfUsers == 2 {
		return 6
	}

	if nrOfUsers == 3 {
		return 5
	}

	if nrOfUsers == 4 {
		return 4
	}

	return 3
}

func (g *Game) InsertCards(cardsToInsert []card.Card) {
	g.Deck.InsertCards(cardsToInsert)
}

func (g *Game) DealCards() {
	//for _, player := range g.Players {
	//	//player.DealCard()
	//
	//}
}

func (g *Game) PickCard() *card.Card {
	//fmt.Println(g.Deck)
	return g.Deck.PickCard()
}
