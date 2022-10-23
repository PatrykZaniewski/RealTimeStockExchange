package service

import (
	config "broker/broker_core/config/env"
	"broker/broker_core/domain/model"
	"broker/broker_core/interface/pubsub/publisher"
	"cloud.google.com/go/firestore"
	"context"
	"log"
	"strconv"
	"time"
)

func ProcessOrder(internalOrder *model.InternalOrder) {
	if internalOrder.ClientId != "mock_client" {
		log.Printf("%s,BROKER_CORE,ORDER_RECEIVED,%s", internalOrder.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	}

	firestoreStockConfig := config.AppConfig.Firestore.StockConfig
	projectId := firestoreStockConfig.ProjectId
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)
	ordersRef := client.Collection("orders").Doc(internalOrder.ClientId)
	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(ordersRef)
		if err != nil {
			return err
		}
		pendingOrders, _ := doc.DataAt("pendingOrders")
		if err != nil {
			return err
		}

		_ = tx.Set(ordersRef, map[string]interface{}{
			"pendingOrders": append(pendingOrders.([]interface{}), map[string]interface{}{
				"assetName":    internalOrder.AssetName,
				"clientId":     internalOrder.ClientId,
				"id":           internalOrder.Id,
				"orderPrice":   internalOrder.OrderPrice,
				"orderType":    internalOrder.OrderType,
				"orderSubtype": internalOrder.OrderSubtype,
				"quantity":     internalOrder.Quantity,
				"requestTime":  time.Now(),
			}),
		}, firestore.MergeAll)

		return nil
	})
	if err != nil {
		// Handle any errors appropriately in this section.
		log.Printf("An error has occurred: %s", err)
	}
	defer client.Close()

	publisher.PublishOrder(internalOrder)
}
