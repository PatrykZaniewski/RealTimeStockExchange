package database

import (
	"cloud.google.com/go/firestore"
	"context"
	"fmt"
	"log"
)

func DatabaseOperation() {
	// Use a service account
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "citric-campaign-349210")
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
		fmt.Print("1")
		res := tx.Set(ref, map[string]interface{}{
			"population": pop.(int64) + 1,
		}, firestore.MergeAll)
		fmt.Print("2")
		return res
	})
	fmt.Println("3")
	if err != nil {
		// Handle any errors appropriately in this section.
		log.Printf("An error has occurred: %s", err)
	}
	//for {
	//	doc, err := iter.Next()
	//	if err == iterator.Done {
	//		break
	//	}
	//	if err != nil {
	//		log.Fatalf("Failed to iterate: %v", err)
	//	}
	//	fmt.Println(doc.Data())
	//}
	defer client.Close()
}
