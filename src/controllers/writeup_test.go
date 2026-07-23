package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSubmitWriteupHandler(t *testing.T) {
	body := `{
		"challenge_id": "chal-123",
		"url": "https://example.com/writeup",
		"content": "Here is my solve..."
	}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/writeups", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	SubmitWriteupHandler(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if resp["id"] == nil {
		t.Error("expected writeup ID in response")
	}
}

func TestSubmitWriteupHandlerMissingChallengeID(t *testing.T) {
	body := `{"url": "https://example.com/writeup"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/writeups", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	SubmitWriteupHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestSubmitWriteupHandlerMissingURLAndContent(t *testing.T) {
	body := `{"challenge_id": "chal-123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/writeups", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	SubmitWriteupHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestListWriteupsHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/writeups", nil)
	w := httptest.NewRecorder()

	ListWriteupsHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if resp["writeups"] == nil {
		t.Error("expected writeups array in response")
	}
}

func TestApproveWriteupHandler(t *testing.T) {
	// First create a writeup
	body := `{"challenge_id": "chal-123", "url": "https://example.com/writeup"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/writeups", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	SubmitWriteupHandler(w, req)

	var created map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &created)
	writeupID := created["id"].(string)

	// Approve it
	req = httptest.NewRequest(http.MethodPost, "/api/v1/writeups/"+writeupID+"/approve", nil)
	w = httptest.NewRecorder()

	ApproveWriteupHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}
}
