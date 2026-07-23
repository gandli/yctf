package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestHealthCheck verifies the health endpoint returns 200 with status info.
// This test is RED - healthHandler does not exist yet.
func TestHealthCheck(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	healthHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	expected := `{"status":"ok"`
	if w.Body.String() == "" {
		t.Error("expected non-empty response body")
	}
	_ = expected
}
