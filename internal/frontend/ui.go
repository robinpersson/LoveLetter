package frontend

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/robinpersson/LoveLetter/internal/card"
	"github.com/robinpersson/LoveLetter/internal/chat"
	"golang.org/x/net/websocket"
)

const (
	WebsocketEndpoint = "ws://:3000/"
	WebsocketOrigin   = "http://"

	MessageWidget  = "messages"
	UsersWidget    = "players"
	InputWidget    = "send"
	CardsWidget    = "your cards"
	ControlsWidget = "game control"
	ActionsWidget  = "actions"
	RulesWidget    = "rules"
	GuardWidget    = "guard action"
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

	if action, err := g.SetView(ActionsWidget, 0, maxY-12, maxX-20, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		action.Title = ActionsWidget
		action.Autoscroll = false
		action.Wrap = true
		action.Editable = false
	}

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

func getRules() string {
	return card.GetPrincessRule() +
		card.GetCountessRule() +
		card.GetkingRule() +
		card.GetChancellorRule() +
		card.GetPrinceRule() +
		card.GetHandmaidRule() +
		card.GetBaronRule() +
		card.GetPriestRule() +
		card.GetGuardRule() +
		card.GetSpyRule()
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

	if err := g.SetKeybinding(InputWidget, gocui.KeyCtrlH, gocui.ModNone, ui.PlayCurrentCard); err != nil {
		return err
	}

	if err := g.SetKeybinding(InputWidget, gocui.KeyCtrlP, gocui.ModNone, ui.PlayPickedCard); err != nil {
		return err
	}

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

func (ui *UI) ShowGuardActionView(g *gocui.Gui, message chat.Message) error {
	//if err := g.SetKeybinding(GuardWidget, gocui.KeyF2, gocui.ModNone, ui.PlayPickedCard); err != nil {
	//	return err
	//}

	maxX, maxY := g.Size()
	g.Cursor = true

	if guard, err := g.SetView(GuardWidget, 0, maxY-17, maxX-33, maxY-5); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		guard.Title = GuardWidget
		guard.Autoscroll = false
		guard.Wrap = true
		//guard.Highlight = true
		guard.Overwrite = true
		guard.Title = "select player to guess on"
		fmt.Fprint(guard, getOpponents(message.GuardInfo))

	}

	for _, u := range message.GuardInfo.Opponents {
		key, pickFunc := ui.getKey(u.Order)
		if err := g.SetKeybinding(InputWidget, key, gocui.ModNone, pickFunc); err != nil {
			return err
		}

	}

	return nil
}

func (ui *UI) getKey(userOrder int) (gocui.Key, func(g *gocui.Gui, v *gocui.View) error) {
	switch userOrder {
	case 1:
		return gocui.KeyF1, ui.PickPlayer1
	case 2:
		return gocui.KeyF2, ui.PickPlayer2
	case 3:
		return gocui.KeyF3, ui.PickPlayer3
	case 4:
		return gocui.KeyF4, ui.PickPlayer3
	case 5:
		return gocui.KeyF5, ui.PickPlayer3
	case 6:
		return gocui.KeyF6, ui.PickPlayer3
	default:
		return 0, nil

	}
}

func getGuessCards() string {
	return "F1 Spy\nF2 Priest\nF3 Baron\nF4 Handmaid\nF5 Prince\nF6 Chancellor\nF7 King\nF8 Countess\nF9 Princess"
}

func (ui *UI) printGuessCards(g *gocui.Gui) {
	view, _ := g.View(GuardWidget)
	view.Clear()
	view.Title = "select card"
	fmt.Fprint(view, getGuessCards())
}

func (ui *UI) PickPlayer1(g *gocui.Gui, _ *gocui.View) error {
	g.DeleteKeybinding(InputWidget, gocui.KeyF1, gocui.ModNone)
	ui.printGuessCards(g)
	ui.GuessCard(1)
	return nil
}

func (ui *UI) PickPlayer2(g *gocui.Gui, v *gocui.View) error {
	g.DeleteKeybinding(InputWidget, gocui.KeyF2, gocui.ModNone)
	ui.printGuessCards(g)
	ui.GuessCard(2)
	return nil
}

func (ui *UI) PickPlayer3(g *gocui.Gui, v *gocui.View) error {
	g.DeleteKeybinding(InputWidget, gocui.KeyF3, gocui.ModNone)
	ui.printGuessCards(g)
	ui.GuessCard(3)
	return nil
}

func (ui *UI) PickPlayer4(g *gocui.Gui, v *gocui.View) error {
	g.DeleteKeybinding(InputWidget, gocui.KeyF4, gocui.ModNone)
	ui.printGuessCards(g)
	ui.GuessCard(4)
	return nil
}

func (ui *UI) PickPlayer5(g *gocui.Gui, v *gocui.View) error {
	g.DeleteKeybinding(InputWidget, gocui.KeyF5, gocui.ModNone)
	ui.printGuessCards(g)
	ui.GuessCard(5)
	return nil
}

func (ui *UI) PickPlayer6(g *gocui.Gui, v *gocui.View) error {
	g.DeleteKeybinding(InputWidget, gocui.KeyF6, gocui.ModNone)
	ui.printGuessCards(g)
	ui.GuessCard(6)
	return nil
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

	clearGuessCardBindings(ui.Gui)
}

func clearGuessCardBindings(g *gocui.Gui) {
	g.DeleteKeybinding(InputWidget, gocui.KeyF1, gocui.ModNone)
	g.DeleteKeybinding(InputWidget, gocui.KeyF2, gocui.ModNone)
	g.DeleteKeybinding(InputWidget, gocui.KeyF3, gocui.ModNone)
	g.DeleteKeybinding(InputWidget, gocui.KeyF4, gocui.ModNone)
	g.DeleteKeybinding(InputWidget, gocui.KeyF5, gocui.ModNone)
	g.DeleteKeybinding(InputWidget, gocui.KeyF6, gocui.ModNone)
	g.DeleteKeybinding(InputWidget, gocui.KeyF7, gocui.ModNone)
	g.DeleteKeybinding(InputWidget, gocui.KeyF8, gocui.ModNone)
	g.DeleteKeybinding(InputWidget, gocui.KeyF9, gocui.ModNone)
}

func getOpponents(info chat.GuardInfo) string {
	var opponents string
	for _, o := range info.Opponents {
		opponents += fmt.Sprintf("F%d. %s", o.Order, o.Name)
	}

	return opponents
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

func (ui *UI) GuessPlayer(playerNumber, card int) error {

	//message := chat.NewMessage(chat.GuardGuess, ui.username, fmt.Sprintf("guessed ")
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

func (ui *UI) WriteMessage(_ *gocui.Gui, v *gocui.View) error {
	message := chat.NewMessage(chat.Regular, ui.username, v.Buffer())

	if err := websocket.JSON.Send(ui.connection, message); err != nil {
		return fmt.Errorf("UI.WriteMessage: %w", err)
	}

	v.SetCursor(0, 0)
	v.Clear()

	return nil
}

func (ui *UI) StartGame(_ *gocui.Gui, v *gocui.View) error {

	message := chat.NewMessage(chat.StartGame, ui.username, "started the game\n")

	if err := websocket.JSON.Send(ui.connection, message); err != nil {
		return fmt.Errorf("UI.WriteMessage: %w", err)
	}

	v.SetCursor(0, 0)
	v.Clear()

	return nil
}

func (ui *UI) PlayCurrentCard(_ *gocui.Gui, v *gocui.View) error {

	message := chat.NewMessage(chat.PlayCurrentCard, ui.username, "Play current card")

	if err := websocket.JSON.Send(ui.connection, message); err != nil {
		return fmt.Errorf("UI.WriteMessage: %w", err)
	}

	v.SetCursor(0, 0)
	v.Clear()

	return nil
}

func (ui *UI) PlayPickedCard(_ *gocui.Gui, v *gocui.View) error {

	message := chat.NewMessage(chat.PlayPickedCard, ui.username, "Play picked card")

	if err := websocket.JSON.Send(ui.connection, message); err != nil {
		return fmt.Errorf("UI.WriteMessage: %w", err)
	}

	v.SetCursor(0, 0)
	v.Clear()

	return nil
}

func (ui *UI) RulesView() error {
	maxX, maxY := ui.Size()
	//name := fmt.Sprintf("v%v", 0)
	//0, 0, maxX-25, maxY-30
	v, err := ui.SetView("Rules", 0, 0+5, maxX, maxY)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Wrap = true

		fmt.Fprintln(v, getRules())
	}
	if _, err := ui.SetCurrentView("Rules"); err != nil {
		return err
	}

	//fmt.Fprint(v, me)
	//views = append(views, name)
	return nil
}

func (ui *UI) ReadMessage() error {
	for {
		var message chat.Message
		if err := websocket.JSON.Receive(ui.connection, &message); err != nil {
			return fmt.Errorf("UI.ReadMessage: %w", err)
		}

		ui.Update(func(g *gocui.Gui) error {
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
				view, err := ui.View(ActionsWidget)
				if err != nil {
					return fmt.Errorf("UI.ReadMessage: %w", err)
				}
				view.Clear()

				fmt.Fprint(view, message.Text)
			case chat.Connected, chat.Disconnected:
				view, err := ui.View(UsersWidget)
				if err != nil {
					return fmt.Errorf("UI.ReadMessage: %w", err)
				}

				view.Clear()
				fmt.Fprint(view, message.Text)
			case chat.Guard:
				ui.ShowGuardActionView(g, message)
				//view.Clear()
				//fmt.Fprint(view, message.Text)
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
