package chat

import (
	"fmt"
	"github.com/robinpersson/LoveLetter/internal/card"
	"time"

	"golang.org/x/net/websocket"
)

type User struct {
	Name                string `json:"name"`
	Order               int
	Connection          *websocket.Conn `json:"-"`
	Egress              chan *Message   `json:"-"`
	Supervisor          *Supervisor     `json:"-"`
	Cards               *Cards          `json:"-"`
	IsInTurn            bool            `json:"-"`
	Tokens              int             `json:"-"`
	IsProtected         bool            `json:"-"`
	Eliminated          bool            `json:"-"`
	ChancellorCardInfos []CardInfo      `json:"-"`
	ChancellorCards     []card.Card     `json:"-"`
}

type Cards struct {
	Current *card.Card
	Played  []card.Card
	Picked  *card.Card
}

func (u *User) GetTokens() string {
	var t string
	for _ = range u.Tokens {
		t += "*"
	}
	return t
}

func (u *User) PickCard() *card.Card {
	c := u.Supervisor.PickCard()
	fmt.Println("CAAAARD: %+v", c)
	u.IsInTurn = true
	u.Cards.Picked = c
	u.SendPickedCard()
	return u.Cards.Picked
}

func (u *User) DiscardCardAndPick() *card.Card {
	c := u.Supervisor.PickCard()
	u.Cards.Played = append(u.Cards.Played, *u.Cards.Current)
	u.Cards.Current = c
	u.SendCurrentCard()
	time.Sleep(time.Millisecond * 200)
	u.Supervisor.SendPlayOrder()
	return u.Cards.Picked
}

func (u *User) DealCard() *card.Card {
	c := *u.Supervisor.PickCard()
	u.Cards.Current = &c
	if u.Order != 1 {
		u.SendCurrentCard()
	}

	return &c
}

func (u *User) DealCard_Temp() *card.Card {

	fmt.Println("666666666 %T", u.Supervisor)

	c := u.Supervisor.PickCard()

	fmt.Println("999999999 %T", c)

	u.Cards.Current = c
	if u.Order != 1 {
		u.SendCurrentCard()
	}

	return c
}

func (u *User) SendCurrentCard() {
	u.Write(u.GetCurrentCardMessage())
}

func (u *User) GetCurrentCardMessage() *Message {
	fmt.Printf("JJJJJJ %+v, %T", u.Cards.Current, u.Cards.Current)

	cc := *u.Cards.Current

	mess := &Message{
		Type:      CardMessage,
		From:      "Game master",
		Text:      fmt.Sprintf("%s\n", cc.ToString()),
		Timestamp: "",
		CurrentCard: CardInfo{
			Value:       cc.Value(),
			Name:        cc.Name(),
			Description: cc.ToString(),
		},
	}

	return mess
}

func (u *User) SendPickedCard() {
	pc := *u.Cards.Picked
	cc := *u.Cards.Current
	mess := &Message{
		Type:      CardMessage,
		From:      "Game control",
		Text:      fmt.Sprintf("%s\n", pc.ToString()),
		Timestamp: "",
		PickedCard: CardInfo{
			Value:       pc.Value(),
			Name:        pc.Name(),
			Description: pc.ToString(),
		},
		CurrentCard: CardInfo{
			Value:       cc.Value(),
			Name:        cc.Name(),
			Description: cc.ToString(),
		},
	}

	u.Supervisor.Broadcast(NewMessage(Regular, u.Name, "picked a card"))
	u.Write(mess)
	u.PrintPlayActions()
}

func (u *User) IsForcedToPlayCountess(currentCard *card.Card, pickedCard *card.Card) (bool, bool) {

	if currentCard == nil {
		return false, false
	}

	if (*currentCard).Name() == "Countess" {
		if (*pickedCard).Name() == "King" ||
			(*pickedCard).Name() == "Prince" {
			return true, true
		}
	}

	if pickedCard == nil {
		return false, false
	}

	if (*pickedCard).Name() == "Countess" {
		if (*currentCard).Name() == "King" ||
			(*currentCard).Name() == "Prince" {
			return true, false
		}
	}

	return false, false
}

