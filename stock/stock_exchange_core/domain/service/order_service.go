package service

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"reflect"
	"sort"
	config "stock/stock_exchange_core/config/env"
	"stock/stock_exchange_core/domain/model"
	"stock/stock_exchange_core/interface/pubsub/publisher"
	"strconv"
	"time"
)

func ProcessLimitOrder() (*model.OrderStatus, *model.Price) {
	//stockOrder *model.StockOrder
	var stockOrder = &model.StockOrder{
		AssetName:    "ASSECO",
		Quantity:     1,
		OrderType:    "SELL",
		OrderSubtype: "MARKET_ORDER",
		OrderPrice:   205.00,
		ClientId:     "broker_client",
		BrokerId:     "common_broker",
		Id:           uuid.New().String(),
	}

	firestoreStockConfig := config.AppConfig.Firestore.StockConfig
	projectId := firestoreStockConfig.ProjectId
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)
	orderBookRef := client.Collection("orderBook").Doc(stockOrder.AssetName)
	newPrices := (*model.Price)(nil)
	orderStatus := (*model.OrderStatus)(nil)

	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(orderBookRef)
		if err != nil {
			return err
		}
		var ordersType string
		if stockOrder.OrderType == "BUY" {
			ordersType = "buyOrders"
		} else {
			ordersType = "sellOrders"
		}
		orders, _ := doc.DataAt(ordersType)
		if err != nil {
			return err
		}

		_ = tx.Set(orderBookRef, map[string]interface{}{
			ordersType: append(orders.([]interface{}), map[string]interface{}{
				"brokerId":    stockOrder.BrokerId,
				"clientId":    stockOrder.ClientId,
				"id":          stockOrder.Id,
				"price":       stockOrder.OrderPrice,
				"quantity":    stockOrder.Quantity,
				"requestTime": time.Now(),
			}),
		}, firestore.MergeAll)

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

		var currentBuyPriceTmp interface{}
		var currentSellPriceTmp interface{}

		var currentBuyPrice float64
		var currentSellPrice float64

		currentBuyPriceTmp, _ = doc.DataAt("buyPrice")
		currentSellPriceTmp, _ = doc.DataAt("sellPrice")
		if reflect.TypeOf(currentBuyPriceTmp).Kind() == reflect.Int64 {
			currentBuyPrice = float64(currentBuyPriceTmp.(int64))
		} else {
			currentBuyPrice = currentBuyPriceTmp.(float64)
		}
		if reflect.TypeOf(currentSellPriceTmp).Kind() == reflect.Int64 {
			currentSellPrice = float64(currentSellPriceTmp.(int64))
		} else {
			currentSellPrice = currentSellPriceTmp.(float64)
		}

		if stockOrder.OrderType == "BUY" && stockOrder.OrderPrice > currentSellPrice {
			_ = tx.Set(orderBookRef, map[string]interface{}{
				"sellPrice": stockOrder.OrderPrice,
			}, firestore.MergeAll)
			newPrices = &model.Price{
				AssetName: stockOrder.AssetName,
				SellPrice: stockOrder.OrderPrice,
				BuyPrice:  currentBuyPrice,
			}
		}

		if stockOrder.OrderType == "SELL" && stockOrder.OrderPrice < currentBuyPrice {
			_ = tx.Set(orderBookRef, map[string]interface{}{
				"buyPrice": stockOrder.OrderPrice,
			}, firestore.MergeAll)
			newPrices = &model.Price{
				AssetName: stockOrder.AssetName,
				SellPrice: currentSellPrice,
				BuyPrice:  stockOrder.OrderPrice,
			}
		}
		return nil
	})
	if err != nil {
		// Handle any errors appropriately in this section.
		log.Printf("An error has occurred: %s", err)
	}
	defer client.Close()

	return orderStatus, newPrices
}

