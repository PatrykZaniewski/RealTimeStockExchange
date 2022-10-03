package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	config "stock/stock_exchange_core/config/env"
	"stock/stock_exchange_core/domain/model"
)

func PublishPrices(newPrice *model.Price) error {
	pubSubConfig := config.AppConfig.PubSub
	projectId := pubSubConfig.Stock.ProjectId
	topicId := pubSubConfig.Stock.Publisher.PricesTopicId

	err := PublishMessage(projectId, topicId, newPrice)
	if err != nil {
		return err
	}
	return nil
}

func PublishMessage(projectId, topicID string, msg interface{}) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		return fmt.Errorf("pubsub: NewClient: %v", err)
	}
	defer client.Close()
	jsonMsg, _ := json.Marshal(msg)

	t := client.Topic(topicID)
	result := t.Publish(ctx, &pubsub.Message{
		Data: jsonMsg,
	})

	id, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("pubsub: result.Get: %v", err)
	}
	fmt.Printf("Published a message; msg ID: %v\n", id)
	return nil
}
