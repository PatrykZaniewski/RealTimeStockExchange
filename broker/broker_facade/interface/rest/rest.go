package rest

import (
	config "broker/broker_facade/config/env"
	"broker/broker_facade/domain/model"
	"broker/broker_facade/interface/pubsub"
	"encoding/json"
	"fmt"
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
	fmt.Println("Endpoint Hit: InternalOrder")
}

func order(w http.ResponseWriter, r *http.Request) {
	//message, _ := ioutil.ReadAll(r.Body)
	//var order model.FacadeOrder
	//json.Unmarshal(message, &order)
	//pubsub.PublishOrder(&order)
	//fmt.Println(order)
	fmt.Println("Endpoint Hit: InternalOrder")
	var client = r.Header.Get("identifier")
	var facadeOrder model.FacadeOrder
	json.NewDecoder(r.Body).Decode(&facadeOrder)
	var internalOrder = model.InternalOrder{
		AssetName:    facadeOrder.AssetName,
		Quantity:     facadeOrder.Quantity,
		OrderType:    facadeOrder.OrderType,
		OrderSubtype: facadeOrder.OrderSubtype,
		ClientId:     client,
		Id:           facadeOrder.Id,
	}

	pubsub.PublishOrder(&internalOrder)
	fmt.Println(internalOrder)
}
