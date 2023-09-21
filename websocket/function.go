package ws

import (
	"server/user"
	"strings"
)

func userLeave(u *user.Info) (err error) {
	delete(users, u.Name)
	var leaveMess []byte
	leaveMess, err = systemMessage([]byte(u.Name + " leaves the chatroom"))
	if err != nil {
		return
	}
	messageChan <- leaveMess

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
