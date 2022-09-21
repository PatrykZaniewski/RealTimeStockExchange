package main

import (
	"fmt"
	"net/http"
)

type userHandler struct{}

func (h *userHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// all users request are going to be routed here
}

func main() {
	c := make(chan bool)
	go PullMsg(c)
	fmt.Println("XDABC")

	mux := http.NewServeMux()
	mux.Handle("/users/", &userHandler{})
	http.ListenAndServe(":8080", mux)
}
