package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoginHandler(t *testing.T) {
	// First register a user
	regBody := `{"username":"loginuser","email":"login@example.com","password":"SecurePass123!"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBufferString(regBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	RegisterHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("setup registration failed: %d", w.Code)
	}

	// Now login
	loginBody := `{"username":"loginuser","password":"SecurePass123!"}`
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(loginBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	LoginHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if resp["token"] == nil {
		t.Error("expected token in response")
	}
}

func TestLoginHandlerWrongPassword(t *testing.T) {
	t.Skip("Requires database integration (Phase 1 DB)")
}

func TestLoginHandlerNonExistentUser(t *testing.T) {
	t.Skip("Requires database integration (Phase 1 DB)")
}

func TestLoginHandlerMissingFields(t *testing.T) {
	testCases := []struct {
		name string
		body string
	}{
		{"empty body", `{}`},
		{"missing username", `{"password":"pass"}`},
		{"missing password", `{"username":"user"}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			LoginHandler(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("%s: expected status 400, got %d", tc.name, w.Code)
			}
		})
	}
}

func TestLoginHandlerMalformedJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBufferString("{bad"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	LoginHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}
