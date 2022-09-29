package main

import (
	config "stock/order_collector/config/env"
	"stock/order_collector/interface/pubsub"
	"stock/order_collector/interface/rest"
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
