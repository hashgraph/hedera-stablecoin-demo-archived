package notification

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
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

func Handler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Err(err).Msg("cannot upgrade to ws://")
		return
	}

	// client should send us an announce
	// that identifies them

	var n struct {
		UserName string `json:"userId"`
	}

	err = conn.ReadJSON(&n)
	if err != nil {
		if err != websocket.ErrCloseSent {
			log.Err(err).Msg("cannot read message from ws")
		}

		return
	}

	// now, remember the connection (and keep it open) so we can reach for it when wanting to
	// send out a notification
	// TODO: multiple sockets per user?

	socketForUser[n.UserName] = conn

	// on close, let us remove ourselves from the map
	// to clean up memory and not try to ping a dead channel

	conn.SetCloseHandler(func(code int, text string) error {
		delete(socketForUser, n.UserName)

		return nil
	})

	select {}
}
