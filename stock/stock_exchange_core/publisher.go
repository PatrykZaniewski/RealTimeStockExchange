package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"log"
)

func PublishMessage(msg string) error {

	projectID, ok := viper.Get("cloud.projectid").(string)

	if !ok {
		log.Fatalf("No value for cloud.projectid")
	}
	ordersTopicId, ok := viper.Get("cloud.orderstopicid").(string)
	if !ok {
		log.Fatalf("No value for cloud.orderstopicid")
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		return fmt.Errorf("pubsub: NewClient: %v", err)
	}
	defer client.Close()

	t := client.Topic(ordersTopicId)
	result := t.Publish(ctx, &pubsub.Message{
		Data: []byte(msg),
	})

	id, err := result.Get(ctx)
	if err != nil {
		return fmt.Errorf("pubsub: result.Get: %v", err)
	}
	fmt.Printf("Published a message; msg ID: %v\n", id)
	return nil
}
