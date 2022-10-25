package rest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	config "stock/order_collector/config/env"
	"stock/order_collector/domain/model"
	"stock/order_collector/domain/service"
	"sync"
)

type PubSubMessage struct {
	Message struct {
		Data []byte `json:"data,omitempty"`
		ID   string `json:"id"`
	} `json:"message"`
	Subscription string `json:"subscription"`
}

func HandleRequests(wg *sync.WaitGroup) {
	defer wg.Done()
	generalConfig := config.AppConfig.General
	http.HandleFunc("/order", order)
	log.Fatal(http.ListenAndServe(":"+generalConfig.Rest.Port, nil))
}

func order(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the order!")
	body, _ := ioutil.ReadAll(r.Body)
	var pr PubSubMessage
	json.Unmarshal(body, &pr)
	var stockOrder model.StockOrder
	json.Unmarshal(pr.Message.Data, &stockOrder)
	service.PublishOrder(&stockOrder)
}
