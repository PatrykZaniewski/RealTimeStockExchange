package publisher

import (
	config "broker/order_executor/config/env"
	"broker/order_executor/domain/model"
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

func PublishOrder(order *model.StockOrder) error {
	pubSubConfig := config.AppConfig.PubSub
	projectId := pubSubConfig.Broker.ProjectId
	topicId := pubSubConfig.Broker.Publisher.BrokerPendingOrdersTopicId

	if order.ClientId != "mock_client" {
		log.Printf("%s,BROKER_ORDER_EXECUTOR,ORDER_SENDING,%s", order.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	}
	err := publishMessage(projectId, topicId, order)
	if order.BrokerId != "mock_broker" && order.ClientId != "mock_client" {
		log.Printf("%s,BROKER_ORDER_EXECUTOR,ORDER_SEND,%s", order.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	}
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
