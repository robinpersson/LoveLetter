// -------------------------------------------------------------------
// Copyright (c) Axis Communications AB, SWEDEN. All rights reserved.
// -------------------------------------------------------------------

package frontend

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/robinpersson/LoveLetter/internal/chat"
	"golang.org/x/net/websocket"
	"sort"
)

var currentCards []chat.CardInfo
var keepCard chat.CardInfo

func (ui *UI) ShowChancellorActionView(_ *gocui.Gui, message chat.Message) error {
	maxX, maxY := ui.Size()
	height := len(message.ChancellorCards) + 7

	if chancellor, err := ui.SetView(ChancellorWidget, 0, maxY-height, maxX-33, maxY-6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		chancellor.Title = "select card to keep"
		chancellor.Highlight = true
		chancellor.SelBgColor = gocui.ColorGreen
		chancellor.BgColor = gocui.ColorGreen
		currentCards = message.ChancellorCards
		fmt.Fprint(chancellor, getChancellorCards(message.ChancellorCards))
	}

	for _, u := range message.ChancellorCards {
		key, pickFunc := ui.chancellor_getKey(u.Index)
		if err := ui.SetKeybinding(InputWidget, key, gocui.ModNone, pickFunc); err != nil {
			return err
		}
	}

	return nil
}

func (ui *UI) ShowChancellorActionOrderView() error {
	ui.DeleteView(ChancellorWidget)
	maxX, maxY := ui.Size()
	height := len(currentCards) + 7

	if chancellor, err := ui.SetView(ChancellorWidget, 0, maxY-height, maxX-33, maxY-6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		chancellor.Title = "select card to put in bottom of deck"
		chancellor.Highlight = true
		chancellor.SelBgColor = gocui.ColorGreen
		chancellor.BgColor = gocui.ColorGreen
		//currentCards = message.ChancellorCards
		fmt.Fprint(chancellor, getChancellorCards(currentCards))
	}

	for i, _ := range currentCards {
		key, pickFunc := ui.chancellor_getKey2(i)
		if err := ui.SetKeybinding(InputWidget, key, gocui.ModNone, pickFunc); err != nil {
			return err
		}
	}

	return nil
}

func (ui *UI) chancellor_getKey2(index int) (gocui.Key, func(g *gocui.Gui, v *gocui.View) error) {
	switch index {
	case 0:
		return gocui.KeyF1, ui.Chancellor_FirstBottom
	case 1:
		return gocui.KeyF2, ui.Chancellor_SecondBottom
	default:
		return 0, nil
	}
}

func (ui *UI) chancellor_getKey(index int) (gocui.Key, func(g *gocui.Gui, v *gocui.View) error) {
	switch index {
	case 0:
		return gocui.KeyF1, ui.Chancellor_KeepCard1
	case 1:
		return gocui.KeyF2, ui.Chancellor_KeepCard2
	case 2:
		return gocui.KeyF3, ui.Chancellor_KeepCard3
	default:
		return 0, nil
	}
}

func (ui *UI) KeepCard(cardIndex int) error {
	ui.clearGuessCardBindings()
	keepCard = currentCards[cardIndex]
	currentCards = append(currentCards[:cardIndex], currentCards[cardIndex+1:]...)

	cardsView, _ := ui.View(CardsWidget)
	cardsView.Clear()
	currentCardPrint := fmt.Sprintf("Current card: %s\nPicked card:", keepCard.Description)
	fmt.Fprint(cardsView, currentCardPrint)

	if len(currentCards) == 0 {
		return ui.SendChancellorCards(false)
	}

	return ui.ShowChancellorActionOrderView()
}

func (ui *UI) SendChancellorCards(firstBottom bool) error {

	if firstBottom {
		sort.Slice(currentCards, func(i, j int) bool {
			return currentCards[i].Index > currentCards[j].Index
		})
	}

	message := chat.Message{
		Type:            chat.InsertChancellorCards,
		From:            ui.username,
		CurrentCard:     keepCard,
		ChancellorCards: currentCards,
	}

	if err := websocket.JSON.Send(ui.connection, message); err != nil {
		return fmt.Errorf("UI.WriteMessage: %w", err)
	}

	ui.DeleteView(ChancellorWidget)
	ui.clearGuessCardBindings()
	return nil
}

func (ui *UI) Chancellor_FirstBottom(g *gocui.Gui, _ *gocui.View) error {
	return ui.SendChancellorCards(true)
}

func (ui *UI) Chancellor_SecondBottom(g *gocui.Gui, _ *gocui.View) error {
	return ui.SendChancellorCards(false)
}

func (ui *UI) Chancellor_KeepCard1(g *gocui.Gui, _ *gocui.View) error {
	return ui.KeepCard(0)
}

func (ui *UI) Chancellor_KeepCard2(g *gocui.Gui, _ *gocui.View) error {
	return ui.KeepCard(1)
}

func (ui *UI) Chancellor_KeepCard3(g *gocui.Gui, _ *gocui.View) error {
	return ui.KeepCard(2)
}

func getChancellorCards(cardInfos []chat.CardInfo) string {
	var opponents string
	for i, o := range cardInfos {
		opponents += fmt.Sprintf("F%d. %s\n", i+1, o.Description)
	}

	return opponents
}
