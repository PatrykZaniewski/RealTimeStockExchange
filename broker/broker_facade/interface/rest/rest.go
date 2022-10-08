package rest

import (
	config "broker/broker_facade/config/env"
	"broker/broker_facade/domain/model"
	"broker/broker_facade/interface/pubsub"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

func HandleRequests(wg *sync.WaitGroup) {
	defer wg.Done()
	generalConfig := config.AppConfig.General
	http.HandleFunc("/", homePage)
	http.HandleFunc("/order", order)
	log.Fatal(http.ListenAndServe(":"+generalConfig.Rest.Port, nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	var order = model.Order{
		AssetName: "ABC",
		Quantity:  13,
	}
	pubsub.PublishOrder(&order)
	fmt.Println("Endpoint Hit: homePage")
}

func order(w http.ResponseWriter, r *http.Request) {
	message, _ := ioutil.ReadAll(r.Body)
	var order model.Order
	json.Unmarshal(message, &order)
	pubsub.PublishOrder(&order)
	fmt.Println(order)
}
