package service

import (
	"broker/order_status_collector/domain/model"
	"broker/order_status_collector/interface/pubsub/publisher"
	"log"
	"strconv"
	"time"
)

func PublishOrderStatus(orderStatus *model.OrderStatus) {
	if orderStatus.BrokerId != "mock_broker" && orderStatus.ClientId != "mock_client" {
		log.Printf("%s,BROKER_ORDER_STATUS_COLLECTOR,STATUS_RECEIVED,%s", orderStatus.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	}
	publisher.PublishOrderStatus(orderStatus)
}
