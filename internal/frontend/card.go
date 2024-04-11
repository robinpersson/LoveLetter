// -------------------------------------------------------------------
// Copyright (c) Axis Communications AB, SWEDEN. All rights reserved.
// -------------------------------------------------------------------

package frontend

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/robinpersson/LoveLetter/internal/chat"
	"golang.org/x/net/websocket"
)

func (ui *UI) ShowPickCardsView(message chat.Message) error {
	maxX, maxY := ui.Size()
	yStart := maxY - int(float64(maxY)*0.84)
	items := 2
	width := maxX - int(float64(maxX)*0.3)

	if action, err := ui.SetView(ActionsWidget, 0, maxY-yStart-items, width, maxY-yStart+1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		action.Title = "Play a card"
		action.Highlight = true
		action.SelBgColor = gocui.ColorGreen
		action.BgColor = gocui.ColorGreen
		_, _ = fmt.Fprint(action, getCardsToPick(message.Cards))
	}

	for i, _ := range message.Cards {
		key, pickFunc := ui.pick_getKey(i)
		if err := ui.SetKeybinding(InputWidget, key, gocui.ModNone, pickFunc); err != nil {
			return err
		}
	}

	return nil
}

func (ui *UI) pick_getKey(userOrder int) (gocui.Key, func(g *gocui.Gui, v *gocui.View) error) {
	switch userOrder {
	case 0:
		return gocui.KeyF1, ui.PlayCurrentCard
	case 1:
		return gocui.KeyF2, ui.PlayPickedCard
	default:
		return 0, nil
	}
}

func getCardsToPick(cardInfos []chat.CardInfo) string {
	var opp string
	for i, o := range cardInfos {
		opp += fmt.Sprintf("F%d. %s\n", i+1, o.Description)
	}

	return opp
}

func (ui *UI) PlayPickedCard(g *gocui.Gui, v *gocui.View) error {
	message := chat.NewMessage(chat.PlayPickedCard, ui.username, "Play picked card")

	if err := websocket.JSON.Send(ui.connection, message); err != nil {
		return fmt.Errorf("UI.WriteMessage: %w", err)
	}

	v.Clear()
	_ = ui.DeleteView(ActionsWidget)
	ui.clearGuessCardBindings()
	return nil
}

func (ui *UI) PlayCurrentCard(g *gocui.Gui, v *gocui.View) error {
	message := chat.NewMessage(chat.PlayCurrentCard, ui.username, "Play current card")

	fmt.Println(ui.connection.IsClientConn())
	if err := websocket.JSON.Send(ui.connection, message); err != nil {

		return fmt.Errorf("UI.WriteMessage: %w", err)
	}

	v.Clear()
	_ = ui.DeleteView(ActionsWidget)
	ui.clearGuessCardBindings()
	return nil
}

func getCards() string {
	return "F1. Spy\nF2. Priest\nF3. Baron\nF4. Handmaid\nF5. Prince\nF6. Chancellor\nF7. King\nF8. Countess\nF9. Princess"
}

func (ui *UI) printGuessCards(g *gocui.Gui, opponentName string) {
	maxX, maxY := ui.Size()

	yStart := maxY - int(float64(maxY)*0.84)
	items := 9
	//itemRow := int(maxY / items)
	width := maxX - int(float64(maxX)*0.2)

	g.SetView(GuardWidget, 0, maxY-yStart-items, width, maxY-yStart+1)
	view, _ := g.View(GuardWidget)
	view.Clear()
	view.Title = fmt.Sprintf("select which card you think %s has", opponentName)
	_, _ = fmt.Fprint(view, getCards())
}
