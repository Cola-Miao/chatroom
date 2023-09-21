package ws

import (
	"github.com/gorilla/websocket"
	"server/user"
)

func privateBroadcast(u *user.Info) {
	for {
		select {
		case message := <-u.Channel:
			if err := u.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				return
			}
		}
	}
}

func Broadcast() {
	for {
		select {
		case message := <-messageChan:
			for _, u := range users {
				u.Channel <- message
			}
		}
	}
}
