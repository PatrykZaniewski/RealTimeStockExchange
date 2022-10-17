package service

import (
	"cloud.google.com/go/firestore"
	"context"
	"log"
	"reflect"
	config "stock/stock_exchange_core/config/env"
	"stock/stock_exchange_core/domain/model"
	"stock/stock_exchange_core/interface/pubsub/publisher"
	"strconv"
	"time"
)

func ProcessMarketOrder() *model.OrderStatus {
	var stockOrder = &model.StockOrder{
		AssetName:    "ASSECO",
		Quantity:     1,
		OrderType:    "BUY",
		OrderSubtype: "MARKET",
		OrderPrice:   160.00,
		ClientId:     "broker_client_1",
		BrokerId:     "common_broker_1",
		Id:           "zxc",
	}

	firestoreStockConfig := config.AppConfig.Firestore.StockConfig
	projectId := firestoreStockConfig.ProjectId
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)
	orderBookRef := client.Collection("orderBook").Doc("ASSECO")
	newPrices := (*model.Price)(nil)

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
				"orderId":     stockOrder.Id,
				"price":       stockOrder.OrderPrice,
				"quantity":    stockOrder.Quantity,
				"requestTime": time.Now(),
			}),
		}, firestore.MergeAll)

		var currentBuyPriceTmp interface{}
		var currentSellPriceTmp interface{}

		var currentBuyPrice float64
		var currentSellPrice float64

		currentBuyPriceTmp, _ = doc.DataAt("buyPrice")
		currentSellPriceTmp, _ = doc.DataAt("sellPrice")
		if reflect.TypeOf(currentBuyPriceTmp).Kind() == reflect.Int64 {
			currentBuyPrice = float64(currentBuyPriceTmp.(int64))
		}
		if reflect.TypeOf(currentSellPriceTmp).Kind() == reflect.Int64 {
			currentSellPrice = float64(currentSellPriceTmp.(int64))
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

	if newPrices != nil {
		publisher.PublishPrices(newPrices)
	}
	return nil
}

func processLimitOrder(stockOrder *model.StockOrder) *model.OrderStatus {
	firestoreStockConfig := config.AppConfig.Firestore.StockConfig
	projectId := firestoreStockConfig.ProjectId
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)
	orderBookRef := client.Collection("orderBook").Doc("abc")
	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(orderBookRef) // tx.Get, NOT ref.Get!
		if err != nil {
			return err
		}
		pop, err := doc.DataAt("population")
		if err != nil {
			return err
		}
		res := tx.Set(orderBookRef, map[string]interface{}{
			"population": pop.(int64) + 1,
		}, firestore.MergeAll)
		return res
	})
	if err != nil {
		// Handle any errors appropriately in this section.
		log.Printf("An error has occurred: %s", err)
	}
	defer client.Close()
	return nil
}

func ProcessOrder(stockOrder *model.StockOrder) {
	if stockOrder.BrokerId != "mock_broker" && stockOrder.ClientId != "mock_client" {
		log.Printf("%s,STOCK_CORE,ORDER_RECEIVED,%s", stockOrder.Id, strconv.FormatInt(time.Now().UnixMicro(), 10))
	}

	//switch stockOrder.OrderType {
	//case "LIMIT":
	//	status := processLimitOrder(stockOrder)
	//	break
	//case "MARKET":
	//	status := processMarketOrder(stockOrder)
	//	break
	//}

	var orderStatus = model.OrderStatus{
		AssetName:      stockOrder.AssetName,
		Quantity:       stockOrder.Quantity,
		OrderType:      stockOrder.OrderType,
		OrderSubtype:   stockOrder.OrderSubtype,
		OrderPrice:     stockOrder.OrderPrice,
		ExecutionPrice: stockOrder.OrderPrice,
		ClientId:       stockOrder.ClientId,
		BrokerId:       stockOrder.BrokerId,
		Id:             stockOrder.Id,
		Status:         model.FULFILLED,
	}

	publisher.PublishOrderStatus(&orderStatus)

	//var price = model.Price{
	//	AssetName: stockOrder.AssetName,
	//	BuyPrice:  math.Round((rand.Float64()*100+600)*100) / 100,
	//	SellPrice: math.Round((rand.Float64()*100+600)*100) / 100,
	//}
	//publisher.PublishPrices(&price)
}
