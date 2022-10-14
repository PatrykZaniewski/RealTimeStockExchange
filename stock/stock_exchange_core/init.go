package main

import (
	config "stock/stock_exchange_core/config/env"
	"stock/stock_exchange_core/interface/database"
	"stock/stock_exchange_core/interface/pubsub/consumer"
	"stock/stock_exchange_core/interface/rest"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	config.ConfigSetup()
	go consumer.InitConsumers(&wg)
	go rest.HandleRequests(&wg)
	database.ProceedDbOperation()
	wg.Wait()
}
