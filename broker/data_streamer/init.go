package main

import (
	config "broker/data_streamer/config/env"
	"broker/data_streamer/interface/pubsub/consumer"
	"broker/data_streamer/interface/rest"
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
