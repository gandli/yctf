package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAdminSubmissionsAuditLog(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/submissions", nil)
	w := httptest.NewRecorder()

	AdminSubmissionsAuditLogHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["submissions"] == nil {
		t.Error("expected submissions field")
	}
}

func TestAdminWriteupReviewQueue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/writeups/pending", nil)
	w := httptest.NewRecorder()

	AdminWriteupReviewQueueHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
}

func TestAdminScoreWriteupHandler(t *testing.T) {
	body := `{"writeup_id":"writeup-123","score":5}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/writeups/score", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	AdminScoreWriteupHandler(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("expected 200 or 404, got %d", w.Code)
	}
}
