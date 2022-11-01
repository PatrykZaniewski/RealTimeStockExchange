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
	"sync"
	"time"
)

func PublishStatusOrder(orderStatus *model.OrderStatus) {
	if orderStatus.BrokerId != "mock_broker" && orderStatus.ClientId != "mock_client" {
		log.Printf("%s,BROKER_CORE,STATUS_RECEIVED,%s", orderStatus.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	}

	firestoreStockConfig := config.AppConfig.Firestore.StockConfig
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: firestoreStockConfig.ProjectId,
	}
	app, _ := firebase.NewApp(ctx, conf, option.WithoutAuthentication())
	client, _ := app.Database(ctx)
	pendingOrderRef := client.NewRef("orders/" + orderStatus.ClientId + "/pendingOrders")
	archiveOrderRef := client.NewRef("orders/" + orderStatus.ClientId + "/executedOrders")
	walletRef := client.NewRef("wallet/" + orderStatus.ClientId + "/" + orderStatus.AssetName)

	if orderStatus.BrokerId != "mock_broker" && orderStatus.ClientId != "mock_client" {
		log.Printf("%s,BROKER_CORE,STATUS_PROCESSING,%s", orderStatus.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	}

	var orderToBeExecuted model.PendingOrder

	removePendingOrder := func(t db.TransactionNode) (interface{}, error) {
		var pendingOrders []model.PendingOrder
		t.Unmarshal(&pendingOrders)

		var newPendingOrders []interface{}

		for _, s := range pendingOrders {
			if s.Id != orderStatus.Id {
				newPendingOrders = append(newPendingOrders, s)
			} else {
				orderToBeExecuted = s
			}
		}

		return newPendingOrders, nil
	}

	addArchiveOrder := func(t db.TransactionNode) (interface{}, error) {
		var archiveOrders []interface{}
		t.Unmarshal(&archiveOrders)

		return append(archiveOrders, map[string]interface{}{
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
		}), nil
	}

	updateWallet := func(t db.TransactionNode) (interface{}, error) {
		var assetCountTmp float64
		t.Unmarshal(&assetCountTmp)
		assetCount := int(assetCountTmp)

		if orderStatus.OrderType == "BUY" {
			return assetCount + 1, nil
		} else {
			return assetCount - 1, nil
		}
	}

	var wg sync.WaitGroup
	wg.Add(3)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		pendingOrderRef.Transaction(ctx, removePendingOrder)
	}(&wg)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		archiveOrderRef.Transaction(ctx, addArchiveOrder)
	}(&wg)
	go func(wg *sync.WaitGroup) {
		defer wg.Done()
		walletRef.Transaction(ctx, updateWallet)
	}(&wg)

	wg.Wait()

	if orderStatus.BrokerId != "mock_broker" && orderStatus.ClientId != "mock_client" {
		log.Printf("%s,BROKER_CORE,STATUS_PROCESSED,%s", orderStatus.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	}

	publisher.PublishOrderStatus(orderStatus)
}
