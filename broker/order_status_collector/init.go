package main

import (
	config "broker/order_status_collector/config/env"
	"broker/order_status_collector/interface/pubsub/consumer"
	"broker/order_status_collector/interface/rest"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	config.ConfigSetup()
	go consumer.InitConsumers(&wg)
	go rest.HandleRequests(&wg)
	wg.Wait()
}
