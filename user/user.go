package user

import (
	"gorm.io/gorm"
	"time"
)

type Info struct {
	gorm.Model
	Name     string `validate:"min=1,max=10" gorm:"index"`
	Password string `validate:"required"`
}

type Message struct {
	Owner   string
	Content []byte
	Time    time.Time
}
