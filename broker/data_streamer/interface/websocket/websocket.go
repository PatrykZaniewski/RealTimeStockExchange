package websocket

import (
	"broker/data_streamer/domain/model"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var upgrade = websocket.Upgrader{}

type ConnectionConfig struct {
	Connection      *websocket.Conn
	ConnectionMutex sync.Mutex
}

var connections = make(map[string]*ConnectionConfig)

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
		connections[string(message)] = &ConnectionConfig{
			Connection: c,
		}
	}
}

func PublishOrderStatusMessage(data *model.OrderStatus) {
	var connection = connections[data.ClientId]
	if connection != nil {
		connection.ConnectionMutex.Lock()
		defer connection.ConnectionMutex.Unlock()
		res, _ := json.Marshal(model.OrderStatusMessage{"ORDER_STATUS", *data})
		connection.Connection.WriteMessage(websocket.TextMessage, res)
		log.Printf("%s,BROKER_DATA_STREAMER,STATUS_SEND,%s", data.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	}
}

func PublishPriceMessage(data *model.Price) {
	res, _ := json.Marshal(model.PriceMessage{"PRICE", *data})
	for _, connection := range connections {
		connection.ConnectionMutex.Lock()
		defer connection.ConnectionMutex.Unlock()
		connection.Connection.WriteMessage(websocket.TextMessage, res)
	}
}
