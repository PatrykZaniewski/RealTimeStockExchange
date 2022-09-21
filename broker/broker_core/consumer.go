package main

import (
	"cloud.google.com/go/pubsub"
	"context"
	"fmt"
	"sync/atomic"
	"time"
)

func test(_ context.Context, msg *pubsub.Message) {
	var received int32

	fmt.Println("Got message: %q\n", string(msg.Data))
	time.Sleep(15 * time.Second)
	fmt.Println("XD")
	atomic.AddInt32(&received, 1)
	msg.Ack()
}

func PullMsg(c chan bool) error {
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

	err = sub.Receive(ctx, test)
	if err != nil {
		return fmt.Errorf("sub.Receive: %v", err)
	}

	c <- true
	return nil
}
