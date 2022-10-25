package rest

import (
	config "broker/data_streamer/config/env"
	"broker/data_streamer/domain/model"
	"broker/data_streamer/domain/service"
	wb "broker/data_streamer/interface/websocket"
	"encoding/json"
	"fmt"
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
	http.HandleFunc("/ws", wb.Websocket)
	http.HandleFunc("/status", status)
	http.HandleFunc("/price", price)
	log.Fatal(http.ListenAndServe(":"+generalConfig.Rest.Port, nil))
}

func status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the price!")
	body, _ := ioutil.ReadAll(r.Body)
	var pr PubSubMessage
	json.Unmarshal(body, &pr)
	var orderStatus model.OrderStatus
	json.Unmarshal(pr.Message.Data, &orderStatus)
	service.PublishStatusOrder(&orderStatus)

	w.WriteHeader(200)
	return
}

func price(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the price!")
	body, _ := ioutil.ReadAll(r.Body)
	var pr PubSubMessage
	json.Unmarshal(body, &pr)
	var price model.Price
	json.Unmarshal(pr.Message.Data, &price)
	service.PublishPrice(&price)

	w.WriteHeader(200)
	return
}
