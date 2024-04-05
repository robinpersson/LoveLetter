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

var selectedGuesPlayer string
var opponents []chat.UserInfo

func (ui *UI) ShowGuardActionView(_ *gocui.Gui, message chat.Message) error {
	maxX, maxY := ui.Size()
	yStart := maxY - int(float64(maxY)*0.84)
	items := len(message.Opponents)
	width := maxX - int(float64(maxX)*0.2)

	if guard, err := ui.SetView(GuardWidget, 0, maxY-yStart-items, width, maxY-yStart+1); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			fmt.Println(err)
			return err
		}
		guard.Title = "select player to guess on"
		guard.Highlight = true
		guard.SelBgColor = gocui.ColorGreen
		guard.BgColor = gocui.ColorGreen
		opponents = message.Opponents
		_, _ = fmt.Fprint(guard, getOpponents(message.Opponents))
	}

	for _, u := range message.Opponents {
		key, pickFunc := ui.guard_getKey(u.Order)
		if err := ui.SetKeybinding(InputWidget, key, gocui.ModNone, pickFunc); err != nil {
			return err
		}
	}

	return nil
}

func getOpponentNameByOrder(order int) string {
	for _, o := range opponents {
		if o.Order == order {
			return o.Name
		}
	}

	return ""
}

func (ui *UI) GuessPlayer(playerNumber, card int) error {

	message := chat.Message{
		Type: chat.GuardGuess,
		From: ui.username,
		GuardGuess: chat.Guess{
			PlayerOrder: playerNumber,
			Card:        card,
		},
	}

	if err := websocket.JSON.Send(ui.connection, message); err != nil {
		return fmt.Errorf("UI.WriteMessage: %w", err)
	}

	ui.clearGuessCardBindings()
	_ = ui.DeleteView(GuardWidget)

	return nil
}

func (ui *UI) clearGuessCardBindings() {
	_ = ui.DeleteKeybinding(InputWidget, gocui.KeyF1, gocui.ModNone)
	_ = ui.DeleteKeybinding(InputWidget, gocui.KeyF2, gocui.ModNone)
	_ = ui.DeleteKeybinding(InputWidget, gocui.KeyF3, gocui.ModNone)
	_ = ui.DeleteKeybinding(InputWidget, gocui.KeyF4, gocui.ModNone)
	_ = ui.DeleteKeybinding(InputWidget, gocui.KeyF5, gocui.ModNone)
	_ = ui.DeleteKeybinding(InputWidget, gocui.KeyF6, gocui.ModNone)
	_ = ui.DeleteKeybinding(InputWidget, gocui.KeyF7, gocui.ModNone)
	_ = ui.DeleteKeybinding(InputWidget, gocui.KeyF8, gocui.ModNone)
	_ = ui.DeleteKeybinding(InputWidget, gocui.KeyF9, gocui.ModNone)
}

func (ui *UI) GuessCard(playerNumber int) {
	//ui.clearGuessCardBindings()
	//SPY
	_ = ui.SetKeybinding(InputWidget, gocui.KeyF1, gocui.ModNone, func(g *gocui.Gui, _ *gocui.View) error {
		return ui.GuessPlayer(playerNumber, 0)
	})
	//PRIEST
	_ = ui.SetKeybinding(InputWidget, gocui.KeyF2, gocui.ModNone, func(g *gocui.Gui, _ *gocui.View) error {
		return ui.GuessPlayer(playerNumber, 2)
	})
	//BARON
	_ = ui.SetKeybinding(InputWidget, gocui.KeyF3, gocui.ModNone, func(g *gocui.Gui, _ *gocui.View) error {
		return ui.GuessPlayer(playerNumber, 3)
	})
	//HANDMAID
	_ = ui.SetKeybinding(InputWidget, gocui.KeyF4, gocui.ModNone, func(g *gocui.Gui, _ *gocui.View) error {
		return ui.GuessPlayer(playerNumber, 4)
	})
	//PRINCE
	_ = ui.SetKeybinding(InputWidget, gocui.KeyF5, gocui.ModNone, func(g *gocui.Gui, _ *gocui.View) error {
		return ui.GuessPlayer(playerNumber, 5)
	})
	//CHANCELLOR
	_ = ui.SetKeybinding(InputWidget, gocui.KeyF6, gocui.ModNone, func(g *gocui.Gui, _ *gocui.View) error {
		return ui.GuessPlayer(playerNumber, 6)
	})
	//KING
	_ = ui.SetKeybinding(InputWidget, gocui.KeyF7, gocui.ModNone, func(g *gocui.Gui, _ *gocui.View) error {
		return ui.GuessPlayer(playerNumber, 7)
	})
	//COUNTESS
	_ = ui.SetKeybinding(InputWidget, gocui.KeyF8, gocui.ModNone, func(g *gocui.Gui, _ *gocui.View) error {
		return ui.GuessPlayer(playerNumber, 8)
	})
	//PRINCESS
	_ = ui.SetKeybinding(InputWidget, gocui.KeyF9, gocui.ModNone, func(g *gocui.Gui, _ *gocui.View) error {
		return ui.GuessPlayer(playerNumber, 9)
	})

	//ui.clearGuessCardBindings()
}

func (ui *UI) Guard_PickPlayer1(g *gocui.Gui, _ *gocui.View) error {
	ui.printGuessCards(g, getOpponentNameByOrder(1))
	ui.GuessCard(1)
	return nil
}

func (ui *UI) Guard_PickPlayer2(g *gocui.Gui, v *gocui.View) error {
	ui.printGuessCards(g, getOpponentNameByOrder(2))
	ui.GuessCard(2)
	return nil
}

func (ui *UI) Guard_PickPlayer3(g *gocui.Gui, v *gocui.View) error {
	ui.printGuessCards(g, getOpponentNameByOrder(3))
	ui.GuessCard(3)
	return nil
}

func (ui *UI) Guard_PickPlayer4(g *gocui.Gui, v *gocui.View) error {
	ui.printGuessCards(g, getOpponentNameByOrder(4))
	ui.GuessCard(4)
	return nil
}

func (ui *UI) Guard_PickPlayer5(g *gocui.Gui, v *gocui.View) error {
	ui.printGuessCards(g, getOpponentNameByOrder(5))
	ui.GuessCard(5)
	return nil
}

func (ui *UI) Guard_PickPlayer6(g *gocui.Gui, v *gocui.View) error {
	ui.printGuessCards(g, getOpponentNameByOrder(6))
	ui.GuessCard(6)
	return nil
}

func (ui *UI) guard_getKey(userOrder int) (gocui.Key, func(g *gocui.Gui, v *gocui.View) error) {
	switch userOrder {
	case 1:
		return gocui.KeyF1, ui.Guard_PickPlayer1
	case 2:
		return gocui.KeyF2, ui.Guard_PickPlayer2
	case 3:
		return gocui.KeyF3, ui.Guard_PickPlayer3
	case 4:
		return gocui.KeyF4, ui.Guard_PickPlayer4
	case 5:
		return gocui.KeyF5, ui.Guard_PickPlayer5
	case 6:
		return gocui.KeyF6, ui.Guard_PickPlayer6
	default:
		return 0, nil
	}
}
