package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	config "stock/stock_exchange_core/config/model"
)

func PublishOrdersStatus(brokerId, msg string) error {
	abc := config.StockConfig.Rest
	fmt.Println(abc)
	return nil
}

type Price struct {
	Value int
	Name  string
}

func PublishPrices() {

}

func PublishMessage(projectId, topicID string, msg interface{}) error {
	projectId = "citric-campaign-349210"
	topicID = "orders"

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		return fmt.Errorf("pubsub: NewClient: %v", err)
	}
	defer client.Close()

	t := client.Topic(topicID)
	//pri := &Price{
	//	Value: 15,
	//	Name:  "ABC",
	//}
	//mes, err := json.Marshal(pri)
	fmt.Println("XD")
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte("A"),
	})

	id, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("pubsub: result.Get: %v", err)
	}
	fmt.Printf("Published a message; msg ID: %v\n", id)
	return nil
}
