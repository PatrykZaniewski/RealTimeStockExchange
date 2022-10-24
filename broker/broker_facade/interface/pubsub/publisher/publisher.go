package publisher

import (
	config "broker/broker_facade/config/env"
	"broker/broker_facade/domain/model"
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"
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

	_, err = result.Get(ctx)
	if err != nil {
		return fmt.Errorf("pubsub: result.Get: %v", err)
	}
	//fmt.Printf("Published a message; msg ID: %v\n", id)
	ref := reflect.ValueOf(msg)
	orderId := reflect.Indirect(ref).FieldByName("Id")
	log.Printf("%s,BROKER_FACADE,ORDER_SEND,%s", orderId, strconv.FormatInt(time.Now().UnixMicro(), 10))
	return nil
}
