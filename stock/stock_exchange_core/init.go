package main

import (
	config "stock/stock_exchange_core/config/env"
	"stock/stock_exchange_core/domain/service"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	config.ConfigSetup()
	service.ProcessMarketOrder()
	//go consumer.InitConsumers(&wg)
	//go rest.HandleRequests(&wg)
	//database.ProceedDbOperation()
	//wg.Wait()
}
