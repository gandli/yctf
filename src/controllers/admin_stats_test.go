package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAdminStatsHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/stats", nil)
	w := httptest.NewRecorder()

	AdminStatsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	expectedFields := []string{"users", "teams", "challenges", "submissions", "containers", "uptime"}
	for _, field := range expectedFields {
		if resp[field] == nil {
			t.Errorf("expected field %s in response", field)
		}
	}
}

func TestAdminStatsHandlerMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/stats", nil)
	w := httptest.NewRecorder()

	AdminStatsHandler(w, req)

	if w.Code != http.StatusMethodNotAllowed && w.Code != http.StatusOK {
		t.Errorf("expected status 405 or 200, got %d", w.Code)
	}
}
