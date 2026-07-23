package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAdminListUsersHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users", nil)
	w := httptest.NewRecorder()
	AdminListUsersHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["users"] == nil {
		t.Error("expected users field")
	}
}

func TestAdminBanUserHandler(t *testing.T) {
	body := `{"user_id":"user-123","action":"ban"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users/ban", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	AdminBanUserHandler(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("expected 200 or 404, got %d", w.Code)
	}
}

func TestAdminPromoteUserHandler(t *testing.T) {
	body := `{"user_id":"user-456","role":"author"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/admin/users/promote", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	AdminPromoteUserHandler(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusNotFound {
		t.Errorf("expected 200 or 404, got %d", w.Code)
	}
}
