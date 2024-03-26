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

func (ui *UI) ShowPickCardsView() error {
	maxX, maxY := ui.Size()
	ui.Cursor = true

	if action, err := ui.SetView(ActionsWidget, 0, maxY-9, maxX-33, maxY-6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		action.Title = "Pick a card"
		action.Highlight = true
		action.SelBgColor = gocui.ColorGreen
		action.BgColor = gocui.ColorGreen
		fmt.Fprint(action, "Current card: F1\nPicked card: F2")
	}

	if err := ui.SetKeybinding(InputWidget, gocui.KeyF1, gocui.ModNone, ui.PlayCurrentCard); err != nil {
		return err
	}

	if err := ui.SetKeybinding(InputWidget, gocui.KeyF2, gocui.ModNone, ui.PlayPickedCard); err != nil {
		return err
	}

	return nil
}

func (ui *UI) PlayPickedCard(g *gocui.Gui, v *gocui.View) error {

	message := chat.NewMessage(chat.PlayPickedCard, ui.username, "Play picked card")

	if err := websocket.JSON.Send(ui.connection, message); err != nil {
		return fmt.Errorf("UI.WriteMessage: %w", err)
	}

	v.SetCursor(0, 0)
	v.Clear()

	deletePickCardBindings(g)

	return nil
}

func (ui *UI) PlayCurrentCard(g *gocui.Gui, v *gocui.View) error {

	message := chat.NewMessage(chat.PlayCurrentCard, ui.username, "Play current card")

	if err := websocket.JSON.Send(ui.connection, message); err != nil {
		return fmt.Errorf("UI.WriteMessage: %w", err)
	}

	v.SetCursor(0, 0)
	v.Clear()

	deletePickCardBindings(g)

	return nil
}

func deletePickCardBindings(g *gocui.Gui) {
	g.DeleteView(ActionsWidget)
	g.DeleteKeybinding(InputWidget, gocui.KeyF1, gocui.ModNone)
	g.DeleteKeybinding(InputWidget, gocui.KeyF2, gocui.ModNone)
}

func getCards() string {
	return "F1 Spy\nF2 Priest\nF3 Baron\nF4 Handmaid\nF5 Prince\nF6 Chancellor\nF7 King\nF8 Countess\nF9 Princess"
}

func (ui *UI) printGuessCards(g *gocui.Gui) {
	maxX, maxY := ui.Size()
	g.SetView(GuardWidget, 0, maxY-16, maxX-33, maxY-6)
	view, _ := g.View(GuardWidget)
	view.Clear()
	view.Title = "select card"
	fmt.Fprint(view, getCards())
}
