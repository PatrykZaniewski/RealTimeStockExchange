package main

import (
	config "broker/price_streamer/config/env"
	"broker/price_streamer/interface/pubsub"
	"broker/price_streamer/interface/rest"
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
