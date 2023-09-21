package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func Encode(text []byte) (encodeText []byte, err error) {
	encodeText, err = bcrypt.GenerateFromPassword(text, bcrypt.DefaultCost)

	return
}

func ComparePwd(hash []byte, pwd []byte) (err error) {
	err = bcrypt.CompareHashAndPassword(hash, pwd)

	return
}
