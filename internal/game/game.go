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
	Deck    card.Deck
	Started bool
}

func NewGame() *Game {
	return &Game{
		Id: rand.Intn(100),
		//Players: nil,
		Deck:    nil,
		Started: false,
	}
}

func (g *Game) StartNewGame() {
	deck := card.NewDeck()
	//deck.Shuffle()
	g.Deck = deck
	g.Started = true
	//g.Players = append(g.Players, player)
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
