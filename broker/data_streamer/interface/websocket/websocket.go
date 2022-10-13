package websocket

import (
	"broker/data_streamer/domain/model"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"time"
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
		res, _ := json.Marshal(model.OrderStatusMessage{"ORDER_STATUS", *data})
		connection.WriteMessage(websocket.TextMessage, res)
		log.Printf("%s,BROKER_DATA_STREAMER,STATUS_SEND,%s", data.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	}
}

func PublishPriceMessage(data *model.Price) {
	res, _ := json.Marshal(model.PriceMessage{"PRICE", *data})
	for _, connection := range connections {
		connection.WriteMessage(websocket.TextMessage, res)
	}
}
