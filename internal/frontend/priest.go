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

var priestOpponents []chat.UserInfo

func (ui *UI) ShowPriestActionView(_ *gocui.Gui, message chat.Message) error {
	maxX, maxY := ui.Size()
	yStart := maxY - int(float64(maxY)*0.84)
	items := len(message.Opponents)
	width := maxX - int(float64(maxX)*0.2)
	priestOpponents = nil
	if priest, err := ui.SetView(PriestWidget, 0, maxY-yStart-items, width, maxY-yStart+1); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			fmt.Println(err)
			return err
		}
		priest.Title = "select the player which card you want to see"
		priest.Highlight = true
		priest.SelBgColor = gocui.ColorGreen
		priest.BgColor = gocui.ColorGreen
		priestOpponents = message.Opponents
		_, _ = fmt.Fprint(priest, getOpponents(message.Opponents))
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
	yStart := maxY - int(float64(maxY)*0.84)
	items := 2
	width := maxX - int(float64(maxX)*0.2)

	//_ = ui.DeleteView(PriestWidget)

	if priest, err := ui.SetView(PriestWidget, 0, maxY-yStart-items, width, maxY-yStart+1); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			fmt.Println(err)
			return err
		}
		//priest.Clear()
		priest.Title = "priest"
		priest.Highlight = true
		priest.SelBgColor = gocui.ColorGreen
		priest.BgColor = gocui.ColorGreen
		_, _ = fmt.Fprint(priest, message.Text+"\nPress F1 to discard")
	}

	err := ui.SetKeybinding(InputWidget, gocui.KeyF1, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
		return ui.sendPriestDiscardMessage(message.OpponentPlayer.Name)

		//return _ = ui.DeleteView(PriestWidget)
	})

	return err
}

func (ui *UI) ViewCard(playerNumber int) error {

	maxX, maxY := ui.Size()
	yStart := maxY - int(float64(maxY)*0.84)
	items := 2
	width := maxX - int(float64(maxX)*0.2)

	ui.clearGuessCardBindings()

	v, _ := ui.View(PriestWidget)
	v.Clear()
	_ = ui.DeleteView(PriestWidget)

	if priest, err := ui.SetView(PriestWidget, 0, maxY-yStart-items, width, maxY-yStart+1); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			fmt.Println(err)
			return err
		}
		//priest.Clear()
		priest.Title = "priest"
		priest.Highlight = true
		priest.SelBgColor = gocui.ColorGreen
		priest.BgColor = gocui.ColorGreen

		card := ""
		playerName := ""
		fmt.Println(priestOpponents)
		for _, u := range priestOpponents {
			if u.Order == playerNumber {
				playerName = u.Name
				card = fmt.Sprintf("%s has %s", u.Name, u.CardInfo.Description)
			}
		}

		_ = ui.SetKeybinding(InputWidget, gocui.KeyF1, gocui.ModNone, func(gui *gocui.Gui, view *gocui.View) error {
			return ui.sendPriestDiscardMessage(playerName)

			//return _ = ui.DeleteView(PriestWidget)
		})

		_, _ = fmt.Fprint(priest, card+"\nPress F1 to discard")
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
	}

	if err := websocket.JSON.Send(ui.connection, message); err != nil {
		return fmt.Errorf("UI.WriteMessage: %w", err)
	}

	_ = ui.DeleteView(PriestWidget)
	ui.clearGuessCardBindings()
	return nil
}

func (ui *UI) Priest_PickPlayer1(g *gocui.Gui, _ *gocui.View) error {
	return ui.ViewCard(1)
}

func (ui *UI) Priest_PickPlayer2(g *gocui.Gui, v *gocui.View) error {
	return ui.ViewCard(2)
}

func (ui *UI) Priest_PickPlayer3(g *gocui.Gui, v *gocui.View) error {
	return ui.ViewCard(3)
}

func (ui *UI) Priest_PickPlayer4(g *gocui.Gui, v *gocui.View) error {
	return ui.ViewCard(4)
}

func (ui *UI) Priest_PickPlayer5(g *gocui.Gui, v *gocui.View) error {
	return ui.ViewCard(5)
}

func (ui *UI) Priest_PickPlayer6(g *gocui.Gui, v *gocui.View) error {
	return ui.ViewCard(6)
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
