package service

import (
	config "broker/broker_core/config/env"
	"broker/broker_core/domain/model"
	"broker/broker_core/interface/pubsub/publisher"
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"
)

func PublishStatusOrder(orderStatus *model.OrderStatus) {
	if orderStatus.BrokerId != "mock_broker" && orderStatus.ClientId != "mock_client" {
		log.Printf("%s,BROKER_CORE,STATUS_RECEIVED,%s", orderStatus.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	}

	firestoreStockConfig := config.AppConfig.Firestore.StockConfig
	projectId := firestoreStockConfig.ProjectId
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)
	ordersRef := client.Collection("orders").Doc(orderStatus.ClientId)
	walletRef := client.Collection("wallet").Doc(orderStatus.ClientId)

	if orderStatus.BrokerId != "mock_broker" && orderStatus.ClientId != "mock_client" {
		log.Printf("%s,BROKER_CORE,STATUS_PROCESSING,%s", orderStatus.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	}

	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		ordersDoc, _ := tx.Get(ordersRef)
		pendingOrdersDb, _ := ordersDoc.DataAt("pendingOrders")

		var pendingOrders []*model.PendingOrder
		pendingOrdersRaw, _ := json.Marshal(pendingOrdersDb)
		json.Unmarshal(pendingOrdersRaw, &pendingOrders)

		var newPendingOrders []interface{}
		var orderToBeExecuted model.PendingOrder

		for _, s := range pendingOrders {
			if s.Id != orderStatus.Id {
				newPendingOrders = append(newPendingOrders, s)
			} else {
				orderToBeExecuted = *s
			}
		}

		executedOrders, _ := ordersDoc.DataAt("executedOrders")

		_ = tx.Set(ordersRef, map[string]interface{}{
			"executedOrders": append(executedOrders.([]interface{}), map[string]interface{}{
				"assetName":      orderToBeExecuted.AssetName,
				"clientId":       orderToBeExecuted.ClientId,
				"id":             orderToBeExecuted.Id,
				"orderPrice":     orderToBeExecuted.OrderPrice,
				"executionPrice": orderStatus.ExecutionPrice,
				"orderType":      orderToBeExecuted.OrderType,
				"orderSubtype":   orderToBeExecuted.OrderSubtype,
				"quantity":       orderToBeExecuted.Quantity,
				"requestTime":    orderToBeExecuted.RequestTime,
				"executionTime":  time.Now(),
			}),
			"pendingOrders": newPendingOrders,
		}, firestore.MergeAll)
		return nil
	})

	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		walletDoc, _ := tx.Get(walletRef)

		clientAssetCount, _ := walletDoc.DataAt(orderStatus.AssetName)

		if orderStatus.OrderType == "BUY" {
			_ = tx.Set(walletRef, map[string]interface{}{
				orderStatus.AssetName: clientAssetCount.(int64) + 1,
			}, firestore.MergeAll)
		} else {
			_ = tx.Set(walletRef, map[string]interface{}{
				orderStatus.AssetName: clientAssetCount.(int64) - 1,
			}, firestore.MergeAll)
		}

		return nil
	})

	if orderStatus.BrokerId != "mock_broker" && orderStatus.ClientId != "mock_client" {
		log.Printf("%s,BROKER_CORE,STATUS_PROCESSED,%s", orderStatus.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	}

	if err != nil {
		// Handle any errors appropriately in this section.
		log.Printf("An error has occurred: %s", err)
	}
	defer client.Close()
	publisher.PublishOrderStatus(orderStatus)
}
