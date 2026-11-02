package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type handlerCase struct {
	name        string
	method      string
	path        string
	wantStatus  int
	wantBody    string // checked only when checkBody is true
	checkBody   bool
	wantHeaders map[string]string
}

func runCases(t *testing.T, h http.HandlerFunc, tests []handlerCase) {
	t.Helper()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, tc.path, nil)
			rec := httptest.NewRecorder()
			h(rec, req)

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

func TestHealth(t *testing.T) {
	runCases(t, Health, []handlerCase{
		{
			name:        "GET",
			method:      http.MethodGet,
			path:        "/health",
			wantStatus:  http.StatusOK,
			wantBody:    `{"status":"ok"}`,
			checkBody:   true,
			wantHeaders: map[string]string{"Content-Type": "application/json"},
		},
		{
			name:        "POST",
			method:      http.MethodPost,
			path:        "/health",
			wantStatus:  http.StatusMethodNotAllowed,
			wantHeaders: map[string]string{"Allow": "GET"},
		},
		{
			name:        "PUT",
			method:      http.MethodPut,
			path:        "/health",
			wantStatus:  http.StatusMethodNotAllowed,
			wantHeaders: map[string]string{"Allow": "GET"},
		},
		{
			name:        "DELETE",
			method:      http.MethodDelete,
			path:        "/health",
			wantStatus:  http.StatusMethodNotAllowed,
			wantHeaders: map[string]string{"Allow": "GET"},
		},
		{
			name:        "HEAD",
			method:      http.MethodHead,
			path:        "/health",
			wantStatus:  http.StatusMethodNotAllowed,
			wantHeaders: map[string]string{"Allow": "GET"},
		},
	})
}

func TestVersion(t *testing.T) {
	runCases(t, Version, []handlerCase{
		{
			name:        "GET",
			method:      http.MethodGet,
			path:        "/version",
			wantStatus:  http.StatusOK,
			wantBody:    `{"version":"1"}`,
			checkBody:   true,
			wantHeaders: map[string]string{"Content-Type": "application/json"},
		},
		{
			name:        "POST",
			method:      http.MethodPost,
			path:        "/version",
			wantStatus:  http.StatusMethodNotAllowed,
			wantHeaders: map[string]string{"Allow": "GET"},
		},
		{
			name:        "PUT",
			method:      http.MethodPut,
			path:        "/version",
			wantStatus:  http.StatusMethodNotAllowed,
			wantHeaders: map[string]string{"Allow": "GET"},
		},
		{
			name:        "DELETE",
			method:      http.MethodDelete,
			path:        "/version",
			wantStatus:  http.StatusMethodNotAllowed,
			wantHeaders: map[string]string{"Allow": "GET"},
		},
		{
			name:        "HEAD",
			method:      http.MethodHead,
			path:        "/version",
			wantStatus:  http.StatusMethodNotAllowed,
			wantHeaders: map[string]string{"Allow": "GET"},
		},
	})
}

func TestHello(t *testing.T) {
	runCases(t, Hello, []handlerCase{
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
			name:        "POST /",
			method:      http.MethodPost,
			path:        "/",
			wantStatus:  http.StatusMethodNotAllowed,
			wantHeaders: map[string]string{"Allow": "GET"},
		},
		{
			name:        "HEAD /",
			method:      http.MethodHead,
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
		{
			// An unknown path must give 404 before the method check runs.
			name:        "POST /unknown",
			method:      http.MethodPost,
			path:        "/unknown",
			wantStatus:  http.StatusNotFound,
			wantHeaders: map[string]string{"Allow": ""},
		},
	})
}

func TestAllowGet(t *testing.T) {
	t.Run("GET passes through", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()

		if !allowGet(rec, req) {
			t.Fatal("allowGet(GET) = false, want true")
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

			if allowGet(rec, req) {
				t.Fatalf("allowGet(%s) = true, want false", method)
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
