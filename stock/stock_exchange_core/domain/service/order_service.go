package service

import (
	"context"
	"firebase.google.com/go/v4"
	"firebase.google.com/go/v4/db"
	"google.golang.org/api/option"
	"log"
	"sort"
	config "stock/stock_exchange_core/config/env"
	"stock/stock_exchange_core/domain/model"
	"stock/stock_exchange_core/interface/pubsub/publisher"
	"strconv"
	"time"
)

func ProcessLimitOrder(stockOrder *model.StockOrder) (*model.OrderStatus, *model.Price) {
	firestoreStockConfig := config.AppConfig.Firestore.StockConfig
	ctx := context.Background()
	newPrices := (*model.Price)(nil)
	orderStatus := (*model.OrderStatus)(nil)

	conf := &firebase.Config{
		DatabaseURL: firestoreStockConfig.ProjectId,
	}
	app, _ := firebase.NewApp(ctx, conf, option.WithoutAuthentication())
	client, _ := app.Database(ctx)

	buyOrdersRef := client.NewRef("orderBook/" + stockOrder.AssetName + "/buyOrders")
	sellOrdersRef := client.NewRef("orderBook/" + stockOrder.AssetName + "/sellOrders")
	buyPriceRef := client.NewRef("orderBook/" + stockOrder.AssetName + "/buyPrice")
	sellPriceRef := client.NewRef("orderBook/" + stockOrder.AssetName + "/sellPrice")

	orderBookAdd := func(t db.TransactionNode) (interface{}, error) {
		var currentValue []interface{}
		t.Unmarshal(&currentValue)
		return append(currentValue, map[string]interface{}{
			"brokerId":    stockOrder.BrokerId,
			"clientId":    stockOrder.ClientId,
			"id":          stockOrder.Id,
			"price":       stockOrder.OrderPrice,
			"quantity":    stockOrder.Quantity,
			"requestTime": time.Now(),
		}), nil
	}

	var priceChanged = false

	priceUpdate := func(t db.TransactionNode) (interface{}, error) {
		var currentValue float64
		t.Unmarshal(&currentValue)

		if stockOrder.OrderType == "BUY" && stockOrder.OrderPrice > currentValue {
			priceChanged = true
			return stockOrder.OrderPrice, nil
		}

		if stockOrder.OrderType == "SELL" && stockOrder.OrderPrice < currentValue {
			priceChanged = true
			return stockOrder.OrderPrice, nil
		}

		return currentValue, nil
	}

	if stockOrder.OrderType == "BUY" {
		buyOrdersRef.Transaction(ctx, orderBookAdd)
		sellPriceRef.Transaction(ctx, priceUpdate)
	} else {
		sellOrdersRef.Transaction(ctx, orderBookAdd)
		buyPriceRef.Transaction(ctx, priceUpdate)
	}

	if priceChanged {
		var sellPrice float64
		var buyPrice float64
		sellPriceRef.Get(ctx, &sellPrice)
		buyPriceRef.Get(ctx, &buyPrice)
		newPrices = &model.Price{
			AssetName: stockOrder.AssetName,
			SellPrice: sellPrice,
			BuyPrice:  buyPrice,
		}
	}

	orderStatus = &model.OrderStatus{
		AssetName:      stockOrder.AssetName,
		Quantity:       stockOrder.Quantity,
		OrderType:      stockOrder.OrderType,
		OrderSubtype:   stockOrder.OrderSubtype,
		OrderPrice:     stockOrder.OrderPrice,
		ExecutionPrice: stockOrder.OrderPrice,
		ClientId:       stockOrder.ClientId,
		BrokerId:       stockOrder.BrokerId,
		Id:             stockOrder.Id,
		Status:         model.CREATED,
	}

	return orderStatus, newPrices
}

