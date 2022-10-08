package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{}

var tmp []*websocket.Conn

func Websocket(w http.ResponseWriter, r *http.Request) {
	fmt.Print("XD")
	c, err := upgrader.Upgrade(w, r, nil)
	fmt.Print("XD2")
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	fmt.Print("XD3")
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
		d := map[string]string{
			"buy_price":  "635.74",
			"name":       "CDPROJECT",
			"sell_price": "539.34",
		}
		res, _ := json.Marshal(d)
		tmp[0].WriteMessage(websocket.TextMessage, res)

	}
}
