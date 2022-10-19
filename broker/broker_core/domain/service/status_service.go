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

	//orderStatus *model.OrderStatus
	//var orderStatus = &model.OrderStatus{
	//	AssetName:      "ABC",
	//	Quantity:       1,
	//	OrderType:      "a",
	//	OrderSubtype:   "b",
	//	OrderPrice:     150.00,
	//	ExecutionPrice: 153.0,
	//	ClientId:       "a",
	//	BrokerId:       "b",
	//	Id:             "de",
	//	Status:         model.CREATED,
	//}

	if orderStatus.ClientId != "mock_client" {
		log.Printf("%s,BROKER_CORE,STATUS_RECEIVED,%s", orderStatus.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	}

	firestoreStockConfig := config.AppConfig.Firestore.StockConfig
	projectId := firestoreStockConfig.ProjectId
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)
	ordersRef := client.Collection("orders").Doc(orderStatus.ClientId)
	walletRef := client.Collection("wallet").Doc(orderStatus.ClientId)

	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		ordersDoc, err := tx.Get(ordersRef)
		walletDoc, err := tx.Get(walletRef)

		if err != nil {
			return err
		}
		var pendingOrders []*model.PendingOrder
		pendingOrdersDb, _ := ordersDoc.DataAt("pendingOrders")
		pendingOrdersRaw, _ := json.Marshal(pendingOrdersDb)
		json.Unmarshal(pendingOrdersRaw, &pendingOrders)
		if err != nil {
			return err
		}

		var newPendingOrders []interface{}
		var orderToBeExecuted *model.PendingOrder

		for _, s := range pendingOrders {
			if s.Id != orderStatus.Id {
				newPendingOrders = append(newPendingOrders, s)
			} else {
				orderToBeExecuted = s
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
		}, firestore.MergeAll)

		_ = tx.Set(ordersRef, map[string]interface{}{
			"pendingOrders": newPendingOrders,
		}, firestore.MergeAll)

		clientAssetCount, _ := walletDoc.DataAt(orderStatus.AssetName)

		if orderStatus.OrderType == "BUY" {
			_ = tx.Set(walletRef, map[string]interface{}{
				orderStatus.AssetName: clientAssetCount.(int64) + 1,
			})
		} else {
			_ = tx.Set(walletRef, map[string]interface{}{
				orderStatus.AssetName: clientAssetCount.(int64) - 1,
			})
		}

		return nil
	})
	if err != nil {
		// Handle any errors appropriately in this section.
		log.Printf("An error has occurred: %s", err)
	}
	defer client.Close()

	publisher.PublishOrderStatus(orderStatus)
}
