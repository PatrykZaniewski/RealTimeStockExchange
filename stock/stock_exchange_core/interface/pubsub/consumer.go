package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	config "stock/stock_exchange_core/config/env"
	"stock/stock_exchange_core/domain/model"
	"sync"
	"sync/atomic"
)

func ordersCallback(_ context.Context, msg *pubsub.Message) {
	var received int32
	var order model.Order
	json.Unmarshal([]byte(msg.Data), &order)
	fmt.Printf("Got message: %q\n\n", string(msg.Data))
	atomic.AddInt32(&received, 1)
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
