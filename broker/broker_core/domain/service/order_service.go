package service

import (
	"broker/broker_core/domain/model"
	"broker/broker_core/interface/pubsub/publisher"
	"log"
	"strconv"
	"time"
)

func PublishOrder(internalOrder *model.InternalOrder) {
	log.Printf("%s,BROKER_CORE,ORDER_RECEIVED,%s", internalOrder.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	publisher.PublishOrder(internalOrder)
}