func (u *User) PrintPlayActions() {
	currentCard := u.Cards.Current
	pickedCard := u.Cards.Picked
	forced, isCurrent := u.IsForcedToPlayCountess(currentCard, pickedCard)

	m := &Message{Type: PickCard}

	if forced {
		if isCurrent {
			m.Cards = []CardInfo{{
				Value:       (*currentCard).Value(),
				Name:        (*currentCard).Name(),
				Description: (*currentCard).ToString(),
				Index:       0,
			}}
		} else {
			m.Cards = []CardInfo{{
				Value:       (*pickedCard).Value(),
				Name:        (*pickedCard).Name(),
				Description: (*pickedCard).ToString(),
				Index:       0,
			}}
		}
	} else {
		m.Cards = []CardInfo{{
			Value:       (*currentCard).Value(),
			Name:        (*currentCard).Name(),
			Description: (*currentCard).ToString(),
			Index:       0,
		}}

		if pickedCard != nil {
			m.Cards = append(m.Cards, CardInfo{
				Value:       (*pickedCard).Value(),
				Name:        (*pickedCard).Name(),
				Description: (*pickedCard).ToString(),
				Index:       1,
			})
		}
	}

	u.Write(m)
}

func NewUser(name string, connection *websocket.Conn, supervisor *Supervisor) *User {
	return &User{
		Name:       name,
		Connection: connection,
		Egress:     make(chan *Message),
		Supervisor: supervisor,
		Cards:      &Cards{},
		//Player:     game.NewPlayer(*supervisor.Game),
	}
}

func (u *User) Read() {
	for {
		message := &Message{}
		if err := websocket.JSON.Receive(u.Connection, message); err != nil {
			// EOF connection closed by the client
			fmt.Printf("Supervisor quit: %v\n", err)
			u.Supervisor.Quit(u)
			break
		}
		fmt.Printf("READ %d %s, User: %s: \n", message.Type, message.Text, u.Name)
		message.SetTime(time.Now())

		switch message.Type {
		case StartGame:
			if !u.Supervisor.Game.Started {
				u.Supervisor.StartGame(u)
			} else {
				u.Write(NewMessage(Regular, "Game control", "Game already started\n"))
			}
		case NewRound:
			u.Supervisor.NewRound(message.LatestWinnerOrder)
		case Regular:
			u.Supervisor.Broadcast(message)
		case PlayCurrentCard:
			u.PlayCurrentCard()
		case PlayPickedCard:
			u.PlayPickedCard()
		case GuardGuess:
			u.GuessCard(message.GuardGuess.PlayerOrder, message.GuardGuess.Card)
		case PriestRequest:
			u.ViewOpponentCard(message.OpponentPlayer)
		case PriestDiscard:
			u.IsInTurn = false
			time.Sleep(time.Millisecond * 500)
			u.Supervisor.NextPlayer(u.Order)
		case CompareHands:
			u.CompareHands(message.OpponentPlayer.Order)
		case DiscardCard:
			u.DiscardCard(message.OpponentPlayer.Order)
		case InsertChancellorCards:
			u.InsertChancellorCards(message)
		case TradeCards:
			u.TradeCards(message.OpponentPlayer.Order)
		case NextPlayer:
			u.Supervisor.NextPlayer(u.Order)
		}

		u.Supervisor.BroadcastDeckCount()
	}
}

func (u *User) getChancellorCardByValue(value int) *card.Card {
	for _, chCard := range u.ChancellorCards {
		if chCard.Value() == value {
			return &chCard
		}
	}

	return nil
}

func (u *User) InsertChancellorCards(message *Message) {
	u.Cards.Current = u.getChancellorCardByValue(message.CurrentCard.Value)
	u.ChancellorCardInfos = nil

	var cardsToInsert []card.Card

	fmt.Println(u.ChancellorCards)

	for _, chCard := range message.ChancellorCards {
		cc := u.getChancellorCardByValue(chCard.Value)
		cardsToInsert = append(cardsToInsert, *cc)
	}

	u.ChancellorCards = nil

	u.Supervisor.InsertCards(cardsToInsert)

	u.Supervisor.BroadcastText(fmt.Sprintf("inserted %d cards in the bottom of the deck", len(cardsToInsert)), u.Name)
	time.Sleep(time.Millisecond * 200)

	u.IsInTurn = false
	u.Supervisor.NextPlayer(u.Order)

}

