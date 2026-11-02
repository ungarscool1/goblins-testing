package main

import (
	"log"
	"net/http"

	"goblins-testing/handler"
)

// newMux routes /health to handler.Health, /version to handler.Version and all other paths to handler.Hello, which serves / and returns 404 for the rest.
func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/version", handler.Version)
	mux.HandleFunc("/", handler.Hello)
	return mux
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", newMux()))
}
