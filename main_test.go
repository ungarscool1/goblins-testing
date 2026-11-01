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
			name:       "GET /",
			method:     http.MethodGet,
			path:       "/",
			wantStatus: http.StatusOK,
			wantBody:   "hello world",
			checkBody:  true,
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
