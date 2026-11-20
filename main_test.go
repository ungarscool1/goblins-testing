package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// TestRoutes checks that newMux sends each path to the right handler.
// Handler behaviour itself is tested in package handler.
func TestRoutes(t *testing.T) {
	tests := []struct {
		name            string
		path            string
		wantStatus      int
		wantBody        string // checked only when wantStatus is 200 and wantBodyPrefix is empty
		wantBodyPrefix  string // when set, the body only has to start with this
		wantContentType string // checked only when set
	}{
		{name: "GET /", path: "/", wantStatus: http.StatusOK, wantBody: "hello world"},
		{name: "GET /health", path: "/health", wantStatus: http.StatusOK, wantBody: `{"status":"ok"}`},
		{
			name:            "GET /uptime",
			path:            "/uptime",
			wantStatus:      http.StatusOK,
			wantBodyPrefix:  `{"seconds":`,
			wantContentType: "application/json",
		},
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
				got := rec.Body.String()
				if tc.wantBodyPrefix != "" {
					if !strings.HasPrefix(got, tc.wantBodyPrefix) {
						t.Errorf("body = %q, want prefix %q", got, tc.wantBodyPrefix)
					}
				} else if got != tc.wantBody {
					t.Errorf("body = %q, want %q", got, tc.wantBody)
				}
			}
			if tc.wantContentType != "" {
				if got := rec.Header().Get("Content-Type"); got != tc.wantContentType {
					t.Errorf("Content-Type = %q, want %q", got, tc.wantContentType)
				}
			}
		})
	}
}