func (u *User) DiscardCard(opponentOrder int) {
	opponent := u.Supervisor.GetPlayerByOrder(opponentOrder)
	opponentCard := *opponent.Cards.Current

	u.Supervisor.BroadcastText(fmt.Sprintf("made %s discard his %s", opponent.Name, opponentCard.ToString()), u.Name)
	time.Sleep(time.Millisecond * 200)
	opponent.DiscardCardAndPick()

	u.IsInTurn = false
	u.Supervisor.NextPlayer(u.Order)
}

func (u *User) CompareHands(opponentOrder int) {
	opponent := u.Supervisor.GetPlayerByOrder(opponentOrder)
	opponentCard := *opponent.Cards.Current
	userCard := *u.Cards.Current

	if opponentCard.Value() > userCard.Value() {
		//Opponent win
		u.Supervisor.BroadcastText(fmt.Sprintf("compared card with %s, %s won", opponent.Name, opponent.Name), u.Name)
		time.Sleep(time.Millisecond * 200)
		u.Supervisor.EliminatePlayer(u)
	} else if opponentCard.Value() < userCard.Value() {
		//Player win
		u.Supervisor.BroadcastText(fmt.Sprintf("compared card with %s, %s won", opponent.Name, u.Name), u.Name)
		time.Sleep(time.Millisecond * 200)
		u.Supervisor.EliminatePlayer(opponent)
	} else {
		//It's a tie
		u.Supervisor.BroadcastText(fmt.Sprintf("compared card with %s, it was a tie", opponent.Name), u.Name)
		time.Sleep(time.Millisecond * 200)
	}

	u.IsInTurn = false
	u.Supervisor.NextPlayer(u.Order)
}

func (u *User) GuessCard(playerOrder, cardNumber int) {
	player := u.Supervisor.GetPlayerByOrder(playerOrder)

	playerCard := *player.Cards.Current
	guessedCard := card.GetCardByValue(cardNumber)

	if playerCard.Value() == cardNumber {
		//Right guess
		u.Supervisor.BroadcastText(fmt.Sprintf("guessed that %s had %s, which was correct", player.Name, guessedCard.ToString()), u.Name)
		time.Sleep(time.Millisecond * 100)
		u.Supervisor.EliminatePlayer(player)

	} else {
		u.Supervisor.BroadcastText(fmt.Sprintf("guessed that %s had %s, which was incorrect", player.Name, guessedCard.ToString()), u.Name)
	}
	time.Sleep(time.Millisecond * 500)
	u.IsInTurn = false
	u.Supervisor.NextPlayer(u.Order)

}

func (u *User) Write(message *Message) {
	//fmt.Printf("WRITE %d %s to %s\n", message.Type, message.Text, u.Name)
	if err := websocket.JSON.Send(u.Connection, message); err != nil {
		// EOF connection closed by the client
		u.Supervisor.Quit(u)
	}
}

func (u *User) broadcastPickedCard(cc card.Card) {
	m := &Message{
		Type:      Regular,
		From:      u.Name,
		Text:      fmt.Sprintf("played %s\n", cc.ToString()),
		Timestamp: time.Now().Format(time.TimeOnly),
	}
	u.Supervisor.Broadcast(m)
}

func (u *User) PlayPickedCard() {
	cc := *u.Cards.Picked //TODO: NULL
	u.Cards.Picked = nil
	u.PlayCard(cc)
}

func (u *User) PlayCurrentCard() {
	cc := *u.Cards.Current
	u.Cards.Current = u.Cards.Picked
	u.Cards.Picked = nil
	u.PlayCard(cc)
}

func (u *User) PlayCard(cc card.Card) {
	u.Cards.Played = append(u.Cards.Played, cc)
	u.IsProtected = cc.Name() == "Handmaid"

	if cc.Name() == "Princess" {
		u.Supervisor.EliminatePlayer(u)
	}

	u.Supervisor.SendPlayOrder()

	u.SendCurrentCard()

	u.broadcastPickedCard(cc)
	time.Sleep(time.Millisecond * 100)

	waitForPlay := u.PrintCardActions(cc)

	if !waitForPlay {
		u.IsInTurn = false
		u.Supervisor.NextPlayer(u.Order)
	}
}

