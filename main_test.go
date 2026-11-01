package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutes(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		path        string
		wantStatus  int
		wantBody    string // checked only when checkBody is true
		checkBody   bool
		wantHeaders map[string]string
	}{
		{
			name:        "GET /health",
			method:      http.MethodGet,
			path:        "/health",
			wantStatus:  http.StatusOK,
			wantBody:    `{"status":"ok"}`,
			checkBody:   true,
			wantHeaders: map[string]string{"Content-Type": "application/json"},
		},
		{
			name:        "GET /",
			method:      http.MethodGet,
			path:        "/",
			wantStatus:  http.StatusOK,
			wantBody:    "hello world",
			checkBody:   true,
			wantHeaders: map[string]string{"Content-Type": "text/plain; charset=utf-8"},
		},
		{
			name:        "POST /health",
			method:      http.MethodPost,
			path:        "/health",
			wantStatus:  http.StatusMethodNotAllowed,
			wantHeaders: map[string]string{"Allow": "GET"},
		},
		{
			name:        "PUT /health",
			method:      http.MethodPut,
			path:        "/health",
			wantStatus:  http.StatusMethodNotAllowed,
			wantHeaders: map[string]string{"Allow": "GET"},
		},
		{
			name:        "DELETE /health",
			method:      http.MethodDelete,
			path:        "/health",
			wantStatus:  http.StatusMethodNotAllowed,
			wantHeaders: map[string]string{"Allow": "GET"},
		},
		{
			name:        "HEAD /health",
			method:      http.MethodHead,
			path:        "/health",
			wantStatus:  http.StatusMethodNotAllowed,
			wantHeaders: map[string]string{"Allow": "GET"},
		},
		{
			name:        "POST /",
			method:      http.MethodPost,
			path:        "/",
			wantStatus:  http.StatusMethodNotAllowed,
			wantHeaders: map[string]string{"Allow": "GET"},
		},
		{
			name:       "GET /unknown",
			method:     http.MethodGet,
			path:       "/unknown",
			wantStatus: http.StatusNotFound,
		},
	}

	mux := newMux()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatus)
			}
			if tc.checkBody {
				if got := rec.Body.String(); got != tc.wantBody {
					t.Errorf("body = %q, want %q", got, tc.wantBody)
				}
			}
			for k, want := range tc.wantHeaders {
				if got := rec.Header().Get(k); got != want {
					t.Errorf("header %s = %q, want %q", k, got, want)
				}
			}
		})
	}
}

func TestRequireGet(t *testing.T) {
	t.Run("GET passes through", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		if !requireGet(rec, req) {
			t.Fatal("requireGet(GET) = false, want true")
		}
		if rec.Body.Len() != 0 {
			t.Errorf("body = %q, want empty", rec.Body.String())
		}
		if got := rec.Header().Get("Allow"); got != "" {
			t.Errorf("Allow header = %q, want unset", got)
		}
	})

	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodHead} {
		t.Run(method+" is rejected", func(t *testing.T) {
			req := httptest.NewRequest(method, "/", nil)
			rec := httptest.NewRecorder()

			if requireGet(rec, req) {
				t.Fatalf("requireGet(%s) = true, want false", method)
			}
			if rec.Code != http.StatusMethodNotAllowed {
				t.Errorf("status = %d, want %d", rec.Code, http.StatusMethodNotAllowed)
			}
			if got := rec.Header().Get("Allow"); got != "GET" {
				t.Errorf("Allow header = %q, want %q", got, "GET")
			}
		})
	}
}
