package main

func main() {
	InfoLogger.Printf("XD")
	err := InitPubSubConsumer()
	if err != nil {
		return
	}
}
