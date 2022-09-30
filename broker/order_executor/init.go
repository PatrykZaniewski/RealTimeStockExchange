package main

import (
	config "broker/order_executor/config/env"
	"broker/order_executor/interface/pubsub"
	"broker/order_executor/interface/rest"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	config.ConfigSetup()
	go pubsub.InitConsumers(&wg)
	go rest.HandleRequests(&wg)
	wg.Wait()
}
