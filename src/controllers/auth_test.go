package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisterHandler(t *testing.T) {
	reqBody := `{"username":"player1","email":"player@example.com","password":"SecurePass123!"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	RegisterHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON response: %v", err)
	}

	if resp["token"] == nil {
		t.Error("expected token in response")
	}
}

func TestRegisterHandlerDuplicate(t *testing.T) {
	t.Skip("Requires database integration (Phase 1 DB)")
}

func TestRegisterHandlerInvalidEmail(t *testing.T) {
	reqBody := `{"username":"player2","email":"not-an-email","password":"SecurePass123!"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	RegisterHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestRegisterHandlerShortPassword(t *testing.T) {
	reqBody := `{"username":"player3","email":"p3@example.com","password":"ab"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	RegisterHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestRegisterHandlerInvalidUsername(t *testing.T) {
	reqBody := `{"username":"ab","email":"p4@example.com","password":"SecurePass123!"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	RegisterHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestRegisterHandlerMissingFields(t *testing.T) {
	testCases := []struct {
		name    string
		body    string
	}{
		{"empty body", `{}`},
		{"missing username", `{"email":"a@b.com","password":"pass"}`},
		{"missing email", `{"username":"user","password":"pass"}`},
		{"missing password", `{"username":"user","email":"a@b.com"}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			RegisterHandler(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("%s: expected status 400, got %d", tc.name, w.Code)
			}
		})
	}
}

func TestRegisterHandlerMalformedJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", strings.NewReader("{invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	RegisterHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400 for malformed JSON, got %d", w.Code)
	}
}
