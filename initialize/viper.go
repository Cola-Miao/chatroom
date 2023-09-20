package initialize

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"server/global"
)

func initViper() (err error) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./config")

	err = v.ReadInConfig()
	if err != nil {
		return
	}

	v.OnConfigChange(func(in fsnotify.Event) {
		global.Viper = v
		Initialize()
		global.Log.Println("The configuration file has changed")
	})
	v.WatchConfig()
	global.Log.Println("Configuration read successful")

	global.Viper = v
	return
}
