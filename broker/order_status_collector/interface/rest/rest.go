package rest

import (
	config "broker/order_status_collector/config/env"
	"broker/order_status_collector/domain/model"
	"broker/order_status_collector/domain/service"
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
	http.HandleFunc("/status", status)
	log.Fatal(http.ListenAndServe(":"+generalConfig.Rest.Port, nil))
}

func status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the status!")
	body, _ := ioutil.ReadAll(r.Body)
	var pr PubSubMessage
	json.Unmarshal(body, &pr)
	var orderStatus model.OrderStatus
	json.Unmarshal(pr.Message.Data, &orderStatus)
	service.PublishOrderStatus(&orderStatus)

	w.WriteHeader(200)
	return
}
