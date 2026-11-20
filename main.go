package main

import (
	"log"
	"net/http"
	"time"

	"goblins-testing/handler"
)

// newMux routes /health to handler.Health, /uptime to handler.Uptime (counting
// from when newMux is called) and all other paths to handler.Hello, which
// serves / and returns 404 for the rest.
func newMux() *http.ServeMux {
	start := time.Now()
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/uptime", handler.Uptime(start, time.Now))
	mux.HandleFunc("/", handler.Hello)
	return mux
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", newMux()))
}
