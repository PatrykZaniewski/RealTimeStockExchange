package main

import (
	config "stock/stock_exchange_core/config/env"
	"stock/stock_exchange_core/interface/database"
)

func main() {
	//var wg sync.WaitGroup
	//wg.Add(2)
	config.ConfigSetup()
	//go pubsub.InitConsumers(&wg)
	//go rest.HandleRequests(&wg)
	database.DatabaseOperation()
	//wg.Wait()
}
