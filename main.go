package main

import (
	"log"
	"net/http"

	"goblins-testing/handler"
)

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/", handler.Hello)
	return mux
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", newMux()))
}
