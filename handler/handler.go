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

// Echo responds to GET /echo?msg=<value> with the plain-text body <value>,
// exactly as decoded from the query string. A missing or empty msg returns
// 400 and other methods return 405. The response is always text/plain with
// X-Content-Type-Options: nosniff so the reflected input is never rendered
// as HTML.
func Echo(w http.ResponseWriter, r *http.Request) {
	if !allowGet(w, r) {
		return
	}
	msg := r.URL.Query().Get("msg")
	if msg == "" {
		http.Error(w, "missing msg query parameter", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	// A write error cannot be reported: the status and headers are already
	// sent, so the error is discarded on purpose.
	_, _ = w.Write([]byte(msg))
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
