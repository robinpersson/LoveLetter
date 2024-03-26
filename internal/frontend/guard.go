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

func (ui *UI) ShowGuardActionView(_ *gocui.Gui, message chat.Message) error {
	maxX, maxY := ui.Size()
	height := len(message.Opponents) + 7

	if guard, err := ui.SetView(GuardWidget, 0, maxY-height, maxX-33, maxY-6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		guard.Title = "select player to guess on"
		guard.Highlight = true
		guard.SelBgColor = gocui.ColorGreen
		guard.BgColor = gocui.ColorGreen
		fmt.Fprint(guard, getOpponents(message.Opponents))
	}

	for _, u := range message.Opponents {
		key, pickFunc := ui.guard_getKey(u.Order)
		if err := ui.SetKeybinding(InputWidget, key, gocui.ModNone, pickFunc); err != nil {
			return err
		}
	}

	return nil
}

func (ui *UI) GuessPlayer(playerNumber, card int) error {
	ui.clearGuessCardBindings()
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

	ui.DeleteView(GuardWidget)

	return nil
}

func (ui *UI) clearGuessCardBindings() {
	ui.DeleteKeybinding(InputWidget, gocui.KeyF1, gocui.ModNone)
	ui.DeleteKeybinding(InputWidget, gocui.KeyF2, gocui.ModNone)
	ui.DeleteKeybinding(InputWidget, gocui.KeyF3, gocui.ModNone)
	ui.DeleteKeybinding(InputWidget, gocui.KeyF4, gocui.ModNone)
	ui.DeleteKeybinding(InputWidget, gocui.KeyF5, gocui.ModNone)
	ui.DeleteKeybinding(InputWidget, gocui.KeyF6, gocui.ModNone)
	ui.DeleteKeybinding(InputWidget, gocui.KeyF7, gocui.ModNone)
	ui.DeleteKeybinding(InputWidget, gocui.KeyF8, gocui.ModNone)
	ui.DeleteKeybinding(InputWidget, gocui.KeyF9, gocui.ModNone)
}

func (ui *UI) GuessCard(playerNumber int) {
	//SPY
	ui.SetKeybinding(InputWidget, gocui.KeyF1, gocui.ModNone, func(g *gocui.Gui, _ *gocui.View) error {
		ui.GuessPlayer(playerNumber, 0)
		return nil
	})
	//PRIEST
	ui.SetKeybinding(InputWidget, gocui.KeyF2, gocui.ModNone, func(g *gocui.Gui, _ *gocui.View) error {
		ui.GuessPlayer(playerNumber, 2)
		return nil
	})
	//BARON
	ui.SetKeybinding(InputWidget, gocui.KeyF3, gocui.ModNone, func(g *gocui.Gui, _ *gocui.View) error {
		ui.GuessPlayer(playerNumber, 3)
		return nil
	})
	//HANDMAID
	ui.SetKeybinding(InputWidget, gocui.KeyF4, gocui.ModNone, func(g *gocui.Gui, _ *gocui.View) error {
		ui.GuessPlayer(playerNumber, 4)
		return nil
	})
	//PRINCE
	ui.SetKeybinding(InputWidget, gocui.KeyF5, gocui.ModNone, func(g *gocui.Gui, _ *gocui.View) error {
		ui.GuessPlayer(playerNumber, 5)
		return nil
	})
	//CHANCELLOR
	ui.SetKeybinding(InputWidget, gocui.KeyF6, gocui.ModNone, func(g *gocui.Gui, _ *gocui.View) error {
		ui.GuessPlayer(playerNumber, 6)
		return nil
	})
	//KING
	ui.SetKeybinding(InputWidget, gocui.KeyF7, gocui.ModNone, func(g *gocui.Gui, _ *gocui.View) error {
		ui.GuessPlayer(playerNumber, 7)
		return nil
	})
	//COUNTESS
	ui.SetKeybinding(InputWidget, gocui.KeyF8, gocui.ModNone, func(g *gocui.Gui, _ *gocui.View) error {
		ui.GuessPlayer(playerNumber, 8)
		return nil
	})
	//PRONCESS
	ui.SetKeybinding(InputWidget, gocui.KeyF9, gocui.ModNone, func(g *gocui.Gui, _ *gocui.View) error {
		ui.GuessPlayer(playerNumber, 9)
		return nil
	})
}

func (ui *UI) Guard_PickPlayer1(g *gocui.Gui, _ *gocui.View) error {
	g.DeleteKeybinding(InputWidget, gocui.KeyF1, gocui.ModNone)
	ui.printGuessCards(g)
	ui.GuessCard(1)
	return nil
}

func (ui *UI) Guard_PickPlayer2(g *gocui.Gui, v *gocui.View) error {
	g.DeleteKeybinding(InputWidget, gocui.KeyF2, gocui.ModNone)
	ui.printGuessCards(g)
	ui.GuessCard(2)
	return nil
}

func (ui *UI) Guard_PickPlayer3(g *gocui.Gui, v *gocui.View) error {
	g.DeleteKeybinding(InputWidget, gocui.KeyF3, gocui.ModNone)
	ui.printGuessCards(g)
	ui.GuessCard(3)
	return nil
}

func (ui *UI) Guard_PickPlayer4(g *gocui.Gui, v *gocui.View) error {
	g.DeleteKeybinding(InputWidget, gocui.KeyF4, gocui.ModNone)
	ui.printGuessCards(g)
	ui.GuessCard(4)
	return nil
}

func (ui *UI) Guard_PickPlayer5(g *gocui.Gui, v *gocui.View) error {
	g.DeleteKeybinding(InputWidget, gocui.KeyF5, gocui.ModNone)
	ui.printGuessCards(g)
	ui.GuessCard(5)
	return nil
}

func (ui *UI) Guard_PickPlayer6(g *gocui.Gui, v *gocui.View) error {
	g.DeleteKeybinding(InputWidget, gocui.KeyF6, gocui.ModNone)
	ui.printGuessCards(g)
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
