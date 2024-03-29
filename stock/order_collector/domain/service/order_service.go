package service

import (
	"log"
	"stock/order_collector/domain/model"
	"stock/order_collector/interface/pubsub/publisher"
	"strconv"
	"time"
)

func PublishOrder(order *model.StockOrder) {
	if order.BrokerId != "mock_broker" && order.ClientId != "mock_client" {
		log.Printf("%s,STOCK_ORDER_COLLECTOR,ORDER_RECEIVED,%s", order.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	}
	publisher.PublishOrder(order)
}
