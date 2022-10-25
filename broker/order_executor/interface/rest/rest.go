package rest

import (
	config "broker/order_executor/config/env"
	"broker/order_executor/domain/model"
	"broker/order_executor/domain/service"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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
	body, _ := ioutil.ReadAll(r.Body)
	var pr PubSubMessage
	json.Unmarshal(body, &pr)
	var internalOrder model.InternalOrder
	json.Unmarshal(pr.Message.Data, &internalOrder)
	service.PublishOrder(&internalOrder)
}
