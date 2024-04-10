package chat

import (
	"fmt"
	"github.com/robinpersson/LoveLetter/internal/card"
	"github.com/robinpersson/LoveLetter/internal/game"
	"net/http"
	"sort"
	"strings"
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

func (s *Supervisor) InsertCards(cardsToInsert []card.Card) {
	s.Game.InsertCards(cardsToInsert)
}

func (s *Supervisor) Join(user *User) error {
	s.mu.Lock()

	if len(s.Users) == 0 {
		user.Admin = true
	}

	s.Users = append(s.Users, user)

	s.mu.Unlock()

	notification := NewMessage(Connected, "System", s.CurrentUsers())
	notification.SetTime(time.Now())

	s.Broadcast(notification)

	user.Write(&Message{Type: GameControl, IsAdmin: user.Admin})

	return nil
}

func (s *Supervisor) SendGameControlMessage(text string) {
	mess := NewMessage(Regular, "Game control", text)
	mess.SetTime(time.Now())
	s.Broadcast(mess)
}

func (s *Supervisor) BroadcastDeckCount() {
	if s.Game != nil && s.Game.Deck != nil {
		mess := NewMessage(Deck, "Game control", fmt.Sprintf("favor tokens: %d\nround: %d\ndeck: %d cards left", s.Game.FavorTokens, s.Game.Round, len(s.Game.Deck.Cards())))
		//mess.SetTime(time.Now())
		s.Broadcast(mess)
	}
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

	s.Broadcast(NewMessage(Regular, user.Name, "disconnected"))
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

		if u.Eliminated {
			users += fmt.Sprintf("%d. \u001B[31;1m%s\u001B[0m %s\n", u.Order, u.Name, u.GetTokens())
		} else if u.IsProtected {
			users += fmt.Sprintf("%d. \u001B[33;1m%s\u001B[0m %s\n", u.Order, u.Name, u.GetTokens())
		} else {
			users += fmt.Sprintf("%d. \u001B[32;1m%s\u001B[0m %s\n", u.Order, u.Name, u.GetTokens())
		}

		if u.Cards != nil && len(u.Cards.Played) > 0 {
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

	return cardStr + "\n\n"
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
func (s *Supervisor) NewRound(winnerOrder int) error {
	s.Broadcast(&Message{
		Type: Clear,
	})

	time.Sleep(time.Millisecond * 300)
	latestWinner := s.GetPlayerByOrder(winnerOrder)

	//s.StartGame(latestWinner)

	//return nil
	s.Game.StartNewGame(len(s.Users))
	s.Game.Round = s.Game.Round + 1

	//s.StartGame(latestWinner)

	for _, user := range s.Users {
		user.Eliminated = false
		user.IsProtected = false
		user.Cards = &Cards{}
		user.IsInTurn = false
	}

	s.Broadcast(NewMessage(Regular, latestWinner.Name, "started new round"))
	time.Sleep(time.Millisecond * 100)

	s.SendPlayOrder()

	latestWinner.DealCard()
	//time.Sleep(time.Millisecond * 100)
	for _, user := range s.Users {

		if user.Order != latestWinner.Order {
			user.DealCard()
		}
	}

	latestWinner.PickCard()

	return nil
}

func (s *Supervisor) StartGame(userStarted *User) error {
	s.Game.StartNewGame(len(s.Users))
	s.Game.Round = 1

	s.Broadcast(NewMessage(Regular, userStarted.Name, "started the game"))
	time.Sleep(time.Millisecond * 100)

	//TODO:
	//rand.Shuffle(len(s.Users), func(i, j int) { s.Users[i], s.Users[j] = s.Users[j], s.Users[i] })

	s.SendGameControlMessage("Randomizing play order\n")
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

func (s *Supervisor) UserExists(userName string) bool {
	for _, user := range s.Users {
		if strings.ToLower(user.Name) == strings.ToLower(userName) {
			return true
		}
	}

	return false
}

func (s *Supervisor) IsUp() func(connection *websocket.Conn) {
	return func(connection *websocket.Conn) {
		//userName := connection.Request().Header.Get("Username")
		//exist := s.UserExists(userName)
		//fmt.Printf("User %s exists: %v\n", userName, exist)
		//user := NewUser(userName, connection, s)
		//fmt.Printf("User %s connected\n", userName)
		//s.Join(user)
		//
		//user.Read()

	}
}

func (s *Supervisor) JoinWS() func(connection *websocket.Conn) {
	return func(connection *websocket.Conn) {
		userName := connection.Request().Header.Get("Username")
		//exist := s.UserExists(userName)
		//fmt.Printf("User %s exists: %v\n", userName, exist)
		user := NewUser(userName, connection, s)
		//fmt.Printf("User %s connected\n", userName)
		s.Join(user)
		//
		user.Read()

	}
}

func (s *Supervisor) UserNameTakenHandshake() func(config *websocket.Config, req *http.Request) error {
	return func(config *websocket.Config, req *http.Request) error {
		userName := req.Header.Get("Username")
		exist := s.UserExists(userName)
		fmt.Println(exist)
		if exist {

			return fmt.Errorf("username exists")
		}

		return nil
	}
}

func (s *Supervisor) UserNameTaken() func(connection *websocket.Conn) {
	return func(connection *websocket.Conn) {
		//userName := connection.Request().Header.Get("Username")
		//exist := s.UserExists(userName)
		//fmt.Println(exist)
		//connection.PayloadType = byte(3)
		////connection.Config().Header.Set("usernameexists", strconv.FormatBool(exist))
		//if exist {
		//	connection.WriteClose(http.StatusBadRequest)
		//} else {
		//	connection.WriteClose(http.StatusOK)
		//}

	}
}

func (s *Supervisor) NextPlayer(order int) {
	gameOver := s.IsGameOver()
	fmt.Println("GAME OVER:", gameOver)
	if gameOver {
		return
	}

	if order == len(s.Users) && !s.Users[0].Eliminated {
		fmt.Println("GOT HERE")
		s.Users[0].PickCard()
		return
	}
	fmt.Println("NOP GOT HERE")
	s.GetPlayer(order + 1)
}

func (s *Supervisor) GetPlayer(order int) {
	fmt.Println("order: ", order)

	for {

		if order > len(s.Users) {
			order = 1
		}

		user := s.GetPlayerByOrder(order)

		if !user.Eliminated {
			user.PickCard()
			break
		}

		order++

	}

	//for _, user := range s.Users {
	//	if user.Order == order && !user.Eliminated {
	//		user.PickCard()
	//		fmt.Println("picked")
	//		break
	//	} else if order <= len(s.Users) {
	//		order++
	//		fmt.Println("order++: ", order)
	//	} else {
	//		order = 1
	//		fmt.Println("order1: ", order)
	//	}
	//}
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
	s.SendPlayOrder()
	s.SendGameControlMessage(fmt.Sprintf("%s is eliminated", player.Name))
	time.Sleep(time.Millisecond * 200)
}

func (s *Supervisor) IsGameOver() bool {
	var usersLeft []*User
	for _, user := range s.Users {
		if !user.Eliminated {
			usersLeft = append(usersLeft, user)
		}
	}

	if len(usersLeft) == 1 {

		user := s.GetPlayerByOrder(usersLeft[0].Order)

		var spyWinner UserInfo
		if hasSpy(*user) {
			user.Tokens += 1
			spyWinner = UserInfo{Name: user.Name}
		}

		// We have a round winner
		user.Tokens += 1

		roundOver := RoundOver{
			Winners:   []UserInfo{{Name: user.Name, Order: user.Order}},
			SpyWinner: spyWinner,
			//WinnerCard: CardInfo{
			//	Value:       (*user.Cards.Current).Value(),
			//	Name:        (*user.Cards.Current).Name(),
			//	Description: (*user.Cards.Current).ToString(),
			//	Index:       0,
			//},
			OutCard: CardInfo{
				Value:       s.Game.Deck.OutCard().Value(),
				Name:        s.Game.Deck.OutCard().Name(),
				Description: s.Game.Deck.OutCard().ToString(),
				Index:       0,
			},
		}

		if user.Tokens >= s.Game.FavorTokens {
			s.BroadcastText("Game over", "Game control")
			time.Sleep(time.Millisecond * 200)
			s.Broadcast(&Message{
				Type:      GameFinished,
				From:      "Game control",
				RoundOver: roundOver,
			})
		} else {
			s.BroadcastText("Round over... wait for summary", "Game control")
			time.Sleep(time.Millisecond * 20000)
			s.Broadcast(&Message{
				Type:      RoundFinished,
				From:      "Game control",
				RoundOver: roundOver,
			})
		}

		return true
	}

	if len(s.Game.Deck.Cards()) == 0 {
		//s.BroadcastText("GAME OVER", "Game control")

		gameOver := false

		var winnerUsers []UserInfo
		winners := getWinners(usersLeft)
		spyCount := 0
		for _, winner := range winners {
			if hasSpy(*winner) {
				spyCount++
			}

			winner.Tokens += 1

			if winner.Tokens >= s.Game.FavorTokens {
				gameOver = true
				// Game finished
			}

			winnerUsers = append(winnerUsers, UserInfo{Name: winner.Name, Order: winner.Order})
		}

		var spyWinner UserInfo

		if spyCount == 1 {
			for _, winner := range winners {
				if hasSpy(*winner) {
					spyWinner = UserInfo{Name: winner.Name}
					winner.Tokens += 1
					if winner.Tokens >= s.Game.FavorTokens {
						gameOver = true
						// Game finished
					}
					break
				}
			}
		}

		roundOver := RoundOver{
			Winners:   winnerUsers,
			SpyWinner: spyWinner,
			WinnerCard: CardInfo{
				Value:       (*winners[0].Cards.Current).Value(),
				Name:        (*winners[0].Cards.Current).Name(),
				Description: (*winners[0].Cards.Current).ToString(),
				Index:       0,
			},
			OutCard: CardInfo{
				Value:       s.Game.Deck.OutCard().Value(),
				Name:        s.Game.Deck.OutCard().Name(),
				Description: s.Game.Deck.OutCard().ToString(),
				Index:       0,
			},
		}

		if gameOver {
			gameWinners := getGameWinners(winners)
			var gameWinnerUsers []UserInfo
			for _, winner := range gameWinners {
				gameWinnerUsers = append(gameWinnerUsers, UserInfo{Name: winner.Name})
			}
			roundOver.GameWinners = gameWinnerUsers
			s.Broadcast(&Message{
				Type:      GameFinished,
				From:      "Game control",
				RoundOver: roundOver,
			})
		} else {
			s.BroadcastText("Round over... wait for summary", "Game control")
			time.Sleep(time.Millisecond * 2000)
			s.Broadcast(&Message{
				Type:      RoundFinished,
				From:      "Game control",
				RoundOver: roundOver,
			})
		}

		return true

	}

	return false
}

func getWinners(users []*User) []*User {
	groupedUsers := make(map[int][]*User)

	for _, u := range users {
		groupedUsers[(*u.Cards.Current).Value()] = append(groupedUsers[(*u.Cards.Current).Value()], u)
	}

	sortedKeys := make([]int, 0, len(groupedUsers))

	for k := range groupedUsers {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sortedKeys)))

	return groupedUsers[sortedKeys[0]]
}

func getGameWinners(users []*User) []*User {
	groupedUsers := make(map[int][]*User)

	for _, u := range users {
		groupedUsers[u.Tokens] = append(groupedUsers[u.Tokens], u)
	}

	sortedKeys := make([]int, 0, len(groupedUsers))

	for k := range groupedUsers {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(sortedKeys)))

	return groupedUsers[sortedKeys[0]]
}

func hasSpy(user User) bool {

	if (*user.Cards.Current).Name() == "Spy" {
		return true
	}

	for _, cc := range user.Cards.Played {
		if cc.Name() == "Spy" {
			return true
		}
	}

	return false
}

func (s *Supervisor) getChancellorCards(current card.Card) []card.Card {
	var chCards []card.Card

	chCards = append(chCards, current)

	card2 := s.Game.PickCard()
	if card2 == nil {
		return chCards
	}
	chCards = append(chCards, *card2)

	card3 := s.Game.PickCard()
	if card3 == nil {
		return chCards
	}
	chCards = append(chCards, *card3)

	return chCards
}

func (s *Supervisor) MapCards(cards []card.Card) []CardInfo {
	var cardInfos []CardInfo

	for i, chCard := range cards {
		cardInfos = append(cardInfos, CardInfo{
			Value:       chCard.Value(),
			Name:        chCard.Name(),
			Description: chCard.ToString(),
			Index:       i,
		})
	}

	return cardInfos
}
