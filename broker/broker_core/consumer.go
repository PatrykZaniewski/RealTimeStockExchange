package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"sync/atomic"
)

func pubSubCallback(_ context.Context, msg *pubsub.Message) {
	var received int32

	fmt.Printf("Got message: %q\n\n", string(msg.Data))
	atomic.AddInt32(&received, 1)
	msg.Ack()
}

func InitPubSubConsumer() error {
	//projectID := os.Getenv("projectId")
	//subID := os.Getenv("subscriptionId")
	projectID := "citric-campaign-349210"
	subID := "abc"
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		fmt.Errorf("pubsub.NewClient: %v", err)
	}
	defer client.Close()

	sub := client.Subscription(subID)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err = sub.Receive(ctx, pubSubCallback)
	if err != nil {
		return fmt.Errorf("sub.Receive: %v", err)
	}

	return nil
}
