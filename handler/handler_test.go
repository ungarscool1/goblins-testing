package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
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

func TestTime(t *testing.T) {
	t.Run("fixed clock", func(t *testing.T) {
		// A non-UTC zone checks the UTC conversion; 100ns checks that the
		// fraction keeps all 9 digits, trailing zeros included.
		fixed := time.Date(2026, 11, 3, 13, 34, 56, 100, time.FixedZone("UTC+1", 60*60))
		orig := now
		now = func() time.Time { return fixed }
		t.Cleanup(func() { now = orig })

		req := httptest.NewRequest(http.MethodGet, "/time", nil)
		rec := httptest.NewRecorder()
		Time(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		for k, want := range map[string]string{
			"Content-Type":  "application/json",
			"Cache-Control": "no-store",
		} {
			if got := rec.Header().Get(k); got != want {
				t.Errorf("header %s = %q, want %q", k, got, want)
			}
		}

		var body struct {
			UTC      string `json:"utc"`
			UnixNano int64  `json:"unix_nano"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("decode body %q: %v", rec.Body.String(), err)
		}
		if want := "2026-11-03T12:34:56.000000100Z"; body.UTC != want {
			t.Errorf("utc = %q, want %q", body.UTC, want)
		}
		if want := fixed.UnixNano(); body.UnixNano != want {
			t.Errorf("unix_nano = %d, want %d", body.UnixNano, want)
		}

		var raw map[string]json.RawMessage
		if err := json.Unmarshal(rec.Body.Bytes(), &raw); err != nil {
			t.Fatalf("decode raw body %q: %v", rec.Body.String(), err)
		}
		rawNano, ok := raw["unix_nano"]
		if !ok {
			t.Fatalf("body %q has no unix_nano field", rec.Body.String())
		}
		if strings.Contains(string(rawNano), `"`) {
			t.Errorf("unix_nano = %s, want a bare JSON number", rawNano)
		}
	})

	t.Run("methods", func(t *testing.T) {
		runCases(t, Time, []handlerCase{
			{
				name:        "POST",
				method:      http.MethodPost,
				path:        "/time",
				wantStatus:  http.StatusMethodNotAllowed,
				wantHeaders: map[string]string{"Allow": "GET"},
			},
			{
				name:        "PUT",
				method:      http.MethodPut,
				path:        "/time",
				wantStatus:  http.StatusMethodNotAllowed,
				wantHeaders: map[string]string{"Allow": "GET"},
			},
			{
				name:        "DELETE",
				method:      http.MethodDelete,
				path:        "/time",
				wantStatus:  http.StatusMethodNotAllowed,
				wantHeaders: map[string]string{"Allow": "GET"},
			},
			{
				name:        "HEAD",
				method:      http.MethodHead,
				path:        "/time",
				wantStatus:  http.StatusMethodNotAllowed,
				wantHeaders: map[string]string{"Allow": "GET"},
			},
		})
	})

	t.Run("real clock", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/time", nil)
		rec := httptest.NewRecorder()

		// Round(0) strips the monotonic reading so comparisons use wall time,
		// like the parsed response.
		before := time.Now().Round(0)
		Time(rec, req)
		after := time.Now().Round(0)

		if rec.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
		}
		var body struct {
			UTC      string `json:"utc"`
			UnixNano int64  `json:"unix_nano"`
		}
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
			t.Fatalf("decode body %q: %v", rec.Body.String(), err)
		}
		if !strings.HasSuffix(body.UTC, "Z") {
			t.Errorf("utc = %q, want a trailing Z", body.UTC)
		}
		parsed, err := time.Parse(time.RFC3339Nano, body.UTC)
		if err != nil {
			t.Fatalf("parse utc %q: %v", body.UTC, err)
		}
		if fromNano := time.Unix(0, body.UnixNano); !fromNano.Equal(parsed) {
			t.Errorf("unix_nano %d is %v, want the same instant as utc %v", body.UnixNano, fromNano, parsed)
		}
		if parsed.Before(before) || parsed.After(after) {
			t.Errorf("utc %v not between %v and %v", parsed, before, after)
		}
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
