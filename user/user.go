package user

import (
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
	"time"
)

type Info struct {
	gorm.Model
	Name     string `validate:"min=1,max=10" gorm:"index"`
	Password string `validate:"required"`
}

type Message struct {
	Owner   string    `json:"owner,omitempty"`
	Content []byte    `json:"content,omitempty"`
	Time    time.Time `json:"time"`
}

type PrivateMessage struct {
	Content []byte
	Conn    *websocket.Conn
}
