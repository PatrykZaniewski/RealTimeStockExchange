package rest

import (
	"fmt"
	"log"
	"net/http"
	config "stock/order_collector/config/env"
	"sync"
)

func HandleRequests(wg *sync.WaitGroup) {
	defer wg.Done()
	generalConfig := config.AppConfig.General
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":"+generalConfig.Rest.Port, nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}
