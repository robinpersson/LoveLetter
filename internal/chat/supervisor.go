package chat

import (
	"fmt"
	"github.com/robinpersson/LoveLetter/internal/card"
	"github.com/robinpersson/LoveLetter/internal/game"
	"math/rand/v2"
	"sort"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

type Supervisor struct {
	Users []*User
	Game  *game.Game
	mu    sync.Mutex
}

func NewSupervisor() *Supervisor {
	return &Supervisor{
		Users: make([]*User, 0),
		Game:  game.NewGame(),
	}
}

func (s *Supervisor) PickCard() *card.Card {
	return s.Game.PickCard()
}

func (s *Supervisor) Join(user *User) {
	s.mu.Lock()

	s.Users = append(s.Users, user)

	s.mu.Unlock()

	notification := NewMessage(Connected, "System", s.CurrentUsers())
	notification.SetTime(time.Now())

	s.Broadcast(notification)
}

func (s *Supervisor) SendGameControlMessage(text string) {
	mess := NewMessage(Regular, "Game control", text)
	mess.SetTime(time.Now())
	s.Broadcast(mess)
}

func (s *Supervisor) SendPlayOrder() {
	notification := NewMessage(Connected, "System", s.CurrentSortedUsers())
	notification.SetTime(time.Now())

	s.Broadcast(notification)
}

func (s *Supervisor) Quit(user *User) {
	s.mu.Lock()

	for i := len(s.Users) - 1; i >= 0; i-- {
		if s.Users[i] == user {
			s.Users = append(s.Users[:i], s.Users[i+1:]...)
		}
	}

	s.mu.Unlock()

	notification := NewMessage(Disconnected, "System", s.CurrentUsers())
	notification.SetTime(time.Now())

	s.Broadcast(notification)
}

func (s *Supervisor) CurrentUsers() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	var users string
	for _, user := range s.Users {
		users += fmt.Sprintf("%s\n", user.Name)
	}

	return users
}

func (s *Supervisor) CurrentSortedUsers() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	var users string
	for _, u := range s.Users {
		if u.IsProtected {
			users += fmt.Sprintf("%d. \u001B[33;1m%s\u001B[0m\n", u.Order, u.Name)
		} else if u.Eliminated {
			users += fmt.Sprintf("\u001B[31;1m%s\u001B[0m\n", u.Name)
		} else {
			users += fmt.Sprintf("%d. \u001B[32;1m%s\u001B[0m\n", u.Order, u.Name)
		}

		if len(u.Cards.Played) > 0 {
			users += addUserCards(u)
		}
	}

	return users
}

func addUserCards(user *User) string {
	var cardStr string

	for _, c := range user.Cards.Played {
		cardStr += fmt.Sprintf("%s ", c.ShortString())
	}

	return cardStr + "\n-----------------------\n"
}

func (s *Supervisor) Broadcast(message *Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, user := range s.Users {
		user.Write(message)
	}

	return nil
}

func (s *Supervisor) BroadcastText(text, from string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	message := &Message{
		Type:      Regular,
		From:      from,
		Text:      text + "\n",
		Timestamp: time.Now().Format(time.TimeOnly),
	}

	for _, user := range s.Users {
		user.Write(message)
	}

	return nil
}

func (s *Supervisor) BroadcastToAllOther(message *Message, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, user := range s.Users {
		if user.Name != name {
			user.Write(message)
		}
	}

	return nil
}

func (s *Supervisor) SendToUser(message *Message, user User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	user.Write(message)

	return nil
}

func (s *Supervisor) StartGame(userStarted *User) error {

	s.Game.StartNewGame()

	s.Broadcast(NewMessage(Regular, userStarted.Name, "started the game"))
	time.Sleep(time.Millisecond * 100)

	rand.Shuffle(len(s.Users), func(i, j int) { s.Users[i], s.Users[j] = s.Users[j], s.Users[i] })

	s.SendGameControlMessage("Randomizing play order")
	time.Sleep(time.Millisecond * 100)

	for i, user := range s.Users {
		user.Order = i + 1
	}

	sort.SliceStable(s.Users, func(i, j int) bool {
		return s.Users[i].Order < s.Users[j].Order
	})

	s.SendPlayOrder()

	for _, user := range s.Users {
		user.DealCard()
	}

	s.Users[0].PickCard()

	return nil
}

func (s *Supervisor) ServeWS() func(connection *websocket.Conn) {
	return func(connection *websocket.Conn) {
		user := NewUser(connection.Request().Header.Get("Username"), connection, s)
		s.Join(user)

		user.Read()
	}
}

func (s *Supervisor) NextPlayer(order int) {
	if order == len(s.Users) && !s.Users[0].Eliminated {
		s.Users[0].PickCard()
		return
	}

	s.GetPlayer(order + 1)
}

func (s *Supervisor) GetPlayer(order int) {
	for _, user := range s.Users {
		if user.Order == order && !user.Eliminated {
			user.PickCard()
			return
		}
	}
}

func (s *Supervisor) GetPlayerByOrder(order int) *User {
	for _, user := range s.Users {
		if user.Order == order {
			return user
		}
	}

	return nil
}

func (s *Supervisor) EliminatePlayer(player *User) {
	player.Eliminated = true
	s.SendGameControlMessage(fmt.Sprintf("%s is eliminated", player.Name))
	s.CheckIfGameIsOver()
}

func (s *Supervisor) CheckIfGameIsOver() {
	var usersLeft []*User
	for _, user := range s.Users {
		if !user.Eliminated {
			usersLeft = append(usersLeft, user)
		}
	}

	if len(usersLeft) == 1 {

		user := s.GetPlayerByOrder(usersLeft[0].Order)
		// We have a winner
		user.Tokens += user.Tokens

		//Start new game
	}
}
