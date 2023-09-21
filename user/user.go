package user

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"time"
)

type Info struct {
	gorm.Model
	Name     string          `validate:"min=1,max=10" gorm:"index"`
	Password string          `validate:"required"`
	Channel  chan []byte     `gorm:"-"`
	Conn     *websocket.Conn `gorm:"-"`
}

type Message struct {
	Owner   string    `json:"owner,omitempty"`
	Content []byte    `json:"content,omitempty"`
	Time    time.Time `json:"time"`
}

func UnmarshalUser(userJSON []byte, u *Info, c *websocket.Conn) (err error) {
	err = json.Unmarshal(userJSON, u)
	if err != nil {
		return
	}
	u.Channel = make(chan []byte, 1)
	u.Conn = c

	return
}