func ProcessMarketOrder(stockOrder *model.StockOrder) (*model.OrderStatus, *model.OrderStatus, *model.Price) {
	//stockOrder *model.StockOrder

	//var stockOrder = &model.StockOrder{
	//	AssetName:    "ASSECO",
	//	Quantity:     1,
	//	OrderType:    "BUY",
	//	OrderSubtype: "MARKET_ORDER",
	//	OrderPrice:   160.00,
	//	ClientId:     "broker_client",
	//	BrokerId:     "common_broker",
	//	Id:           "zxc",
	//}
	firestoreStockConfig := config.AppConfig.Firestore.StockConfig
	projectId := firestoreStockConfig.ProjectId
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)
	orderBookRef := client.Collection("orderBook").Doc(stockOrder.AssetName)
	archiveRef := client.Collection("archiveOrders").Doc(stockOrder.AssetName)
	firstBrokerRef := client.Collection("brokers").Doc(stockOrder.BrokerId)
	newPrices := (*model.Price)(nil)
	triggeringOrderStatus := (*model.OrderStatus)(nil)
	triggeredOrderStatus := (*model.OrderStatus)(nil)

	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {

		orderBookDoc, err := tx.Get(orderBookRef)
		archiveOrdersDoc, _ := tx.Get(archiveRef)
		firstBrokerDoc, _ := tx.Get(firstBrokerRef)
		if err != nil {
			return err
		}
		var ordersType string
		if stockOrder.OrderType == "BUY" {
			ordersType = "sellOrders"
		} else {
			ordersType = "buyOrders"
		}

		var priceType string
		if stockOrder.OrderType == "BUY" {
			priceType = "buyPrice"
		} else {
			priceType = "sellPrice"
		}

		var orders []*model.OrderBookOrder
		ordersDb, _ := orderBookDoc.DataAt(ordersType)
		ordersRaw, _ := json.Marshal(ordersDb)
		json.Unmarshal(ordersRaw, &orders)

		if err != nil {
			return err
		}

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

		limitOrderToBeExecuted := orders[0]

		secondBrokerRef := client.Collection("brokers").Doc(orders[0].BrokerId)
		secondBrokerDoc, _ := tx.Get(secondBrokerRef)

		_ = tx.Set(orderBookRef, map[string]interface{}{
			ordersType: orders[1:],
			priceType:  orders[1].Price,
		}, firestore.MergeAll)

		if stockOrder.OrderType == "BUY" {
			data, _ := orderBookDoc.DataAt("sellPrice")
			var currentSellPrice float64
			if reflect.TypeOf(data).Kind() == reflect.Int64 {
				currentSellPrice = float64(data.(int64))
			} else if reflect.TypeOf(data).Kind() == reflect.Float64 {
				currentSellPrice = data.(float64)
			}
			newPrices = &model.Price{
				AssetName: stockOrder.AssetName,
				BuyPrice:  orders[1].Price,
				SellPrice: currentSellPrice,
			}
		}

		if stockOrder.OrderType == "SELL" {
			data, _ := orderBookDoc.DataAt("buyPrice")
			var currentBuyPrice float64
			if reflect.TypeOf(data).Kind() == reflect.Int64 {
				currentBuyPrice = float64(data.(int64))
			} else if reflect.TypeOf(data).Kind() == reflect.Float64 {
				currentBuyPrice = data.(float64)
			}
			newPrices = &model.Price{
				AssetName: stockOrder.AssetName,
				BuyPrice:  currentBuyPrice,
				SellPrice: orders[1].Price,
			}
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

		archiveOrders, _ := archiveOrdersDoc.DataAt(ordersType)

		_ = tx.Set(archiveRef, map[string]interface{}{
			ordersType: append(archiveOrders.([]interface{}), []interface{}{limitOrderToBeExecutedMap, upcomingOrderToBeExecutedMap}...),
		}, firestore.MergeAll)

		if stockOrder.BrokerId != limitOrderToBeExecuted.BrokerId {

			firstBrokerAssetCount, _ := firstBrokerDoc.DataAt(stockOrder.AssetName)
			secondBrokerAssetCount, _ := secondBrokerDoc.DataAt(stockOrder.AssetName)

			if stockOrder.OrderType == "BUY" {
				_ = tx.Set(firstBrokerRef, map[string]interface{}{
					stockOrder.AssetName: firstBrokerAssetCount.(int64) + 1,
				}, firestore.MergeAll)

				_ = tx.Set(firstBrokerRef, map[string]interface{}{
					stockOrder.AssetName: secondBrokerAssetCount.(int64) - 1,
				}, firestore.MergeAll)
			} else {
				_ = tx.Set(firstBrokerRef, map[string]interface{}{
					stockOrder.AssetName: firstBrokerAssetCount.(int64) - 1,
				}, firestore.MergeAll)

				_ = tx.Set(firstBrokerRef, map[string]interface{}{
					stockOrder.AssetName: secondBrokerAssetCount.(int64) + 1,
				}, firestore.MergeAll)
			}
		}

		return nil
	})

	if err != nil {
		// Handle any errors appropriately in this section.
		log.Printf("An error has occurred: %s", err)
	}
	defer client.Close()

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
		//_, newPrices := ProcessLimitOrder(stockOrder)
		//if newPrices != nil {
		//	publisher.PublishPrices(newPrices)
		//}
		//publisher.PublishOrderStatus(orderStatus)
		break
	case "MARKET_ORDER":
		triggeringOrderStatus, _, newPrices := ProcessMarketOrder(stockOrder)
		publisher.PublishOrderStatus(triggeringOrderStatus)
		//publisher.PublishOrderStatus(triggeredOrderStatus)
		publisher.PublishPrices(newPrices)
		break
	}

}
