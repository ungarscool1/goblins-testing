// Package handler contains the HTTP handlers for the service.
package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

// Hello responds to GET / with "hello world". Other paths return 404 and
// other methods return 405.
func Hello(w http.ResponseWriter, r *http.Request) {
	// Hello is registered on the catch-all "/" pattern, so it receives every
	// unmatched path. Check the path first: an unknown path is 404 whatever
	// the method, and answering 405 with Allow: GET would wrongly suggest the
	// resource exists.
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if !allowGet(w, r) {
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Write([]byte("hello world"))
}

// Health responds to GET with the JSON body {"status":"ok"}. Other methods
// return 405.
func Health(w http.ResponseWriter, r *http.Request) {
	if !allowGet(w, r) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

// Uptime returns a handler that responds to GET with the JSON body
// {"seconds":N}, where N is the number of whole seconds between start and
// now(), rounded down and never negative. Other methods get a JSON 405.
func Uptime(start time.Time, now func() time.Time) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if methodNotAllowedJSON(w, r) {
			return
		}
		seconds := int64(now().Sub(start) / time.Second)
		if seconds < 0 {
			seconds = 0
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct {
			Seconds int64 `json:"seconds"`
		}{seconds})
	}
}

// methodNotAllowedJSON writes a 405 Method Not Allowed response and returns
// true if r is not a GET request. The response has an "Allow: GET" header,
// Content-Type application/json and the body {"error":"method not allowed"}.
// For GET it writes nothing and returns false.
func methodNotAllowedJSON(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == http.MethodGet {
		return false
	}
	w.Header().Set("Allow", http.MethodGet)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusMethodNotAllowed)
	json.NewEncoder(w).Encode(struct {
		Error string `json:"error"`
	}{"method not allowed"})
	return true
}

// allowGet reports whether r is a GET request. If it is not, it writes a
// 405 Method Not Allowed response with an "Allow: GET" header and returns false.
func allowGet(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return false
	}
	return true
}
