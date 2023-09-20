package initialize

func Initialize() (err error) {
	initLog()
	if err = initViper(); err != nil {
		return
	}
	if err = initMysql(); err != nil {
		return
	}
	initValidator()

	return
}
