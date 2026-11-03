package main

import (
	"encoding/json"
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
		// wantBody is compared exactly when wantStatus is 200 and
		// wantJSONKeys is empty.
		wantBody string
		// wantJSONKeys, if set, replaces the exact body check: the body must
		// decode as a JSON object containing these keys. Used when the body
		// changes between calls.
		wantJSONKeys []string
	}{
		{name: "GET /", path: "/", wantStatus: http.StatusOK, wantBody: "hello world"},
		{name: "GET /health", path: "/health", wantStatus: http.StatusOK, wantBody: `{"status":"ok"}`},
		{name: "GET /time", path: "/time", wantStatus: http.StatusOK, wantJSONKeys: []string{"utc", "unix_nano"}},
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
			if tc.wantStatus != http.StatusOK {
				return
			}
			if len(tc.wantJSONKeys) > 0 {
				var got map[string]json.RawMessage
				if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
					t.Fatalf("decode body %q: %v", rec.Body.String(), err)
				}
				for _, k := range tc.wantJSONKeys {
					if _, ok := got[k]; !ok {
						t.Errorf("body %q has no %q field", rec.Body.String(), k)
					}
				}
				return
			}
			if got := rec.Body.String(); got != tc.wantBody {
				t.Errorf("body = %q, want %q", got, tc.wantBody)
			}
		})
	}
}
