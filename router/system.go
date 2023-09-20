package router

import (
	"fmt"
	"net/http"
	"server/global"
)

func HealthHF(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintln(w, "OK")
	if err != nil {
		global.Log.Println(err)
	}

	return
}
