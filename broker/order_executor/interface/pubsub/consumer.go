package pubsub

import (
	config "broker/order_executor/config/env"
	"broker/order_executor/domain/model"
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"
)

func InitConsumers(wg *sync.WaitGroup) error {
	defer wg.Done()
	initBrokerConsumer()
	return nil
}

func initBrokerConsumer() error {
	brokersPubSubConfig := config.AppConfig.PubSub.Broker
	initConsumer(brokersPubSubConfig.ProjectId, brokersPubSubConfig.Consumer.BrokerInternalOrdersSubId, ordersCallback)
	return nil
}

func initConsumer(projectId, subId string, callback func(context.Context, *pubsub.Message)) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		_ = fmt.Errorf("PubSub connection error: %v", err)
	}
	defer client.Close()

	sub := client.Subscription(subId)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	err = sub.Receive(ctx, callback)
	if err != nil {
		return fmt.Errorf("sub.Receive: %v", err)
	}

	return nil
}

func ordersCallback(_ context.Context, msg *pubsub.Message) {
	appConfig := config.AppConfig

	var internalOrder model.InternalOrder
	json.Unmarshal(msg.Data, &internalOrder)
	log.Printf("%s,BROKER_ORDER_EXECUTOR,ORDER_RECEIVED,%s", internalOrder.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
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
	PublishOrder(&stockOrder)
	fmt.Printf("Got message: %q\n\n", string(msg.Data))
	msg.Ack()
}
