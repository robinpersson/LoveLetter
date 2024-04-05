// -------------------------------------------------------------------
// Copyright (c) Axis Communications AB, SWEDEN. All rights reserved.
// -------------------------------------------------------------------

package frontend

import (
	"errors"
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/robinpersson/LoveLetter/internal/chat"
	"golang.org/x/net/websocket"
)

func (ui *UI) ShowBaronActionView(_ *gocui.Gui, message chat.Message) error {
	maxX, maxY := ui.Size()
	yStart := maxY - int(float64(maxY)*0.84)
	items := len(message.Opponents)
	width := maxX - int(float64(maxX)*0.2)

	if baron, err := ui.SetView(BaronWidget, 0, maxY-yStart-items, width, maxY-yStart+1); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			fmt.Println(err)
			return err
		}
		baron.Title = "select the player you want to compare hands with"
		baron.Highlight = true
		baron.SelBgColor = gocui.ColorGreen
		baron.BgColor = gocui.ColorGreen
		_, _ = fmt.Fprint(baron, getOpponents(message.Opponents))
	}

	for _, u := range message.Opponents {
		key, pickFunc := ui.baron_getKey(u.Order)
		if err := ui.SetKeybinding(InputWidget, key, gocui.ModNone, pickFunc); err != nil {
			return err
		}
	}

	return nil
}

func (ui *UI) baron_getKey(userOrder int) (gocui.Key, func(g *gocui.Gui, v *gocui.View) error) {
	switch userOrder {
	case 1:
		return gocui.KeyF1, ui.Baron_PickPlayer1
	case 2:
		return gocui.KeyF2, ui.Baron_PickPlayer2
	case 3:
		return gocui.KeyF3, ui.Baron_PickPlayer3
	case 4:
		return gocui.KeyF4, ui.Baron_PickPlayer4
	case 5:
		return gocui.KeyF5, ui.Baron_PickPlayer5
	case 6:
		return gocui.KeyF6, ui.Baron_PickPlayer6
	default:
		return 0, nil
	}
}

func (ui *UI) CompareCard(playerNumber int) error {
	message := chat.Message{
		Type:           chat.CompareHands,
		From:           ui.username,
		OpponentPlayer: chat.UserInfo{Order: playerNumber},
	}

	if err := websocket.JSON.Send(ui.connection, message); err != nil {
		return fmt.Errorf("UI.WriteMessage: %w", err)
	}

	_ = ui.DeleteView(BaronWidget)
	ui.clearGuessCardBindings()

	return nil
}

func (ui *UI) Baron_PickPlayer1(g *gocui.Gui, _ *gocui.View) error {
	return ui.CompareCard(1)
}

func (ui *UI) Baron_PickPlayer2(g *gocui.Gui, v *gocui.View) error {
	return ui.CompareCard(2)
}

func (ui *UI) Baron_PickPlayer3(g *gocui.Gui, v *gocui.View) error {
	return ui.CompareCard(3)
}

func (ui *UI) Baron_PickPlayer4(g *gocui.Gui, v *gocui.View) error {
	return ui.CompareCard(4)
}

func (ui *UI) Baron_PickPlayer5(g *gocui.Gui, v *gocui.View) error {
	return ui.CompareCard(5)
}

func (ui *UI) Baron_PickPlayer6(g *gocui.Gui, v *gocui.View) error {
	return ui.CompareCard(6)
}
