package service

import (
	"cloud.google.com/go/firestore"
	"context"
	"encoding/json"
	"log"
	"reflect"
	"sort"
	config "stock/stock_exchange_core/config/env"
	"stock/stock_exchange_core/domain/model"
	"stock/stock_exchange_core/interface/pubsub/publisher"
	"strconv"
	"time"
)

func ProcessLimitOrder() *model.OrderStatus {
	//stockOrder *model.StockOrder
	var stockOrder = &model.StockOrder{
		AssetName:    "ASSECO",
		Quantity:     1,
		OrderType:    "SELL",
		OrderSubtype: "MARKET",
		OrderPrice:   210.00,
		ClientId:     "broker_client_1",
		BrokerId:     "common_broker_1",
		Id:           "zxc",
	}

	firestoreStockConfig := config.AppConfig.Firestore.StockConfig
	projectId := firestoreStockConfig.ProjectId
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)
	orderBookRef := client.Collection("orderBook").Doc(stockOrder.AssetName)
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

func ProcessMarketOrder() *model.OrderStatus {
	//stockOrder *model.StockOrder

	var stockOrder = &model.StockOrder{
		AssetName:    "ASSECO",
		Quantity:     1,
		OrderType:    "SELL",
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
	orderBookRef := client.Collection("orderBook").Doc(stockOrder.AssetName)
	newPrices := (*model.Price)(nil)

	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {

		doc, err := tx.Get(orderBookRef)
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
		ordersDb, _ := doc.DataAt(ordersType)
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

		//limitOrderToBeExecuted := orders[0]

		_ = tx.Set(orderBookRef, map[string]interface{}{
			ordersType: orders[1:],
			priceType:  orders[1].Price,
		}, firestore.MergeAll)

		//var currentBuyPriceTmp interface{}
		//var currentSellPriceTmp interface{}
		//
		//var currentBuyPrice float64
		//var currentSellPrice float64
		//
		//currentBuyPriceTmp, _ = doc.DataAt("buyPrice")
		//currentSellPriceTmp, _ = doc.DataAt("sellPrice")
		//if reflect.TypeOf(currentBuyPriceTmp).Kind() == reflect.Int64 {
		//	currentBuyPrice = float64(currentBuyPriceTmp.(int64))
		//}
		//if reflect.TypeOf(currentSellPriceTmp).Kind() == reflect.Int64 {
		//	currentSellPrice = float64(currentSellPriceTmp.(int64))
		//}
		//
		//if stockOrder.OrderType == "BUY" && stockOrder.OrderPrice > currentSellPrice {
		//	_ = tx.Set(orderBookRef, map[string]interface{}{
		//		"sellPrice": stockOrder.OrderPrice,
		//	}, firestore.MergeAll)
		//	newPrices = &model.Price{
		//		AssetName: stockOrder.AssetName,
		//		SellPrice: stockOrder.OrderPrice,
		//		BuyPrice:  currentBuyPrice,
		//	}
		//}
		//
		//if stockOrder.OrderType == "SELL" && stockOrder.OrderPrice < currentBuyPrice {
		//	_ = tx.Set(orderBookRef, map[string]interface{}{
		//		"buyPrice": stockOrder.OrderPrice,
		//	}, firestore.MergeAll)
		//	newPrices = &model.Price{
		//		AssetName: stockOrder.AssetName,
		//		SellPrice: currentSellPrice,
		//		BuyPrice:  stockOrder.OrderPrice,
		//	}
		//}
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
