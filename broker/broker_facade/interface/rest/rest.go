package rest

import (
	config "broker/broker_facade/config/env"
	"broker/broker_facade/domain/model"
	"broker/broker_facade/interface/pubsub/publisher"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func HandleRequests(wg *sync.WaitGroup) {
	defer wg.Done()
	generalConfig := config.AppConfig.General
	http.HandleFunc("/", homePage)
	http.HandleFunc("/order", order)
	log.Printf("Container started at %s", generalConfig.Rest.Port)
	log.Fatal(http.ListenAndServe(":"+generalConfig.Rest.Port, nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: InternalOrder")
}

func order(w http.ResponseWriter, r *http.Request) {
	var client = r.Header.Get("identifier")
	var facadeOrder model.FacadeOrder
	json.NewDecoder(r.Body).Decode(&facadeOrder)
	if client != "mock_client" {
		log.Printf("%s,BROKER_FACADE,ORDER_RECEIVED,%s", facadeOrder.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	}
	var internalOrder = model.InternalOrder{
		AssetName:    facadeOrder.AssetName,
		Quantity:     facadeOrder.Quantity,
		OrderType:    facadeOrder.OrderType,
		OrderSubtype: facadeOrder.OrderSubtype,
		OrderPrice:   facadeOrder.OrderPrice,
		ClientId:     client,
		Id:           facadeOrder.Id,
	}

	publisher.PublishOrder(&internalOrder)
}
