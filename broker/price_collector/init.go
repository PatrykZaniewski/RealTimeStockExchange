package main

import (
	config "broker/price_collector/config/env"
	"broker/price_collector/interface/pubsub/consumer"
	"broker/price_collector/interface/rest"
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
