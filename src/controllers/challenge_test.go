package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateChallengeHandler(t *testing.T) {
	body := `{
		"title": "SQL Injection 101",
		"description": "Learn basic SQL injection",
		"category": "web",
		"points": 100,
		"container_image": "yctf/sqli101:latest"
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/challenges", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateChallengeHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if resp["id"] == nil {
		t.Error("expected challenge ID in response")
	}
}

func TestCreateChallengeHandlerInvalidCategory(t *testing.T) {
	body := `{"title": "Bad Cat", "category": "invalid", "points": 100}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/challenges", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateChallengeHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestCreateChallengeHandlerMissingTitle(t *testing.T) {
	body := `{"category": "web", "points": 100}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/challenges", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateChallengeHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestCreateChallengeHandlerZeroPoints(t *testing.T) {
	body := `{"title": "Free Points", "category": "misc", "points": 0}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/challenges", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateChallengeHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestCreateChallengeHandlerMalformedJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/challenges", bytes.NewBufferString("{bad"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateChallengeHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestListChallengesHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/challenges", nil)
	w := httptest.NewRecorder()

	ListChallengesHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if resp["challenges"] == nil {
		t.Error("expected challenges array in response")
	}
}

func TestListChallengesHandlerWithCategory(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/challenges?category=web", nil)
	w := httptest.NewRecorder()

	ListChallengesHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestGetChallengeHandler(t *testing.T) {
	// First create a challenge
	body := `{"title": "Test Chal", "category": "web", "points": 100}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/challenges", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	CreateChallengeHandler(w, req)

	var created map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &created)
	challengeID := created["id"].(string)

	// Now get it
	req = httptest.NewRequest(http.MethodGet, "/api/v1/challenges/"+challengeID, nil)
	w = httptest.NewRecorder()

	GetChallengeHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}

func TestGetChallengeHandlerNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/challenges/nonexistent-id", nil)
	w := httptest.NewRecorder()

	GetChallengeHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}
