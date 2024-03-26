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

func (ui *UI) ShowPriestActionView(_ *gocui.Gui, message chat.Message) error {
	maxX, maxY := ui.Size()
	height := len(message.Opponents) + 7

	if priest, err := ui.SetView(PriestWidget, 0, maxY-height, maxX-33, maxY-6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		priest.Title = "select the player which card you want to see"
		priest.Highlight = true
		priest.SelBgColor = gocui.ColorGreen
		priest.BgColor = gocui.ColorGreen
		fmt.Fprint(priest, getOpponents(message.Opponents))
	}

	for _, u := range message.Opponents {
		key, pickFunc := ui.priest_getKey(u.Order)
		if err := ui.SetKeybinding(InputWidget, key, gocui.ModNone, pickFunc); err != nil {
			return err
		}
	}

	return nil
}

func (ui *UI) ShowPriestResponseActionView(_ *gocui.Gui, message chat.Message) error {
	maxX, maxY := ui.Size()

	ui.DeleteView(PriestWidget)

	if priest, err := ui.SetView(PriestWidget, 0, maxY-9, maxX-33, maxY-6); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		//priest.Clear()
		priest.Title = "priest"
		priest.Highlight = true
		priest.SelBgColor = gocui.ColorGreen
		priest.BgColor = gocui.ColorGreen
		fmt.Fprint(priest, message.Text+"\nPress F1 to discard")
	}

	err := ui.SetKeybinding(InputWidget, gocui.KeyF1, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		ui.DeleteKeybinding(InputWidget, gocui.KeyF1, gocui.ModNone)
		ui.sendPriestDiscardMessage(message.PriestPlayer.Name)
		return ui.DeleteView(PriestWidget)
	})

	if err != nil {
		return err
	}

	return nil
}

func (ui *UI) ViewCard(playerNumber int) error {
	message := chat.Message{
		Type:         chat.PriestRequest,
		From:         ui.username,
		PriestPlayer: chat.UserInfo{Order: playerNumber},
	}

	if err := websocket.JSON.Send(ui.connection, message); err != nil {
		return fmt.Errorf("UI.WriteMessage: %w", err)
	}

	return nil
}

func (ui *UI) sendPriestDiscardMessage(name string) error {
	message := chat.Message{
		Type: chat.Regular,
		From: ui.username,
		Text: fmt.Sprintf("looked at %s card\n", name),
	}

	if err := websocket.JSON.Send(ui.connection, message); err != nil {
		return fmt.Errorf("UI.WriteMessage: %w", err)
	}

	message = chat.Message{
		Type: chat.PriestDiscard,
		From: ui.username,
		//Text: fmt.Sprintf("looked at %s card", name),
	}

	if err := websocket.JSON.Send(ui.connection, message); err != nil {
		return fmt.Errorf("UI.WriteMessage: %w", err)
	}

	return nil
}

func (ui *UI) Priest_PickPlayer1(g *gocui.Gui, _ *gocui.View) error {
	g.DeleteKeybinding(InputWidget, gocui.KeyF1, gocui.ModNone)
	ui.ViewCard(1)
	return nil
}

func (ui *UI) Priest_PickPlayer2(g *gocui.Gui, v *gocui.View) error {
	g.DeleteKeybinding(InputWidget, gocui.KeyF2, gocui.ModNone)
	ui.ViewCard(2)
	return nil
}

func (ui *UI) Priest_PickPlayer3(g *gocui.Gui, v *gocui.View) error {
	g.DeleteKeybinding(InputWidget, gocui.KeyF3, gocui.ModNone)
	ui.ViewCard(3)
	return nil
}

func (ui *UI) Priest_PickPlayer4(g *gocui.Gui, v *gocui.View) error {
	g.DeleteKeybinding(InputWidget, gocui.KeyF4, gocui.ModNone)
	ui.ViewCard(4)
	return nil
}

func (ui *UI) Priest_PickPlayer5(g *gocui.Gui, v *gocui.View) error {
	g.DeleteKeybinding(InputWidget, gocui.KeyF5, gocui.ModNone)
	ui.ViewCard(5)
	return nil
}

func (ui *UI) Priest_PickPlayer6(g *gocui.Gui, v *gocui.View) error {
	g.DeleteKeybinding(InputWidget, gocui.KeyF6, gocui.ModNone)
	ui.ViewCard(6)
	return nil
}

func (ui *UI) priest_getKey(userOrder int) (gocui.Key, func(g *gocui.Gui, v *gocui.View) error) {
	switch userOrder {
	case 1:
		return gocui.KeyF1, ui.Priest_PickPlayer1
	case 2:
		return gocui.KeyF2, ui.Priest_PickPlayer2
	case 3:
		return gocui.KeyF3, ui.Priest_PickPlayer3
	case 4:
		return gocui.KeyF4, ui.Priest_PickPlayer4
	case 5:
		return gocui.KeyF5, ui.Priest_PickPlayer5
	case 6:
		return gocui.KeyF6, ui.Priest_PickPlayer6
	default:
		return 0, nil
	}
}
