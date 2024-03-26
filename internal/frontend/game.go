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

func (ui *UI) StartGame(g *gocui.Gui, v *gocui.View) error {

	message := chat.NewMessage(chat.StartGame, ui.username, "started the game\n")

	if err := websocket.JSON.Send(ui.connection, message); err != nil {
		return fmt.Errorf("UI.WriteMessage: %w", err)
	}

	v.SetCursor(0, 0)
	v.Clear()

	g.DeleteKeybinding(InputWidget, gocui.KeyCtrlS, gocui.ModNone)

	return nil
}
