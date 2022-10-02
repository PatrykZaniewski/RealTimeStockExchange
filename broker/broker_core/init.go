package main

import (
	config "broker/broker_core/config/env"
	"broker/broker_core/interface/pubsub"
	"broker/broker_core/interface/rest"
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
