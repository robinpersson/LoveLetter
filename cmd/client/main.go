package main

import (
	"github.com/robinpersson/LoveLetter/internal/frontend"
	"log"
	"os"
)

func main() {
	//myFigure := figure.NewFigure("Love Letter", "caligraphy", true)
	//myFigure.Print()
	//
	//fmt.Print("\n")
	//reader := bufio.NewReader(os.Stdin)
	//var username string
	//for {
	//	fmt.Print("Enter name: ")
	//	username, _ = reader.ReadString('\n')
	//	username = strings.Replace(username, "\n", "", -1)
	//	break
	//}

	ui, err := frontend.NewUI()
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()
	username := os.Args[1]
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
