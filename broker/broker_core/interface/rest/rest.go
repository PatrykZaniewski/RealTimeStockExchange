package rest

import (
	config "broker/broker_core/config/env"
	"broker/broker_core/domain/model"
	"broker/broker_core/domain/service"
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
	http.HandleFunc("/order", order)
	http.HandleFunc("/price", price)
	http.HandleFunc("/status", status)
	log.Fatal(http.ListenAndServe(":"+generalConfig.Rest.Port, nil))
}

func order(w http.ResponseWriter, r *http.Request) {
	log.Println("Welcome to the order!")
	//service.PublishStatusOrder()
	body, _ := ioutil.ReadAll(r.Body)
	//log.Println(string(body))
	var pr PubSubMessage
	json.Unmarshal(body, &pr)
	//log.Println(pr)
	var order model.InternalOrder
	json.Unmarshal(pr.Message.Data, &order)
	service.ProcessOrder(&order)

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

	w.WriteHeader(200)
	return
}

func status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the status!")
	body, _ := ioutil.ReadAll(r.Body)
	var pr PubSubMessage
	json.Unmarshal(body, &pr)
	var status model.OrderStatus
	json.Unmarshal(pr.Message.Data, &status)
	service.PublishStatusOrder(&status)

	w.WriteHeader(200)
	return
}
