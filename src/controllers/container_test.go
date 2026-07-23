package controllers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStartContainerHandler(t *testing.T) {
	body := `{"challenge_id":"chal-123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/challenges/start", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	StartContainerHandler(w, req)

	if w.Code != http.StatusCreated && w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 201 or 500, got %d: %s", w.Code, w.Body.String())
	}
}

func TestStartContainerHandlerMissingChallengeID(t *testing.T) {
	body := `{}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/challenges/start", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	StartContainerHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestStopContainerHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/instances/abc123/stop", nil)
	w := httptest.NewRecorder()

	StopContainerHandler(w, req)

	// Should get 404 since instance doesn't exist (or 500 if docker not available)
	if w.Code != http.StatusNotFound && w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 404 or 500, got %d", w.Code)
	}
}

func TestContainerStatusHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/instances/abc123/status", nil)
	w := httptest.NewRecorder()

	ContainerStatusHandler(w, req)

	if w.Code != http.StatusNotFound && w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 404 or 500, got %d", w.Code)
	}
}

func TestListMyContainersHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/v1/instances", nil)
	w := httptest.NewRecorder()

	ListMyContainersHandler(w, req)

	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("expected status 200 or 500, got %d", w.Code)
	}
}

func TestContainerHandlerMalformedJSON(t *testing.T) {
	body := `{bad`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/challenges/start", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	StartContainerHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}
