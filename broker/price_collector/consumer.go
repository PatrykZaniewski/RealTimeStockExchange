package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"sync/atomic"
)

func ordersCallback(_ context.Context, msg *pubsub.Message) {
	var received int32

	fmt.Printf("Got message: %q\n\n", string(msg.Data))
	atomic.AddInt32(&received, 1)
	msg.Ack()
}

func InitPubSubConsumer() error {
	projectID, ok := viper.Get("cloud.projectid").(string)

	if !ok {
		log.Fatalf("No value for cloud.projectid")
	}
	ordersSubscriptionId, ok := viper.Get("cloud.orderssubscriptionid").(string)
	if !ok {
		log.Fatalf("No value for cloud.orderssubscriptionid")
	}

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, projectID)
	if err != nil {
		_ = fmt.Errorf("PubSub connection error: %v", err)
	}
	defer client.Close()

	sub := client.Subscription(ordersSubscriptionId)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err = sub.Receive(ctx, ordersCallback)
	if err != nil {
		return fmt.Errorf("sub.Receive: %v", err)
	}

	return nil
}
