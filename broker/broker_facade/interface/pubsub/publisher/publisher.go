package publisher

import (
	config "broker/broker_facade/config/env"
	"broker/broker_facade/domain/model"
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

func PublishOrder(data *model.InternalOrder) error {
	pubSubConfig := config.AppConfig.PubSub
	projectId := pubSubConfig.Broker.ProjectId
	topicId := pubSubConfig.Broker.Publisher.BrokerInternalClientOrdersTopicId

	if data.ClientId != "mock_client" {
		log.Printf("%s,BROKER_FACADE,ORDER_SENDING,%s", data.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	}
	err := publishMessage(projectId, topicId, data)
	if data.ClientId != "mock_client" {
		log.Printf("%s,BROKER_FACADE,ORDER_SEND,%s", data.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	}
	if err != nil {
		return err
	}
	return nil
}

func publishMessage(projectId, topicID string, msg interface{}) error {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectId)
	defer client.Close()
	if err != nil {
		return fmt.Errorf("pubsub: NewClient: %v", err)
	}

	jsonMsg, _ := json.Marshal(msg)
	t := client.Topic(topicID)
	t.PublishSettings.DelayThreshold = -1 * time.Millisecond
	t.PublishSettings.CountThreshold = 1
	t.PublishSettings.ByteThreshold = 1
	result := t.Publish(ctx, &pubsub.Message{
		Data: jsonMsg,
	})

	_, err = result.Get(ctx)
	if err != nil {
		return fmt.Errorf("pubsub: result.Get: %v", err)
	}
	//fmt.Printf("Published a message; msg ID: %v\n", id)
	return nil
}
