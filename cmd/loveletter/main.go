package main

import (
	"bufio"
	"embed"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/hajimehoshi/go-mp3"
	"github.com/hajimehoshi/oto"
	"github.com/robinpersson/LoveLetter/internal/frontend"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"os"
	"strings"
)

//go:embed resources/intro.mp3
var embedFS embed.FS

func main() {
	fmt.Print("\033[H\033[2J")
	go PlayIntro()
	username := ""
	address := ""

	if len(os.Args) == 2 {
		username = os.Args[1]
	} else {
		myFigure := figure.NewFigure("Love Letter", "larry3d", true)
		myFigure.Print()

		fmt.Print("\n")
		reader := bufio.NewReader(os.Stdin)

		for {
			fmt.Print("Enter address: ")
			address, _ = reader.ReadString('\n')
			address = strings.Replace(address, "\n", "", -1)
			if !isServerUp(address) {
				fmt.Printf("Unable to connect to %s\n", address)
			} else {
				break
			}
		}

		fmt.Print("\n")
		reader = bufio.NewReader(os.Stdin)
		for {
			fmt.Print("Enter name: ")
			username, _ = reader.ReadString('\n')
			username = strings.Replace(username, "\n", "", -1)

			if len(username) > 0 {
				taken := isUserNameTaken(address, username)
				if !taken {
					break
				}
				fmt.Print("Name already taken\n")
			} else {
				fmt.Print("Name cannot be empty\n")
			}

		}

	}

	ui, err := frontend.NewUI()
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()

	ui.Connect(username, address)

	ui.SetUsername(username)
	ui.SetManagerFunc(ui.Layout)
	ui.SetKeyBindings(ui.Gui)

	if err = ui.Serve(); err != nil {
		log.Fatal(err)
	}
}

func isServerUp(address string) bool {
	config, err := websocket.NewConfig(fmt.Sprintf("ws://%s:3000/isup", address), frontend.WebsocketOrigin)
	if err != nil {
		return false
	}

	_, err = websocket.DialConfig(config)
	if err != nil {
		return false
	}

	return true
}

func isUserNameTaken(address, userName string) bool {
	config, err := websocket.NewConfig(fmt.Sprintf("ws://%s:3000/username", address), frontend.WebsocketOrigin)
	config.Header.Set("Username", userName)

	if err != nil {
		return true
	}

	_, err = websocket.DialConfig(config)

	if err != nil {
		return true
	}

	return false
}

func PlayIntro() error {

	f, err := embedFS.Open("resources/intro.mp3")
	if err != nil {
		return err
	}
	defer f.Close()

	d, err := mp3.NewDecoder(f)
	if err != nil {
		return err
	}

	c, err := oto.NewContext(d.SampleRate(), 2, 2, 8192)
	if err != nil {
		return err
	}
	defer c.Close()

	p := c.NewPlayer()

	defer p.Close()

	//fmt.Printf("Length: %d[bytes]\n", d.Length())

	if _, err := io.Copy(p, d); err != nil {
		return err
	}
	return nil
}
