package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCORSMiddleware(t *testing.T) {
	handler := CORSMiddleware("http://localhost:5173")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	origin := w.Header().Get("Access-Control-Allow-Origin")
	if origin != "http://localhost:5173" {
		t.Errorf("expected allowed origin, got %s", origin)
	}
}

func TestCORSMiddlewarePreflight(t *testing.T) {
	handler := CORSMiddleware("http://localhost:5173")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodOptions, "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	req.Header.Set("Access-Control-Request-Method", "POST")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204 for preflight, got %d", w.Code)
	}

	method := w.Header().Get("Access-Control-Allow-Methods")
	if method == "" {
		t.Error("expected allowed methods header")
	}
}

func TestCORSMiddlewareDisallowedOrigin(t *testing.T) {
	handler := CORSMiddleware("http://localhost:5173")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://evil.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	origin := w.Header().Get("Access-Control-Allow-Origin")
	if origin == "http://evil.com" {
		t.Error("disallowed origin should not be reflected")
	}
}

func TestCORSMiddlewareMultipleOrigins(t *testing.T) {
	handler := CORSMiddleware("http://localhost:5173", "http://localhost:3000")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Test first allowed origin
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "http://localhost:5173" {
		t.Error("should allow first origin")
	}

	// Test second allowed origin
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w = httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Header().Get("Access-Control-Allow-Origin") != "http://localhost:3000" {
		t.Error("should allow second origin")
	}
}

func TestCORSMiddlewareHeaders(t *testing.T) {
	handler := CORSMiddleware("http://localhost:5173")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	headers := w.Header().Get("Access-Control-Allow-Headers")
	if headers == "" {
		t.Error("expected allow-headers header")
	}

	methods := w.Header().Get("Access-Control-Allow-Methods")
	if methods == "" {
		t.Error("expected allow-methods header")
	}
}
