package pubsub

import (
	config "broker/broker_core/config/env"
	"broker/broker_core/domain/model"
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

func InitConsumers(wg *sync.WaitGroup) error {
	defer wg.Done()
	initBrokerConsumer()
	return nil
}

func initBrokerConsumer() error {
	brokersPubSubConfig := config.AppConfig.PubSub.Broker
	initConsumer(brokersPubSubConfig.ProjectId, brokersPubSubConfig.Consumer.BrokerInternalClientOrdersSubId, ordersCallback)
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
	var order model.InternalOrder
	json.Unmarshal(msg.Data, &order)
	PublishOrder(&order)
	fmt.Printf("Got message: %q\n\n", string(msg.Data))
	msg.Ack()
}
