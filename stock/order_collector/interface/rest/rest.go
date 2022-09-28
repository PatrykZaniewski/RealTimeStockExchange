package rest

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

func HandleRequests(wg *sync.WaitGroup) {
	defer wg.Done()
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":5005", nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}
