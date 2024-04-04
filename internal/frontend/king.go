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

func (ui *UI) ShowKingActionView(_ *gocui.Gui, message chat.Message) error {
	maxX, maxY := ui.Size()
	height := len(message.Opponents) + 7

	if king, err := ui.SetView(KingWidget, 0, maxY-height, maxX-33, maxY-6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		king.Title = "select user to swap card with"
		king.Highlight = true
		king.SelBgColor = gocui.ColorGreen
		king.BgColor = gocui.ColorGreen
		fmt.Fprint(king, getOpponents(message.Opponents))
	}

	for _, u := range message.Opponents {
		key, pickFunc := ui.king_getKey(u.Order)
		if err := ui.SetKeybinding(InputWidget, key, gocui.ModNone, pickFunc); err != nil {
			return err
		}
	}

	return nil
}

func (ui *UI) King_PickPlayer1(g *gocui.Gui, _ *gocui.View) error {
	return ui.TradeCards(1)
}

func (ui *UI) King_PickPlayer2(g *gocui.Gui, v *gocui.View) error {
	return ui.TradeCards(2)
}

func (ui *UI) King_PickPlayer3(g *gocui.Gui, v *gocui.View) error {
	return ui.TradeCards(3)
}

func (ui *UI) King_PickPlayer4(g *gocui.Gui, v *gocui.View) error {
	return ui.TradeCards(4)
}

func (ui *UI) King_PickPlayer5(g *gocui.Gui, v *gocui.View) error {
	return ui.TradeCards(5)
}

func (ui *UI) King_PickPlayer6(g *gocui.Gui, v *gocui.View) error {
	return ui.TradeCards(6)
}

func (ui *UI) king_getKey(userOrder int) (gocui.Key, func(g *gocui.Gui, v *gocui.View) error) {
	switch userOrder {
	case 1:
		return gocui.KeyF1, ui.King_PickPlayer1
	case 2:
		return gocui.KeyF2, ui.King_PickPlayer2
	case 3:
		return gocui.KeyF3, ui.King_PickPlayer3
	case 4:
		return gocui.KeyF4, ui.King_PickPlayer4
	case 5:
		return gocui.KeyF5, ui.King_PickPlayer5
	case 6:
		return gocui.KeyF6, ui.King_PickPlayer6
	default:
		return 0, nil
	}
}

func (ui *UI) TradeCards(playerNumber int) error {
	message := chat.Message{
		Type:           chat.TradeCards,
		From:           ui.username,
		OpponentPlayer: chat.UserInfo{Order: playerNumber},
	}

	if err := websocket.JSON.Send(ui.connection, message); err != nil {
		return fmt.Errorf("UI.WriteMessage: %w", err)
	}

	ui.DeleteView(KingWidget)
	ui.clearGuessCardBindings()

	return nil
}
