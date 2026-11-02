package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestRoutes checks that newMux sends each path to the right handler.
// Handler behaviour itself is tested in package handler.
func TestRoutes(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		wantStatus int
		wantBody   string // checked only when wantStatus is 200
	}{
		{name: "GET /", path: "/", wantStatus: http.StatusOK, wantBody: "hello world"},
		{name: "GET /health", path: "/health", wantStatus: http.StatusOK, wantBody: `{"status":"ok"}`},
		{name: "GET /nope", path: "/nope", wantStatus: http.StatusNotFound},
	}

	mux := newMux()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tc.path, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tc.wantStatus)
			}
			if tc.wantStatus == http.StatusOK {
				if got := rec.Body.String(); got != tc.wantBody {
					t.Errorf("body = %q, want %q", got, tc.wantBody)
				}
			}
		})
	}
}
