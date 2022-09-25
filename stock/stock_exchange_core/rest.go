package main

import (
	"fmt"
	"log"
	"net/http"
)

func HandleRequests() {
	http.HandleFunc("/", homePage)
	log.Fatal(http.ListenAndServe(":5005", nil))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}
