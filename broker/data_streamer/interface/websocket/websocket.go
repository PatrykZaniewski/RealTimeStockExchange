package websocket

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrade = websocket.Upgrader{}

var connections = make(map[string]*websocket.Conn)

func Websocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("error:", err)
			break
		}
		log.Printf("Register: %s", message)
		connections[string(message)] = c
	}
}

func PublishMessage(data map[string]string) {
	res, _ := json.Marshal(data)
	connections[data["identifier"]].WriteMessage(websocket.TextMessage, res)
}
