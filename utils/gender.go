package utils

import (
	"encoding/json"
	"server/user"
	"time"
)

func GenerateMessage(owner string, content []byte) (messJSON []byte, err error) {
	mess := &user.Message{
		Owner:   owner,
		Content: content,
		Time:    time.Now(),
	}
	messJSON, err = json.Marshal(mess)

	return
}