func (u *User) TradeCards(playerOrder int) {
	opponent := u.Supervisor.GetPlayerByOrder(playerOrder)
	opponentCard := opponent.Cards.Current
	playerCard := u.Cards.Current

	opponent.Cards.Current = playerCard
	u.Cards.Current = opponentCard

	u.SendCurrentCard()

	u.Supervisor.BroadcastText(fmt.Sprintf(" traded card with %s", opponent.Name), u.Name)
	time.Sleep(time.Millisecond * 200)

	u.IsInTurn = false
	u.Supervisor.NextPlayer(u.Order)
}

func (u *User) PrintCardActions(c card.Card) bool {
	if c.Name() == "Guard" {
		opponents := u.getOpponents()
		if len(opponents) == 0 {
			return false
		}

		m := &Message{
			Type:      Guard,
			From:      "Game control",
			Text:      "guard",
			Opponents: opponents,
		}
		u.Write(m)
		return c.ActionText() != ""
	}

	if c.Name() == "Priest" {
		opponents := u.getOpponents()
		if len(opponents) == 0 {
			return false
		}

		m := &Message{
			Type:      Priest,
			From:      "Game control",
			Text:      "priest",
			Opponents: opponents,
		}
		u.Write(m)
		return c.ActionText() != ""
	}

	if c.Name() == "Baron" {
		opponents := u.getOpponents()
		if len(opponents) == 0 {
			return false
		}

		m := &Message{
			Type:      Baron,
			From:      "Game control",
			Text:      "baron",
			Opponents: opponents,
		}
		u.Write(m)
		return c.ActionText() != ""
	}

	if c.Name() == "Prince" {
		m := &Message{
			Type:      Prince,
			From:      "Game control",
			Text:      "prince",
			Opponents: u.getOpponentsIncludedCurrent(),
		}
		u.Write(m)
		return c.ActionText() != ""
	}

	if c.Name() == "Chancellor" {
		u.ChancellorCards = u.Supervisor.getChancellorCards(*u.Cards.Current)
		u.ChancellorCardInfos = u.Supervisor.MapCards(u.ChancellorCards)

		m := &Message{
			Type:            Chancellor,
			From:            "Game control",
			Text:            "chancellor",
			ChancellorCards: u.ChancellorCardInfos,
		}
		u.Write(m)

		if len(u.ChancellorCards) != 1 {
			time.Sleep(time.Millisecond * 100)
			u.Supervisor.BroadcastText(fmt.Sprintf("picked %d cards", len(u.ChancellorCards)-1), u.Name)
		}

		return c.ActionText() != ""
	}

	if c.Name() == "King" {
		m := &Message{
			Type:      King,
			From:      "Game control",
			Text:      "king",
			Opponents: u.getOpponents(),
		}
		u.Write(m)
		return c.ActionText() != ""
	}

	actionText := c.ActionText()
	if actionText != "" {
		m := &Message{
			Type: PickCard,
			From: "Game control",
			Text: actionText,
		}
		u.Write(m)
	}

	return actionText != ""
}

func (u *User) getOpponents() []UserInfo {
	var opponents []UserInfo
	for i, user := range u.Supervisor.Users {
		if user.Order != u.Order && !user.IsProtected {
			opponents = append(opponents, UserInfo{Name: user.Name, Number: i, Order: user.Order})
		}
	}

	return opponents
}

func (u *User) getOpponentsIncludedCurrent() []UserInfo {
	var opponents []UserInfo
	for i, user := range u.Supervisor.Users {
		if !user.IsProtected {
			opponents = append(opponents, UserInfo{Name: user.Name, Number: i, Order: user.Order})
		}
	}

	return opponents
}

func (u *User) ViewOpponentCard(opponentInfo UserInfo) {
	opponent := u.Supervisor.GetPlayerByOrder(opponentInfo.Order)
	opponentCard := *opponent.Cards.Current
	m := &Message{
		Type: PriestResponse,
		From: "Game control",
		Text: fmt.Sprintf("%s has %s", opponent.Name, opponentCard.ToString()),
		OpponentPlayer: UserInfo{
			Name: opponent.Name,
		},
	}
	u.Write(m)
}

func createCardMessage(card card.Card) *Message {
	return &Message{
		Type:      CardMessage,
		From:      "Game master",
		Text:      "deal card",
		Timestamp: "",
		CurrentCard: CardInfo{
			Value:       card.Value(),
			Name:        card.Name(),
			Description: card.ToString(),
		},
	}
}
