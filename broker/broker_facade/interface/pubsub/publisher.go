package pubsub

import (
	config "broker/broker_facade/config/env"
	"broker/broker_facade/domain/model"
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
)

func PublishOrder(data *model.InternalOrder) error {
	pubSubConfig := config.AppConfig.PubSub
	projectId := pubSubConfig.Broker.ProjectId
	topicId := pubSubConfig.Broker.Publisher.BrokerInternalClientOrdersTopicId

	err := publishMessage(projectId, topicId, data)
	if err != nil {
		return err
	}
	return nil
}

func publishMessage(projectId, topicID string, msg interface{}) error {
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
