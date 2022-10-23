package publisher

import (
	config "broker/broker_core/config/env"
	"broker/broker_core/domain/model"
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

func PublishOrder(order *model.InternalOrder) error {
	pubSubConfig := config.AppConfig.PubSub
	projectId := pubSubConfig.Broker.ProjectId
	topicId := pubSubConfig.Broker.Publisher.BrokerPendingOrdersTopicId

	err := publishMessage(projectId, topicId, order)
	if order.ClientId != "mock_client" {
		log.Printf("%s,BROKER_CORE,ORDER_SEND,%s", order.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	}

	if err != nil {
		return err
	}
	return nil
}

func PublishOrderStatus(order *model.OrderStatus) error {
	pubSubConfig := config.AppConfig.PubSub
	projectId := pubSubConfig.Broker.ProjectId
	topicId := pubSubConfig.Broker.Publisher.BrokerInternalCoreOrdersStatusTopicId

	err := publishMessage(projectId, topicId, order)
	if order.BrokerId != "mock_broker" && order.ClientId != "mock_client" {
		log.Printf("%s,BROKER_CORE,STATUS_SEND,%s", order.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
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
	result := t.Publish(ctx, &pubsub.Message{
		Data: jsonMsg,
	})

	_, err = result.Get(ctx)
	//if err != nil {
	//	return fmt.Errorf("pubsub: result.Get: %v", err)
	//}
	//fmt.Printf("Published a message; msg ID: %v\n", id)
	return nil
}
