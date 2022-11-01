package service

import (
	config "broker/broker_core/config/env"
	"broker/broker_core/domain/model"
	"broker/broker_core/interface/pubsub/publisher"
	"context"
	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"google.golang.org/api/option"
	"log"
	"strconv"
	"time"
)

func ProcessOrder(internalOrder *model.InternalOrder) {
	if internalOrder.ClientId != "mock_client" {
		log.Printf("%s,BROKER_CORE,ORDER_RECEIVED,%s", internalOrder.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	}

	firestoreStockConfig := config.AppConfig.Firestore.StockConfig
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: firestoreStockConfig.ProjectId,
	}
	app, _ := firebase.NewApp(ctx, conf, option.WithoutAuthentication())
	client, _ := app.Database(ctx)
	pendingOrderRef := client.NewRef("orders/" + internalOrder.ClientId + "/pendingOrders")

	if internalOrder.ClientId != "mock_client" {
		log.Printf("%s,BROKER_CORE,ORDER_PROCESSING,%s", internalOrder.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	}

	orderAdd := func(t db.TransactionNode) (interface{}, error) {
		var currentValue []interface{}
		t.Unmarshal(&currentValue)
		return append(currentValue, map[string]interface{}{
			"assetName":    internalOrder.AssetName,
			"clientId":     internalOrder.ClientId,
			"id":           internalOrder.Id,
			"orderPrice":   internalOrder.OrderPrice,
			"orderType":    internalOrder.OrderType,
			"orderSubtype": internalOrder.OrderSubtype,
			"quantity":     internalOrder.Quantity,
			"requestTime":  time.Now(),
		}), nil
	}

	pendingOrderRef.Transaction(ctx, orderAdd)

	if internalOrder.ClientId != "mock_client" {
		log.Printf("%s,BROKER_CORE,ORDER_PROCESSED,%s", internalOrder.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	}

	publisher.PublishOrder(internalOrder)
}
