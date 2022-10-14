package service

import (
	"broker/price_collector/domain/model"
	"broker/price_collector/interface/pubsub/publisher"
)

func PublishPrices(orderStatus *model.Price) {
	//log.Printf("%s,BROKER_ORDER_STATUS_COLLECTOR,STATUS_RECEIVED,%s", orderStatus.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	publisher.PublishPrices(orderStatus)
}
