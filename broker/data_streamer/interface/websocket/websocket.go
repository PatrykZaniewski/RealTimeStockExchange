package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{}

var tmp []*websocket.Conn

func Websocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	tmp = append(tmp, c)
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, []byte("Witam"))
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func PublishMessage(data string) {
	if len(tmp) > 0 {
		tmp[0].WriteMessage(2, []byte("XD"))

	}
}
