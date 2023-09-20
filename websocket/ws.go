package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"server/dao"
	"server/global"
	"server/user"
	"time"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  256,
	WriteBufferSize: 256,
}
var users = make(map[string]*websocket.Conn)
var messageChan = make(chan []byte, 8)

func auth(c *websocket.Conn) (u *user.Info, err error) {
	u = new(user.Info)
	_, userJSON, err := c.ReadMessage()
	if err != nil {
		return
	}
	if err = json.Unmarshal(userJSON, u); err != nil {
		return
	}
	if err = global.Validator.Struct(u); err != nil {
		return
	}
	if err = dao.UserAuth(u); err != nil {
		return
	}

	loginMess, err := systemMessage([]byte("Login successful"))
	if err != nil {
		global.Log.Println(err)
	}
	joinMess, err := systemMessage([]byte(u.Name + " joined the chatroom"))
	if err != nil {
		global.Log.Println(err)
	}
	messageChan <- joinMess
	if err = c.WriteMessage(websocket.TextMessage, loginMess); err != nil {
		global.Log.Println(err)
	}
	users[u.Name] = c

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

func generateMessage(owner string, content []byte) (messJSON []byte, err error) {
	mess := &user.Message{
		Owner:   owner,
		Content: content,
		Time:    time.Now(),
	}
	messJSON, err = json.Marshal(mess)

	return
}

func writeErrMessage(conn *websocket.Conn, err error) {
	errMess, err := systemMessage([]byte(err.Error()))
	if err != nil {
		global.Log.Println(errMess)
	}
	if err = conn.WriteMessage(websocket.TextMessage, errMess); err != nil {
		global.Log.Println(err)
	}

	return
}

func Broadcast() {
	for {
		select {
		case message := <-messageChan:
			for _, conn := range users {
				conn.WriteMessage(websocket.TextMessage, message)
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
	defer conn.Close()

	u, err := auth(conn)
	if err != nil {
		errMess, err := systemMessage([]byte(err.Error()))
		if err != nil {
			global.Log.Println(errMess)
		}
		if err = conn.WriteMessage(websocket.TextMessage, errMess); err != nil {
			global.Log.Println(err)
		}

		return
	}

	defer func() {
		delete(users, u.Name)
		leaveMess, err := systemMessage([]byte(u.Name + " leaves the chat room"))
		if err != nil {
			global.Log.Println(err)
			return
		}
		messageChan <- leaveMess
	}()

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			writeErrMessage(conn, err)
			return
		}

		message, err := generateMessage(u.Name, p)
		if err != nil {
			writeErrMessage(conn, err)

			return
		}
		messageChan <- message
	}
}
