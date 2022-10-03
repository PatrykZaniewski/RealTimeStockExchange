package rest

import (
	config "broker/broker_facade/config/env"
	"broker/broker_facade/domain/model"
	"broker/broker_facade/interface/pubsub"
	"fmt"
	"log"
	"net/http"
	"sync"
)

func HandleRequests(wg *sync.WaitGroup) {
	defer wg.Done()
	generalConfig := config.AppConfig.General
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":"+generalConfig.Rest.Port, nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	var order = model.Order{
		AssetName: "ABC",
		Quantity:  13.3,
	}
	pubsub.PublishOrder(&order)
	fmt.Println("Endpoint Hit: homePage")
}
