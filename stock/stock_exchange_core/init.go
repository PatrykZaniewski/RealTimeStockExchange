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
	//var stockOrder = &model.StockOrder{
	//	AssetName:    "ASSECO",
	//	Quantity:     1,
	//	OrderType:    "SELL",
	//	OrderSubtype: "MARKET_ORDER",
	//	OrderPrice:   205.00,
	//	ClientId:     "broker_client",
	//	BrokerId:     "mock_broker",
	//	Id:           uuid.New().String(),
	//}
	//for i := 1; i <= 40; i++ {
	//	service.ProcessLimitOrder(stockOrder)
	//}
	//service.ProcessMarketOrder()
	go consumer.InitConsumers(&wg)
	go rest.HandleRequests(&wg)
	//database.ProceedDbOperation()
	wg.Wait()

}
