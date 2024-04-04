package chat

import (
	"fmt"
	"time"
)

type MessageType int8

const (
	Regular MessageType = iota
	Connected
	Disconnected
	StartGame
	NewRound
	Clear
	CardMessage
	PickCard
	PlayCurrentCard
	PlayPickedCard
	Guard
	GuardGuess
	Priest
	PriestRequest
	PriestResponse
	PriestDiscard
	Baron
	CompareHands
	Prince
	DiscardCard
	Chancellor
	InsertChancellorCards
	King
	TradeCards
	Deck
	NextPlayer
	RoundFinished
	GameFinished
)

type CardInfo struct {
	Value       int    `json:"value"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Index       int    `json:"index"`
}

type UserInfo struct {
	Name   string
	Number int
	Order  int
}

type Guess struct {
	PlayerOrder int
	Card        int
}

type RoundOver struct {
	Winners     []UserInfo
	SpyWinner   UserInfo
	WinnerCard  CardInfo
	OutCard     CardInfo
	GameWinners []UserInfo
}

type Message struct {
	Type              MessageType `json:"type"`
	From              string      `json:"from"`
	Text              string      `json:"text"`
	Timestamp         string      `json:"timestamp"`
	Cards             []CardInfo  `json:"cards"`
	CurrentCard       CardInfo    `json:"cardInfo"`
	PickedCard        CardInfo    `json:"pickedCard"`
	Opponents         []UserInfo  `json:"opponents"`
	GuardGuess        Guess       `json:"guardGuess"`
	OpponentPlayer    UserInfo    `json:"opponentPlayer"`
	ChancellorCards   []CardInfo  `json:"chancellorCards"`
	RoundOver         RoundOver   `json:"roundOver"`
	LatestWinnerOrder int         `json:"latesWinnerOrder"`
}

func NewMessage(msgType MessageType, from string, text string) *Message {
	return &Message{
		Type:      msgType,
		From:      from,
		Text:      text + "\n",
		Timestamp: time.Now().Format(time.TimeOnly),
	}
}

func (m *Message) SetTime(v time.Time) {
	m.Timestamp = v.Format(time.TimeOnly)
}

func (m *Message) Formatted() string {

	switch m.Type {
	case Regular:
		return fmt.Sprintf("%v %v: %v", m.Timestamp, m.From, m.Text)
	case CardMessage:
		currentCard := fmt.Sprintf("Current card: %s\n", m.CurrentCard.Description)
		pickedCard := fmt.Sprintf("Picked card: %s\n", m.PickedCard.Description)
		return currentCard + pickedCard
	default:
		return fmt.Sprintf("%v %v: %v", m.Timestamp, m.From, m.Text)
	}
}

func SendGameControlMessage(supervisor *Supervisor, text string) {

}
