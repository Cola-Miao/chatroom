package initialize

import (
	"github.com/go-playground/validator/v10"
	"server/global"
)

func initValidator() {
	global.Validator = validator.New()
}
