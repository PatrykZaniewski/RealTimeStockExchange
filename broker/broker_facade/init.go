package main

import (
	config "broker/broker_facade/config/env"
	"broker/broker_facade/interface/rest"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	config.ConfigSetup()
	go rest.HandleRequests(&wg)
	wg.Wait()
}
