package service

import (
	"log"
	"math"
	"math/rand"
	"stock/stock_exchange_core/domain/model"
	"stock/stock_exchange_core/interface/pubsub/publisher"
	"strconv"
	"time"
)

func ProcessOrder(stockOrder *model.StockOrder) {
	log.Printf("%s,STOCK_CORE,ORDER_RECEIVED,%s", stockOrder.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))

	var orderStatus = model.OrderStatus{
		AssetName:    stockOrder.AssetName,
		Quantity:     stockOrder.Quantity,
		OrderType:    stockOrder.OrderType,
		OrderSubtype: stockOrder.OrderSubtype,
		OrderPrice:   stockOrder.OrderPrice,
		ClientId:     stockOrder.ClientId,
		BrokerId:     stockOrder.BrokerId,
		Id:           stockOrder.Id,
		Status:       model.FULFILLED,
	}

	publisher.PublishOrderStatus(&orderStatus)

	var price = model.Price{
		AssetName: stockOrder.AssetName,
		BuyPrice:  float32(math.Round((rand.Float64()*100+600)*100) / 100),
		SellPrice: float32(math.Round((rand.Float64()*100+600)*100) / 100),
	}
	publisher.PublishPrices(&price)
}
