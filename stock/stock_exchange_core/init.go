package main

import (
	config "stock/stock_exchange_core/config/env"
	"stock/stock_exchange_core/interface/pubsub/consumer"
	"stock/stock_exchange_core/interface/rest"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	config.ConfigSetup()
	//for i := 1; i <= 40; i++ {
	//	service.ProcessLimitOrder()
	//}
	//service.ProcessMarketOrder()
	go consumer.InitConsumers(&wg)
	go rest.HandleRequests(&wg)
	//database.ProceedDbOperation()
	wg.Wait()

}