func ProcessMarketOrder(stockOrder *model.StockOrder) (*model.OrderStatus, *model.OrderStatus, *model.Price) {
	firestoreStockConfig := config.AppConfig.Firestore.StockConfig
	ctx := context.Background()

	conf := &firebase.Config{
		DatabaseURL: firestoreStockConfig.ProjectId,
	}
	app, _ := firebase.NewApp(ctx, conf, option.WithoutAuthentication())
	client, _ := app.Database(ctx)

	buyOrdersRef := client.NewRef("orderBook/" + stockOrder.AssetName + "/buyOrders")
	sellOrdersRef := client.NewRef("orderBook/" + stockOrder.AssetName + "/sellOrders")
	buyPriceRef := client.NewRef("orderBook/" + stockOrder.AssetName + "/buyPrice")
	sellPriceRef := client.NewRef("orderBook/" + stockOrder.AssetName + "/sellPrice")
	archiveRef := client.NewRef("archiveOrders/" + stockOrder.AssetName)

	newPrices := (*model.Price)(nil)
	triggeringOrderStatus := (*model.OrderStatus)(nil)
	triggeredOrderStatus := (*model.OrderStatus)(nil)

	var limitOrderToBeExecuted model.OrderBookOrder

	if stockOrder.BrokerId != "mock_broker" && stockOrder.ClientId != "mock_client" {
		log.Printf("%s,STOCK_CORE,ORDER_PROCESSING,%s", stockOrder.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	}

	var newPrice float64

	orderExecute := func(t db.TransactionNode) (interface{}, error) {
		var orders []model.OrderBookOrder
		t.Unmarshal(&orders)
		sort.Slice(orders, func(i, j int) bool {
			if stockOrder.OrderType == "BUY" {
				if orders[i].Price == orders[j].Price {
					return orders[i].RequestTime.Before(orders[j].RequestTime)
				} else {
					return orders[i].Price < orders[j].Price
				}
			} else {
				if orders[i].Price == orders[j].Price {
					return orders[i].RequestTime.Before(orders[j].RequestTime)
				} else {
					return orders[i].Price > orders[j].Price
				}
			}

		})
		limitOrderToBeExecuted = orders[0]
		newPrice = orders[1].Price
		return orders[1:], nil
	}

	updatePrice := func(t db.TransactionNode) (interface{}, error) {
		return newPrice, nil
	}

	limitOrderToBeExecutedMap := map[string]interface{}{
		"brokerId":      limitOrderToBeExecuted.BrokerId,
		"clientId":      limitOrderToBeExecuted.ClientId,
		"id":            limitOrderToBeExecuted.Id,
		"price":         limitOrderToBeExecuted.Price,
		"quantity":      limitOrderToBeExecuted.Quantity,
		"requestTime":   limitOrderToBeExecuted.RequestTime,
		"executionTime": time.Now(),
	}

	upcomingOrderToBeExecutedMap := map[string]interface{}{
		"brokerId":      stockOrder.BrokerId,
		"clientId":      stockOrder.ClientId,
		"id":            stockOrder.Id,
		"price":         stockOrder.OrderPrice,
		"quantity":      stockOrder.Quantity,
		"requestTime":   time.Now(),
		"executionTime": time.Now(),
	}

	updateArchive := func(t db.TransactionNode) (interface{}, error) {
		var currentValue []interface{}
		t.Unmarshal(&currentValue)

		return append(currentValue, []interface{}{limitOrderToBeExecutedMap, upcomingOrderToBeExecutedMap}...), nil
	}

	if stockOrder.OrderType == "BUY" {
		sellOrdersRef.Transaction(ctx, orderExecute)
		buyPriceRef.Transaction(ctx, updatePrice)
	} else {
		buyOrdersRef.Transaction(ctx, orderExecute)
		sellPriceRef.Transaction(ctx, updatePrice)
	}
	archiveRef.Transaction(ctx, updateArchive)

	var newBuyPrice float64
	var newSellPrice float64
	buyPriceRef.Get(ctx, &newBuyPrice)
	sellPriceRef.Get(ctx, &newSellPrice)

	newPrices = &model.Price{
		AssetName: stockOrder.AssetName,
		SellPrice: newSellPrice,
		BuyPrice:  newBuyPrice,
	}

	triggeringOrderStatus = &model.OrderStatus{
		AssetName:      stockOrder.AssetName,
		Quantity:       stockOrder.Quantity,
		OrderType:      stockOrder.OrderType,
		OrderSubtype:   stockOrder.OrderSubtype,
		OrderPrice:     stockOrder.OrderPrice,
		ExecutionPrice: limitOrderToBeExecuted.Price,
		ClientId:       stockOrder.ClientId,
		BrokerId:       stockOrder.BrokerId,
		Id:             stockOrder.Id,
		Status:         model.FULFILLED,
	}

	var triggeredOrderType string

	if stockOrder.OrderType == "BUY" {
		triggeredOrderType = "SELL"
	} else {
		triggeredOrderType = "BUY"
	}

	triggeredOrderStatus = &model.OrderStatus{
		AssetName:      stockOrder.AssetName,
		Quantity:       limitOrderToBeExecuted.Quantity,
		OrderType:      triggeredOrderType,
		OrderSubtype:   "MARKET_ORDER",
		OrderPrice:     limitOrderToBeExecuted.Price,
		ExecutionPrice: limitOrderToBeExecuted.Price,
		ClientId:       limitOrderToBeExecuted.ClientId,
		BrokerId:       limitOrderToBeExecuted.BrokerId,
		Id:             limitOrderToBeExecuted.Id,
		Status:         model.FULFILLED,
	}

	//firstBrokerRef := client.Collection("brokers").Doc(triggeringOrderStatus.BrokerId)
	//secondBrokerRef := client.Collection("brokers").Doc(triggeredOrderStatus.BrokerId)
	//
	//err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
	//	docs, _ := tx.GetAll([]*firestore.DocumentRef{firstBrokerRef, secondBrokerRef})
	//	firstBrokerDoc := docs[0]
	//	secondBrokerDoc := docs[1]
	//
	//	if stockOrder.BrokerId != limitOrderToBeExecuted.BrokerId {
	//
	//		firstBrokerAssetCount, _ := firstBrokerDoc.DataAt(stockOrder.AssetName)
	//		secondBrokerAssetCount, _ := secondBrokerDoc.DataAt(stockOrder.AssetName)
	//
	//		if stockOrder.OrderType == "BUY" {
	//			_ = tx.Set(firstBrokerRef, map[string]interface{}{
	//				stockOrder.AssetName: firstBrokerAssetCount.(int64) + 1,
	//			}, firestore.MergeAll)
	//
	//			_ = tx.Set(firstBrokerRef, map[string]interface{}{
	//				stockOrder.AssetName: secondBrokerAssetCount.(int64) - 1,
	//			}, firestore.MergeAll)
	//		} else {
	//			_ = tx.Set(firstBrokerRef, map[string]interface{}{
	//				stockOrder.AssetName: firstBrokerAssetCount.(int64) - 1,
	//			}, firestore.MergeAll)
	//
	//			_ = tx.Set(firstBrokerRef, map[string]interface{}{
	//				stockOrder.AssetName: secondBrokerAssetCount.(int64) + 1,
	//			}, firestore.MergeAll)
	//		}
	//	}
	//
	//	return nil
	//})

	if stockOrder.BrokerId != "mock_broker" && stockOrder.ClientId != "mock_client" {
		log.Printf("%s,STOCK_CORE,ORDER_PROCESSED,%s", stockOrder.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	}

	if newPrices != nil {
		publisher.PublishPrices(newPrices)
	}
	return triggeringOrderStatus, triggeredOrderStatus, newPrices
}

func ProcessOrder(stockOrder *model.StockOrder) {
	if stockOrder.BrokerId != "mock_broker" && stockOrder.ClientId != "mock_client" {
		log.Printf("%s,STOCK_CORE,ORDER_RECEIVED,%s", stockOrder.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	}
	switch stockOrder.OrderSubtype {
	case "LIMIT_ORDER":
		_, _ = ProcessLimitOrder(stockOrder)
		//publisher.PublishOrderStatus(orderStatus)
		break
	case "MARKET_ORDER":
		triggeringOrderStatus, _, _ := ProcessMarketOrder(stockOrder)
		publisher.PublishOrderStatus(triggeringOrderStatus)
		//publisher.PublishOrderStatus(triggeredOrderStatus)
		break
	}

}
