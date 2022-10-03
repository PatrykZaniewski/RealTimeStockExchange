package pubsub

import (
	config "broker/price_streamer/config/env"
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"sync"
)

func ordersCallback(_ context.Context, msg *pubsub.Message) {
	fmt.Printf("Got message: %q\n\n", string(msg.Data))
	msg.Ack()
}

func initOrdersConsumer() error {
	pubSubConfig := config.AppConfig.PubSub
	projectId := pubSubConfig.Broker.ProjectId
	subscriptionId := pubSubConfig.Broker.Consumer.BrokerInternalPricesSubId

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
