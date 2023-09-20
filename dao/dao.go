package dao

import (
	"errors"
	"gorm.io/gorm"
	"server/global"
	"server/user"
)

func UserAuth(u *user.Info) (err error) {
	var result user.Info
	if err = global.DB.Model(&user.Info{}).Where("name = ?", u.Name).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = global.DB.Model(&user.Info{}).Create(u).Error
		}
	} else {
		if u.Password != result.Password {
			err = errors.New("wrong password")
		}
	}

	return
}
