package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"server/dao"
	"server/global"
	"server/user"
	"server/utils"
	"strings"
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

func userAuth(c *websocket.Conn) (u *user.Info, err error) {
	u = new(user.Info)
	_, userJSON, err := c.ReadMessage()
	if err != nil {
		return
	}
	if err = user.UnmarshalUser(userJSON, u, c); err != nil {
		return
	}
	go privateBroadcast(u)
	if err = global.Validator.Struct(u); err != nil {
		return
	}
	if err = dao.UserAuth(u); err != nil {
		return
	}
	if err = authSuccess(u); err != nil {
		return
	}

	return
}

func authSuccess(u *user.Info) (err error) {
	loginMess, err := systemMessage([]byte("Login successful"))
	if err != nil {
		return
	}
	joinMess, err := systemMessage([]byte(u.Name + " joined the chatroom"))
	if err != nil {
		return
	}
	messageChan <- joinMess
	u.Channel <- loginMess
	users[u.Name] = u

	ul := usersList()
	listMess, err := systemMessage(ul)
	if err != nil {
		return
	}
	u.Channel <- listMess

	return
}

func usersList() (content []byte) {
	var builder strings.Builder
	builder.Write([]byte("UsersList: "))
	for u := range users {
		builder.Write([]byte(u))
		builder.WriteByte(',')
	}

	content = []byte(builder.String())

	return
}

func systemMessage(content []byte) (messJSON []byte, err error) {
	mess := &user.Message{
		Owner:   "Server",
		Content: content,
		Time:    time.Now(),
	}
	messJSON, err = json.Marshal(mess)

	return
}

func writeErrMessage(channel chan []byte, err error) {
	errMess, err := systemMessage([]byte(err.Error()))
	if err != nil {
		global.Log.Println(errMess)
	}
	channel <- errMess

	return
}

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

func userLeave(u *user.Info) (err error) {
	delete(users, u.Name)
	var leaveMess []byte
	leaveMess, err = systemMessage([]byte(u.Name + " leaves the chat room"))
	if err != nil {
		return
	}
	messageChan <- leaveMess

	return
}

func readMessage(u *user.Info, c *websocket.Conn) (message []byte, err error) {
	var p []byte
	_, p, err = c.ReadMessage()
	if err != nil {
		writeErrMessage(u.Channel, err)
		return
	}
	message, err = utils.GenerateMessage(u.Name, p)

	return
}
