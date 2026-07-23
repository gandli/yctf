package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimitMiddleware(t *testing.T) {
	handler := RateLimitMiddleware(5, 1*time.Second)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// First 5 calls should succeed
	for i := 0; i < 5; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "10.0.0.1:1234"
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("call %d: expected 200, got %d", i+1, w.Code)
		}
	}

	// 6th call should be rate limited
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "10.0.0.1:1234"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("expected 429, got %d", w.Code)
	}
}

func TestRateLimitMiddlewareDifferentIPs(t *testing.T) {
	handler := RateLimitMiddleware(2, 1*time.Second)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// IP 1 uses 2 calls
	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		req.RemoteAddr = "10.0.0.2:1234"
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("IP1 call %d: expected 200, got %d", i+1, w.Code)
		}
	}

	// IP 1 is now rate limited
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "10.0.0.2:1234"
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("IP1 should be rate limited, got %d", w.Code)
	}

	// IP 2 should still work
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = "10.0.0.3:1234"
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("IP2 should not be rate limited, got %d", w.Code)
	}
}

func TestRateLimitMiddlewareWindowReset(t *testing.T) {
	handler := RateLimitMiddleware(1, 100*time.Millisecond)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	ip := "10.0.0.4:1234"

	// First call succeeds
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = ip
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("first call: expected 200, got %d", w.Code)
	}

	// Second call fails (rate limited)
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = ip
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusTooManyRequests {
		t.Errorf("second call: expected 429, got %d", w.Code)
	}

	// Wait for window to reset
	time.Sleep(150 * time.Millisecond)

	// Third call should succeed again
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	req.RemoteAddr = ip
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("after reset: expected 200, got %d", w.Code)
	}
}
