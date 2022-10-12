package pubsub

import (
	config "broker/data_streamer/config/env"
	"broker/data_streamer/domain/model"
	"broker/data_streamer/interface/websocket"
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

func pricesCallback(_ context.Context, msg *pubsub.Message) {
	fmt.Printf("Got message: %q\n\n", string(msg.Data))
	var price model.Price
	json.Unmarshal(msg.Data, &price)
	websocket.PublishPriceMessage(&price)
	msg.Ack()
}

func orderStatusCallback(_ context.Context, msg *pubsub.Message) {
	fmt.Printf("Got message: %q\n\n", string(msg.Data))
	var orderStatus model.OrderStatus
	json.Unmarshal(msg.Data, &orderStatus)
	websocket.PublishOrderStatusMessage(&orderStatus)
	msg.Ack()
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

func initPriceConsumer() error {
	brokersPubSubConfig := config.AppConfig.PubSub.Broker
	initConsumer(brokersPubSubConfig.ProjectId, brokersPubSubConfig.Consumer.BrokerInternalPricesSubId, pricesCallback)
	return nil
}

func initOrderStatusConsumer() error {
	brokersPubSubConfig := config.AppConfig.PubSub.Broker
	initConsumer(brokersPubSubConfig.ProjectId, brokersPubSubConfig.Consumer.BrokerInternalCoreOrdersStatusSubId, orderStatusCallback)
	return nil
}

func InitConsumers(wg *sync.WaitGroup) error {
	defer wg.Done()
	initOrderStatusConsumer()
	initPriceConsumer()
	return nil
}
