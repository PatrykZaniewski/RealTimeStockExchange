package pubsub

import (
	config "broker/price_collector/config/env"
	"broker/price_collector/domain/model"
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
	PublishPrices(&price)
	msg.Ack()
}

func initOrdersConsumer() error {
	pubSubConfig := config.AppConfig.PubSub
	projectId := pubSubConfig.Stock.ProjectId
	subscriptionId := pubSubConfig.Stock.Consumer.CorePricesSubId

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		_ = fmt.Errorf("PubSub connection error: %v", err)
	}
	defer client.Close()

	sub := client.Subscription(subscriptionId)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err = sub.Receive(ctx, pricesCallback)
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
