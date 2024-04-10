package transport

import (
	"errors"
	"github.com/robinpersson/LoveLetter/internal/chat"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

func Serve(supervisor *chat.Supervisor) {
	// Use websocket.Server because we want to accept non-browser clients,
	// which do not send an Origin header. websocket.Handler does check
	// the Origin header by default.
	http.Handle("/join", websocket.Server{
		Handler: supervisor.JoinWS(),
		// Set a Server.Handshake to nil - does not check the origin.
		// We can always provide a custom handshake method to access
		// the handshake http request and implement origin check or
		// other custom logic before the connection is established.
		Handshake: nil,
	})

	http.Handle("/isup", websocket.Server{
		Handler: supervisor.IsUp(),
		// Set a Server.Handshake to nil - does not check the origin.
		// We can always provide a custom handshake method to access
		// the handshake http request and implement origin check or
		// other custom logic before the connection is established.
		Handshake: nil,
	})

	http.Handle("/username", websocket.Server{
		Handler: supervisor.UserNameTaken(),
		// Set a Server.Handshake to nil - does not check the origin.
		// We can always provide a custom handshake method to access
		// the handshake http request and implement origin check or
		// other custom logic before the connection is established.
		Handshake: supervisor.UserNameTakenHandshake(),
	})

	err := http.ListenAndServe(":3000", nil)
	if !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
