package main

import (
	config "stock/order_collector/config/env"
	"stock/order_collector/interface/pubsub/consumer"
	"stock/order_collector/interface/rest"
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
