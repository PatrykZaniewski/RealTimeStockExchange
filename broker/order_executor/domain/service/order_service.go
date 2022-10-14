package service

import (
	config "broker/order_executor/config/env"
	"broker/order_executor/domain/model"
	"broker/order_executor/interface/pubsub/publisher"
	"log"
	"strconv"
	"time"
)

func PublishOrder(internalOrder *model.InternalOrder) {
	log.Printf("%s,BROKER_ORDER_EXECUTOR,ORDER_RECEIVED,%s", internalOrder.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))

	appConfig := config.AppConfig
	var stockOrder = model.StockOrder{
		AssetName:    internalOrder.AssetName,
		Quantity:     internalOrder.Quantity,
		OrderType:    internalOrder.OrderType,
		OrderSubtype: internalOrder.OrderSubtype,
		OrderPrice:   internalOrder.OrderPrice,
		ClientId:     internalOrder.ClientId,
		BrokerId:     appConfig.Identity,
		Id:           internalOrder.Id,
	}
	publisher.PublishOrder(&stockOrder)
}
