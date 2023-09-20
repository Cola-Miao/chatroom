package main

import (
	"net/http"
	"server/global"
	"server/initialize"
	"server/router"
	ws "server/websocket"
)

func main() {
	if err := initialize.Initialize(); err != nil {
		global.Log.Fatalln(err)
	}
	go ws.Broadcast()

	http.HandleFunc("/health", router.HealthHF)
	http.HandleFunc("/ws", ws.WebsocketHF)

	global.Log.Println("Server running...")
	global.Log.Fatalln(http.ListenAndServe(":5912", nil))
}
