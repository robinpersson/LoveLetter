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

func (ui *UI) ShowPrinceActionView(_ *gocui.Gui, message chat.Message) error {
	maxX, maxY := ui.Size()
	yStart := maxY - int(float64(maxY)*0.84)
	items := len(message.Opponents)
	width := maxX - int(float64(maxX)*0.2)

	if prince, err := ui.SetView(PrinceWidget, 0, maxY-yStart-items, width, maxY-yStart+1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		prince.Title = "select the player that must discard his card"
		prince.Highlight = true
		prince.SelBgColor = gocui.ColorGreen
		prince.BgColor = gocui.ColorGreen
		_, _ = fmt.Fprint(prince, getOpponents(message.Opponents))
	}

	for _, u := range message.Opponents {
		key, pickFunc := ui.prince_getKey(u.Order)
		if err := ui.SetKeybinding(InputWidget, key, gocui.ModNone, pickFunc); err != nil {
			return err
		}
	}

	return nil
}

func (ui *UI) Prince_PickPlayer1(g *gocui.Gui, _ *gocui.View) error {
	return ui.DiscardCard(1)
}

func (ui *UI) Prince_PickPlayer2(g *gocui.Gui, v *gocui.View) error {
	return ui.DiscardCard(2)
}

func (ui *UI) Prince_PickPlayer3(g *gocui.Gui, v *gocui.View) error {
	return ui.DiscardCard(3)
}

func (ui *UI) Prince_PickPlayer4(g *gocui.Gui, v *gocui.View) error {
	return ui.DiscardCard(4)
}

func (ui *UI) Prince_PickPlayer5(g *gocui.Gui, v *gocui.View) error {
	return ui.DiscardCard(5)
}

func (ui *UI) Prince_PickPlayer6(g *gocui.Gui, v *gocui.View) error {
	return ui.DiscardCard(6)
}

func (ui *UI) prince_getKey(userOrder int) (gocui.Key, func(g *gocui.Gui, v *gocui.View) error) {
	switch userOrder {
	case 1:
		return gocui.KeyF1, ui.Prince_PickPlayer1
	case 2:
		return gocui.KeyF2, ui.Prince_PickPlayer2
	case 3:
		return gocui.KeyF3, ui.Prince_PickPlayer3
	case 4:
		return gocui.KeyF4, ui.Prince_PickPlayer4
	case 5:
		return gocui.KeyF5, ui.Prince_PickPlayer5
	case 6:
		return gocui.KeyF6, ui.Prince_PickPlayer6
	default:
		return 0, nil
	}
}

func (ui *UI) DiscardCard(playerNumber int) error {

	message := chat.Message{
		Type:           chat.DiscardCard,
		From:           ui.username,
		OpponentPlayer: chat.UserInfo{Order: playerNumber},
	}

	if err := websocket.JSON.Send(ui.connection, message); err != nil {
		return fmt.Errorf("UI.WriteMessage: %w", err)
	}

	_ = ui.DeleteView(PrinceWidget)
	ui.clearGuessCardBindings()

	return nil
}
