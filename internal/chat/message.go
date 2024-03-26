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
	CardMessage
	ActionsMessage
	PlayCurrentCard
	PlayPickedCard
	Guard
	GuardGuess
	Priest
	PriestRequest
	PriestResponse
	PriestDiscard
)

type CardInfo struct {
	Value       int    `json:"value"`
	Name        string `json:"name"`
	Description string `json:"description"`
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

type Message struct {
	Type         MessageType `json:"type"`
	From         string      `json:"from"`
	Text         string      `json:"text"`
	Timestamp    string      `json:"timestamp"`
	CurrentCard  CardInfo    `json:"cardInfo"`
	PickedCard   CardInfo    `json:"pickedCard"`
	Opponents    []UserInfo  `json:"opponents"`
	GuardGuess   Guess       `json:"guardGuess"`
	PriestPlayer UserInfo    `json:"priestPlayer"`
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
