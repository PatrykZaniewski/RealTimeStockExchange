package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	config "stock/stock_exchange_core/config/env"
	"stock/stock_exchange_core/domain/model"
	"strconv"
	"sync"
	"time"
)

func ordersCallback(_ context.Context, msg *pubsub.Message) {
	fmt.Printf("Got message: %q\n\n", string(msg.Data))
	var stockOrder model.StockOrder
	json.Unmarshal(msg.Data, &stockOrder)
	log.Printf("%s,STOCK_CORE,ORDER_RECEIVED,%s", stockOrder.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))

	var price = model.Price{
		AssetName: stockOrder.AssetName,
		BuyPrice:  float32(math.Round((rand.Float64()*100+600)*100) / 100),
		SellPrice: float32(math.Round((rand.Float64()*100+600)*100) / 100),
	}
	PublishPrices(&price)

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
	PublishOrderStatus(&orderStatus)

	msg.Ack()
}

func initOrdersConsumer() error {
	pubSubConfig := config.AppConfig.PubSub
	projectId := pubSubConfig.Stock.ProjectId
	subscriptionId := pubSubConfig.Stock.Consumer.InternalOrdersSubId

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		_ = fmt.Errorf("PubSub connection error: %v", err)
	}
	defer client.Close()

	sub := client.Subscription(subscriptionId)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err = sub.Receive(ctx, ordersCallback)
	if err != nil {
		return fmt.Errorf("sub.Receive: %v", err)
	}

	return nil
}

func InitConsumers(wg *sync.WaitGroup) error {
	defer wg.Done()
	err := initOrdersConsumer()
	if err != nil {
		return err
	}
	return nil
}
