package websocket

import (
	"broker/data_streamer/domain/model"
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

func PublishOrderStatusMessage(data *model.OrderStatus) {
	var connection = connections[data.ClientId]
	if connection != nil {
		res, _ := json.Marshal(data)
		connection.WriteMessage(websocket.TextMessage, res)
	}
}

func PublishPriceMessage(data *model.Price) {
	res, _ := json.Marshal(data)
	for _, connection := range connections {
		connection.WriteMessage(websocket.TextMessage, res)
	}
}
