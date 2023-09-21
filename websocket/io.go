package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"server/global"
	"server/user"
	"server/utils"
	"time"
)

func systemMessage(content []byte) (messJSON []byte, err error) {
	mess := &user.Message{
		Owner:   "Server",
		Content: content,
		Time:    time.Now(),
	}
	messJSON, err = json.Marshal(mess)

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

func writeErrMessage(channel chan []byte, err error) {
	errMess, err := systemMessage([]byte(err.Error()))
	if err != nil {
		global.Log.Println(errMess)
	}
	channel <- errMess

	return
}
