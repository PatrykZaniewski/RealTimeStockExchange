package consumer

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	config "stock/order_collector/config/env"
	"stock/order_collector/domain/model"
	"stock/order_collector/domain/service"
	"sync"
)

func InitConsumers(wg *sync.WaitGroup) error {
	defer wg.Done()
	initBrokerConsumers()
	initMockOrdersConsumer()
	return nil
}

func initBrokerConsumers() error {
	brokersPubSubConfigs := config.AppConfig.PubSub.Broker
	for _, brokerPubSubConfig := range brokersPubSubConfigs {
		initConsumer(brokerPubSubConfig.ProjectId, brokerPubSubConfig.Consumer.BrokerPendingOrdersSubId, ordersCallback)
	}
	return nil
}

func initMockOrdersConsumer() error {
	stockPubSubConfig := config.AppConfig.PubSub.Stock
	initConsumer(stockPubSubConfig.ProjectId, stockPubSubConfig.Consumer.BrokerMockOrdersSubId, ordersCallback)
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
	var order model.StockOrder
	json.Unmarshal(msg.Data, &order)
	service.PublishOrder(&order)
	fmt.Printf("Got message: %q\n\n", string(msg.Data))
	msg.Ack()
}
