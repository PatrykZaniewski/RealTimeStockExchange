package database

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"log"
	config "stock/stock_exchange_core/config/env"
)

func ProceedDbOperation() {
	firestoreStockConfig := config.AppConfig.Firestore.StockConfig
	projectId := firestoreStockConfig.ProjectId
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectId)
	ref := client.Collection("order_book").Doc("abc")
	err = client.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(ref) // tx.Get, NOT ref.Get!
		if err != nil {
			return err
		}
		pop, err := doc.DataAt("population")
		tmp, err := ref.Get(ctx)
		fmt.Println(tmp)
		if err != nil {
			return err
		}
		res := tx.Set(ref, map[string]interface{}{
			"population": pop.(int64) + 1,
		}, firestore.MergeAll)
		return res
	})
	if err != nil {
		// Handle any errors appropriately in this section.
		log.Printf("An error has occurred: %s", err)
	}
	defer client.Close()
}
