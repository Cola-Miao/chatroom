package ws

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"server/global"
	"server/user"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  256,
	WriteBufferSize: 256,
}

var (
	users       = make(map[string]*user.Info)
	messageChan = make(chan []byte, 8)
)

func WebsocketHF(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}
	defer func() {
		time.Sleep(time.Second)
		conn.Close()
	}()

	u, err := userAuth(conn)
	if err != nil {
		var errMess []byte
		errMess, err = systemMessage([]byte(err.Error()))
		if err != nil {
			global.Log.Println(errMess)
		}
		u.Channel <- errMess

		return
	}

	defer func() {
		if err = userLeave(u); err != nil {
			global.Log.Println(err)
		}
	}()

	for {
		var message []byte
		message, err = readMessage(u, conn)
		if err != nil {
			writeErrMessage(u.Channel, err)
			return
		}
		messageChan <- message
	}
}
