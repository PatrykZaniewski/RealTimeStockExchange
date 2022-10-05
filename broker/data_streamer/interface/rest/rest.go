package rest

import (
	config "broker/data_streamer/config/env"
	wb "broker/data_streamer/interface/websocket"
	"fmt"
	"log"
	"net/http"
	"sync"
)

func HandleRequests(wg *sync.WaitGroup) {
	defer wg.Done()
	generalConfig := config.AppConfig.General
	http.HandleFunc("/ws", wb.Websocket)
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":"+generalConfig.Rest.Port, nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}
