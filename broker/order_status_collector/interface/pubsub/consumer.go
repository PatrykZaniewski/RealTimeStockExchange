package pubsub

import (
	config "broker/order_status_collector/config/env"
	"broker/order_status_collector/domain/model"
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

func ordersStatusCallback(_ context.Context, msg *pubsub.Message) {
	fmt.Printf("Got message: %q\n\n", string(msg.Data))
	var orderStatus model.OrderStatus
	json.Unmarshal(msg.Data, &orderStatus)
	PublishOrderStatus(&orderStatus)
	msg.Ack()
}

func initOrdersStatusConsumer() error {
	pubSubConfig := config.AppConfig.PubSub
	projectId := pubSubConfig.Broker.ProjectId
	subscriptionId := pubSubConfig.Broker.Consumer.BrokerOrdersSubId

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		_ = fmt.Errorf("PubSub connection error: %v", err)
	}
	defer client.Close()

	sub := client.Subscription(subscriptionId)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err = sub.Receive(ctx, ordersStatusCallback)
	if err != nil {
		return fmt.Errorf("sub.Receive: %v", err)
	}

	return nil
}

func InitConsumers(wg *sync.WaitGroup) error {
	defer wg.Done()
	err := initOrdersStatusConsumer()
	if err != nil {
		return err
	}
	return nil
}
