package main

//
//import (
//	"fmt"
//	"github.com/robinpersson/LoveLetter/internal/game"
//)
//
//func main() {
//	//deck := card.NewDeck()
//	//deck.Shuffle()
//	//cards := deck.Cards()
//	//for _, c := range cards {
//	//	fmt.Println(fmt.Sprintf("%s", c.Name()))
//	//}
//	//fmt.Println(len(cards))
//	//fmt.Println("Outcard", deck.OutCard())
//
//	g := game.NewGame()
//	g.StartNewGame()
//
//	p1 := game.NewPlayer("robin", g)
//	p2 := game.NewPlayer("nisse", g)
//
//	g.JoinGame(&p1, "")
//	g.JoinGame(&p2, "")
//
//	g.Start()
//	g.DealCards()
//
//	fmt.Printf("P1 current card: %s\n", p1.CurrentCard.ToString())
//	fmt.Printf("P2 current card: %s\n", p2.CurrentCard.ToString())
//
//	p1.PickCard()
//	p2.PickCard()
//
//	fmt.Printf("P1 picked card: %s\n", p1.PickedCard.ToString())
//	fmt.Printf("P2 picked card: %s\n", p2.PickedCard.ToString())
//	fmt.Printf("deck: %d\n", len(g.Deck.Cards()))
//
//}
