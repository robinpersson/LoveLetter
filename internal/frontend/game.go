// -------------------------------------------------------------------
// Copyright (c) Axis Communications AB, SWEDEN. All rights reserved.
// -------------------------------------------------------------------

package frontend

import (
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/jroimartin/gocui"
	"github.com/robinpersson/LoveLetter/internal/chat"
	"golang.org/x/net/websocket"
	"time"
)

func (ui *UI) StartGame(g *gocui.Gui, v *gocui.View) error {

	message := chat.NewMessage(chat.StartGame, ui.username, "started the game\n")

	if err := websocket.JSON.Send(ui.connection, message); err != nil {
		return fmt.Errorf("UI.WriteMessage: %w", err)
	}

	v.SetCursor(0, 0)
	v.Clear()

	_ = g.DeleteKeybinding(InputWidget, gocui.KeyCtrlS, gocui.ModNone)

	return nil
}

func (ui *UI) NextPlayer() error {

	message := chat.NewMessage(chat.NextPlayer, ui.username, "")

	if err := websocket.JSON.Send(ui.connection, message); err != nil {
		return fmt.Errorf("UI.WriteMessage: %w", err)
	}

	return nil
}

func (ui *UI) ShowRoundFinishedView(_ *gocui.Gui, message chat.Message) error {
	maxX, maxY := ui.Size()
	//name := fmt.Sprintf("v%v", 0)
	//0, 0, maxX-25, maxY-30
	//ui.DeleteAllViews()

	v, err := ui.SetView("RoundOver", 0, 0, maxX, maxY)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Title = ""
		v.Frame = true

		myFigure := figure.NewFigure("Round Over", "larry3d", true)

		winnerText := "\n\nWinner(s) are:\n"

		for _, winner := range message.RoundOver.Winners {
			wFig := figure.NewFigure(winner.Name, "avatar", true)
			winnerText += wFig.String()
			winnerText += "\n"
		}

		if message.RoundOver.WinnerCard.Name == "" {
			winnerText += "\nby being the last man standing\n\n"
		} else {
			winnerText += fmt.Sprintf("\nwith highest card: %s\n\n", message.RoundOver.WinnerCard.Description)
		}

		if message.RoundOver.SpyWinner.Name != "" {
			winnerText += fmt.Sprintf("extra point to %s for having an [0]Spy\n\n", message.RoundOver.SpyWinner.Name)
		}

		winnerText += fmt.Sprintf("face down card was: %s\n\n", message.RoundOver.OutCard.Description)

		winnerText += "next round will start shortly..."
		_, _ = fmt.Fprint(v, myFigure.String()+winnerText)

		ui.startCountDown(message.RoundOver.Winners[0].Order, message.RoundOver.Winners[0].Name)
	}
	if _, err := ui.SetCurrentView("RoundOver"); err != nil {
		return err
	}
	return nil
}

func (ui *UI) startCountDown(winnerOrder int, winnerName string) {

	ticker := time.NewTicker(1 * time.Second)
	quit := make(chan struct{})
	countdown := 15
	go func() {
		for {
			select {
			case <-ticker.C:
				//fmt.Printf("\x0cOn %d/10", countdown)
				//fmt.Printf("\rNext round starts in %d", countdown)
				//fmt.Print(countdown)
				countdown--
				if countdown == 0 {
					close(quit)
				}
			case <-quit:
				ticker.Stop()
				if ui.username == winnerName {

					ui.StartNewRound(winnerOrder)
				}

				return
			}
		}
	}()

}

func (ui *UI) ShowGameFinishedView(_ *gocui.Gui, message chat.Message) error {
	maxX, maxY := ui.Size()
	//name := fmt.Sprintf("v%v", 0)
	//0, 0, maxX-25, maxY-30
	v, err := ui.SetView("GameOver", 0, 0, maxX-1, maxY-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true
		v.Title = "Round over"
		v.Frame = false

		fmt.Fprintln(v, "Game over")
	}
	if _, err := ui.SetCurrentView("GameOver"); err != nil {
		return err
	}
	return nil
}

func (ui *UI) StartNewRound(winnerOrder int) error {
	message := chat.Message{
		Type:              chat.NewRound,
		LatestWinnerOrder: winnerOrder,
	}

	if err := websocket.JSON.Send(ui.connection, message); err != nil {
		return fmt.Errorf("UI.WriteMessage: %w", err)
	}

	return nil
}

func (ui *UI) Clear() error {
	ui.SetManagerFunc(ui.Layout)
	return ui.SetKeyBindings(ui.Gui)
}

func (ui *UI) GameControl(isAdmin bool) error {
	v, err := ui.View(ControlsWidget)

	if err != nil {
		return err
	}

	v.Clear()

	if isAdmin {
		_, _ = fmt.Fprint(v, "start new game: Ctrl+S\ntoggle rules: Ctrl+R")
	} else {
		_, _ = fmt.Fprint(v, "toggle rules: Ctrl+R")
		return ui.DeleteKeybinding(InputWidget, gocui.KeyCtrlS, gocui.ModNone)
	}

	return err
}
