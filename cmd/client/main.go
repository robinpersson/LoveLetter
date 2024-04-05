package main

import (
	"bufio"
	"fmt"
	"github.com/common-nighthawk/go-figure"
	"github.com/robinpersson/LoveLetter/internal/frontend"
	"log"
	"os"
	"strings"
)

func main() {
	username := ""
	if len(os.Args) == 2 {
		username = os.Args[1]
	} else {
		myFigure := figure.NewFigure("Love Letter", "larry3d", true)
		myFigure.Print()

		fmt.Print("\n")
		reader := bufio.NewReader(os.Stdin)
		var username string
		for {
			fmt.Print("Enter name: ")
			username, _ = reader.ReadString('\n')
			username = strings.Replace(username, "\n", "", -1)
			break
		}
	}

	ui, err := frontend.NewUI()
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()
	ui.SetUsername(username)

	if err = ui.Connect(username); err != nil {
		log.Fatal(err)
	}

	ui.SetManagerFunc(ui.Layout)
	ui.SetKeyBindings(ui.Gui)

	if err = ui.Serve(); err != nil {
		log.Fatal(err)
	}
}
