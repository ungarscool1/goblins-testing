// Package handler contains the HTTP handlers for the service.
package handler

import "net/http"

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
