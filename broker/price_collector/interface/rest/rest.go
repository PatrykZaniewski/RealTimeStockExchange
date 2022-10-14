package rest

import (
	config "broker/price_collector/config/env"
	"broker/price_collector/domain/model"
	"broker/price_collector/domain/service"
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
	http.HandleFunc("/price", price)
	log.Fatal(http.ListenAndServe(":"+generalConfig.Rest.Port, nil))
}

func price(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the price!")
	body, _ := ioutil.ReadAll(r.Body)
	var pr PubSubMessage
	json.Unmarshal(body, &pr)
	var orderStatus model.Price
	service.PublishPrices(&orderStatus)

	w.WriteHeader(200)
	return
}
