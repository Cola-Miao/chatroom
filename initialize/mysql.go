package initialize

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"server/global"
	"server/user"
)

type mysqlConfig struct {
	User   string
	Pass   string
	DBName string
	IP     string
	Port   int
}

func initMysql() (err error) {
	var mc mysqlConfig
	err = global.Viper.UnmarshalKey("mysql", &mc)
	if err != nil {
		return
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		mc.User, mc.Pass, mc.IP, mc.Port, mc.DBName)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return
	}
	err = db.AutoMigrate(&user.Info{})
	if err != nil {
		return
	}
	global.Log.Println("Mysql connection successful")

	global.DB = db
	return
}
