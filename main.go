package main

import (
	"log"
	"net/http"
)

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/", helloHandler)
	return mux
}

// requireGet reports whether r is a GET request. If it is not, it writes a
// 405 Method Not Allowed response with an "Allow: GET" header and returns false.
func requireGet(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if !requireGet(w, r) {
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("hello world"))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	if !requireGet(w, r) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", newMux()))
}
