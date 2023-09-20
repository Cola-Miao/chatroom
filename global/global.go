package global

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"log"
)

var (
	DB        *gorm.DB
	Log       *log.Logger
	Viper     *viper.Viper
	Validator *validator.Validate
)
