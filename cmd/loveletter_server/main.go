package main

import (
	"github.com/robinpersson/LoveLetter/internal/chat"
	"github.com/robinpersson/LoveLetter/internal/transport"
)

func main() {
	transport.Serve(chat.NewSupervisor())
}
