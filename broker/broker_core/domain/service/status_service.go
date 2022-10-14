package service

import (
	"broker/broker_core/domain/model"
	"broker/broker_core/interface/pubsub/publisher"
	"log"
	"strconv"
	"time"
)

func PublishStatusOrder(orderStatus *model.OrderStatus) {
	log.Printf("%s,BROKER_CORE,STATUS_RECEIVED,%s", orderStatus.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	publisher.PublishOrderStatus(orderStatus)
}
