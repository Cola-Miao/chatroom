package dao

import (
	"errors"
	"gorm.io/gorm"
	"server/global"
	"server/user"
	"server/utils"
)

func UserAuth(u *user.Info) (err error) {
	var result user.Info
	if err = global.DB.Model(&user.Info{}).Where("name = ?", u.Name).First(&result).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			var encodePwd []byte
			encodePwd, err = utils.Encode([]byte(u.Password))
			if err != nil {
				global.Log.Println(err)
				return
			}
			u.Password = string(encodePwd)
			err = global.DB.Model(&user.Info{}).Create(u).Error
		}
	} else {
		if err = utils.ComparePwd([]byte(result.Password), []byte(u.Password)); err != nil {
			err = errors.New("wrong password")
		}
	}

	return
}
