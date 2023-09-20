package initialize

import (
	"fmt"
	"log"
	"os"
	"server/global"
	"time"
)

func openLogDir() (fp *os.File, err error) {
	date := time.Now().Format("20060102")
	err = os.MkdirAll("./log/"+date, 0777)
	if err != nil {
		return
	}

	path := fmt.Sprintf("./log/%s/%s.log", date, date)
	fp, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)

	return
}

func initLog() {
	fp, err := openLogDir()
	if err != nil {
		log.Println("Log initialization failed", err)
		global.Log = log.Default()

		return
	}

	l := log.New(fp, "", log.Ldate|log.Ltime|log.Lshortfile)
	l.Println("Log initialization successful")

	global.Log = l
	return
}
