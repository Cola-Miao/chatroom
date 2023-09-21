package ws

import (
	"errors"
	"github.com/gorilla/websocket"
	"server/dao"
	"server/global"
	"server/user"
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
	if _, ok := users[u.Name]; ok {
		err = errors.New("the user is logged in")
		return
	}
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
