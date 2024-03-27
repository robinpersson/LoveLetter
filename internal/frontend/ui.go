package frontend

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/robinpersson/LoveLetter/internal/chat"
	"golang.org/x/net/websocket"
)

const (
	WebsocketEndpoint = "ws://:3000/"
	WebsocketOrigin   = "http://"

	MessageWidget    = "messages"
	UsersWidget      = "players"
	InputWidget      = "send"
	CardsWidget      = "your cards"
	ControlsWidget   = "game control"
	ActionsWidget    = "actions"
	RulesWidget      = "rules"
	GuardWidget      = "guard action"
	PriestWidget     = "priest action"
	BaronWidget      = "baron action"
	PrinceWidget     = "prince action"
	ChancellorWidget = "chancellor action"
)

type UI struct {
	*gocui.Gui
	rulesOpen  bool
	username   string
	connection *websocket.Conn
}

func NewUI() (*UI, error) {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, fmt.Errorf("NewUI: %w", err)
	}

	return &UI{Gui: gui, rulesOpen: false}, nil
}

func (ui *UI) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	g.Cursor = true

	if controls, err := g.SetView(ControlsWidget, 0, 0, maxX-25, maxY-30); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		controls.Title = ControlsWidget
		controls.Autoscroll = true
		controls.Wrap = true
		fmt.Fprint(controls, "- start new game: Ctrl+S\n- toggle rules: Ctrl+R")

	}

	if messages, err := g.SetView(MessageWidget, 0, 5, maxX-15, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		messages.Title = MessageWidget
		messages.Autoscroll = true
		messages.Wrap = true
	}

	//if action, err := g.SetView(ActionsWidget, 0, maxY-12, maxX-20, maxY-1); err != nil {
	//	if err != gocui.ErrUnknownView {
	//		return err
	//	}
	//	action.Title = ActionsWidget
	//	action.Autoscroll = false
	//	action.Wrap = true
	//	action.Editable = false
	//}

	if input, err := g.SetView(InputWidget, 0, maxY-5, maxX-20, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		input.Title = InputWidget
		input.Autoscroll = false
		input.Wrap = true
		input.Editable = true
	}

	if users, err := g.SetView(UsersWidget, maxX-33, 0, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		users.Title = UsersWidget
		users.Autoscroll = false
		users.Wrap = true
	}

	if cards, err := g.SetView(CardsWidget, maxX-33, maxY-10, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		cards.Title = CardsWidget
		cards.Autoscroll = false
		cards.Wrap = true
	}

	g.SetCurrentView(InputWidget)

	return nil
}

func (ui *UI) SetKeyBindings(g *gocui.Gui) error {
	if err := g.SetKeybinding(InputWidget, gocui.KeyCtrlC, gocui.ModNone, ui.Quit); err != nil {
		return err
	}

	if err := g.SetKeybinding(InputWidget, gocui.KeyEnter, gocui.ModNone, ui.WriteMessage); err != nil {
		return err
	}

	if err := g.SetKeybinding(InputWidget, gocui.KeyCtrlS, gocui.ModNone, ui.StartGame); err != nil {
		return err
	}

	//if err := g.SetKeybinding(InputWidget, gocui.KeyCtrlH, gocui.ModNone, ui.PlayCurrentCard); err != nil {
	//	return err
	//}
	//
	//if err := g.SetKeybinding(InputWidget, gocui.KeyCtrlP, gocui.ModNone, ui.PlayPickedCard); err != nil {
	//	return err
	//}

	if err := g.SetKeybinding("", gocui.KeyCtrlR, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			if ui.rulesOpen {
				ui.rulesOpen = false
				return ui.DeleteView("Rules")
			}
			ui.rulesOpen = true
			return ui.RulesView()
		}); err != nil {
		return err
	}

	return nil
}

func (ui *UI) SetUsername(username string) {
	ui.username = username
}

func (ui *UI) SetConnection(connection *websocket.Conn) {
	ui.connection = connection
}

func (ui *UI) Connect(username string) error {
	config, err := websocket.NewConfig(WebsocketEndpoint, WebsocketOrigin)
	if err != nil {
		return err
	}

	config.Header.Set("Username", username)

	connection, err := websocket.DialConfig(config)
	if err != nil {
		return err
	}

	ui.SetConnection(connection)

	return nil
}

func (ui *UI) WriteMessage(_ *gocui.Gui, v *gocui.View) error {
	message := chat.NewMessage(chat.Regular, ui.username, v.Buffer())

	if err := websocket.JSON.Send(ui.connection, message); err != nil {
		return fmt.Errorf("UI.WriteMessage: %w", err)
	}

	v.SetCursor(0, 0)
	v.Clear()

	return nil
}

func (ui *UI) ReadMessage() error {
	for {
		var message chat.Message
		if err := websocket.JSON.Receive(ui.connection, &message); err != nil {
			return fmt.Errorf("UI.ReadMessage: %w", err)
		}

		ui.Update(func(g *gocui.Gui) error {
			//ui.DeleteView(PriestWidget)
			//ui.DeleteView(BaronWidget)
			//ui.DeleteView(PrinceWidget)
			//ui.DeleteView(GuardWidget)
			switch message.Type {
			case chat.Regular, chat.StartGame:
				view, err := ui.View(MessageWidget)
				if err != nil {
					return fmt.Errorf("UI.ReadMessage: %w", err)
				}
				fmt.Fprint(view, message.Formatted())
			case chat.CardMessage:
				view, err := ui.View(CardsWidget)
				if err != nil {
					return fmt.Errorf("UI.ReadMessage: %w", err)
				}
				view.Clear()
				fmt.Fprint(view, message.Formatted())
			case chat.ActionsMessage:
				ui.ShowPickCardsView()
			case chat.Connected, chat.Disconnected:
				view, err := ui.View(UsersWidget)
				if err != nil {
					return fmt.Errorf("UI.ReadMessage: %w", err)
				}
				view.Clear()
				fmt.Fprint(view, message.Text)
			case chat.Guard:
				ui.ShowGuardActionView(g, message)
			case chat.Priest:
				ui.ShowPriestActionView(g, message)
			case chat.PriestResponse:
				ui.ShowPriestResponseActionView(g, message)
			case chat.Baron:
				// if no opponents what to do?
				ui.ShowBaronActionView(g, message)
			case chat.Prince:
				ui.ShowPrinceActionView(g, message)
			case chat.Chancellor:
				ui.ShowChancellorActionView(g, message)
			}

			return nil
		})
	}

	return nil
}

func (ui *UI) Quit(_ *gocui.Gui, _ *gocui.View) error {
	return gocui.ErrQuit
}

func (ui *UI) Serve() error {
	go ui.ReadMessage()

	if err := ui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return fmt.Errorf("UI.Serve: %w", err)
	}

	return nil
}
