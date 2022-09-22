package main

import "log"

func main() {
	EnvInit()
	err := PublishMessage("XD")
	if err != nil {
		return
	}
	err = InitPubSubConsumer()
	if err != nil {
		log.Fatalf("PubSub init error occured: %s", err)
	}
}
