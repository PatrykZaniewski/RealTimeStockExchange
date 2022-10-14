package service

import (
	"broker/data_streamer/domain/model"
	"broker/data_streamer/interface/websocket"
	"log"
	"strconv"
	"time"
)

func PublishStatusOrder(orderStatus *model.OrderStatus) {
	log.Printf("%s,BROKER_DATA_STREAMER,STATUS_RECEIVED,%s", orderStatus.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	websocket.PublishOrderStatusMessage(orderStatus)
}

func PublishPrice(price *model.Price) {
	//log.Printf("%s,BROKER_DATA_STREAMER,STATUS_RECEIVED,%s", orderStatus.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	websocket.PublishPriceMessage(price)
}
