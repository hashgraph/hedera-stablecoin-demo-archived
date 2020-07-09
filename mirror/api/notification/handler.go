package notification

import (
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// allow all origins
		// should we try and block this?
		return true
	},
}

// map of users to web sockets
var socketForUser = map[string]*websocket.Conn{}

// https://echo.labstack.com/cookbook/websocket
func Handler(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	defer ws.Close()

	// client should send us an announce
	// that identifies them

	var n struct {
		UserName string `json:"userId"`
	}

	err = ws.ReadJSON(&n)
	if err != nil {
		return err
	}

	// now, remember the connection (and keep it open) so we can reach for it when wanting to
	// send out a notification
	// TODO: multiple sockets per user?

	socketForUser[n.UserName] = ws

	// on close, let us remove ourselves from the map
	// to clean up memory and not try to ping a dead channel

	ws.SetCloseHandler(func(code int, text string) error {
		delete(socketForUser, n.UserName)

		return nil
	})

	select {}
}
