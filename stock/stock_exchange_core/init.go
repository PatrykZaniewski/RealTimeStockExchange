package main

import (
	"log"
	configModel "stock/stock_exchange_core/config/model"
)

func main() {
	configModel.ConfigSetup()
	err := PublishMessage("XD", "XD", "XD")
	if err != nil {
		return
	}
	err = InitPubSubConsumer()
	if err != nil {
		log.Fatalf("PubSub init error occured: %s", err)
	}
	HandleRequests()
}
