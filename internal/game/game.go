package game

import (
	"github.com/robinpersson/LoveLetter/internal/card"
	"math/rand"
	"time"
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
	deck := card.NewDeck()
	//deck.Shuffle()
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

//func (g *Game) JoinGame(player *Player, id string) {
//	g.Players = append(g.Players, player)
//}

func (g *Game) Start() {
	rand.NewSource(time.Now().UnixNano())
	//rand.Shuffle(len(g.Players), func(i, j int) { g.Players[i], g.Players[j] = g.Players[j], g.Players[i] })
	g.Started = true
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
