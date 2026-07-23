package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gandli/yctf/utils"
)

func TestSubmitFlagHandler(t *testing.T) {
	body := `{"title": "Flag Test", "category": "web", "points": 100}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/challenges", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	CreateChallengeHandler(w, req)

	var chal map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &chal)
	challengeID := chal["id"].(string)

	secret := "dev-flag-secret-change-in-production"
	correctFlag := utils.GenerateFlag("team-placeholder", challengeID, secret)

	subBody := `{"challenge_id":"` + challengeID + `","flag":"` + correctFlag + `"}`
	req = httptest.NewRequest(http.MethodPost, "/api/v1/submit", bytes.NewBufferString(subBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	SubmitFlagHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	if resp["correct"] != true {
		t.Error("expected correct=true")
	}
}

func TestSubmitFlagHandlerIncorrect(t *testing.T) {
	body := `{"title": "Flag Test 2", "category": "web", "points": 100}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/challenges", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	CreateChallengeHandler(w, req)

	var chal map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &chal)
	challengeID := chal["id"].(string)

	subBody := `{"challenge_id":"` + challengeID + `","flag":"flag{wrong}"}`
	req = httptest.NewRequest(http.MethodPost, "/api/v1/submit", bytes.NewBufferString(subBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	SubmitFlagHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["correct"] == true {
		t.Error("expected correct=false for wrong flag")
	}
}

func TestSubmitFlagHandlerMissingFields(t *testing.T) {
	testCases := []struct {
		name string
		body string
	}{
		{"empty body", `{}`},
		{"missing challenge_id", `{"flag":"flag{test}"}`},
		{"missing flag", `{"challenge_id":"abc"}`},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/submit", bytes.NewBufferString(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			SubmitFlagHandler(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("%s: expected status 400, got %d", tc.name, w.Code)
			}
		})
	}
}

func TestSubmitFlagHandlerMalformedJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/submit", bytes.NewBufferString("{bad"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	SubmitFlagHandler(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}
}

func TestSubmitFlagHandlerNonExistentChallenge(t *testing.T) {
	subBody := `{"challenge_id":"nonexistent","flag":"flag{test}"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/submit", bytes.NewBufferString(subBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	SubmitFlagHandler(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}
}

func TestSubmitFlagHandlerDuplicateCorrect(t *testing.T) {
	body := `{"title": "Dup Test", "category": "web", "points": 100}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/challenges", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	CreateChallengeHandler(w, req)

	var chal map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &chal)
	challengeID := chal["id"].(string)

	secret := "dev-flag-secret-change-in-production"
	correctFlag := utils.GenerateFlag("team-placeholder", challengeID, secret)

	subBody := `{"challenge_id":"` + challengeID + `","flag":"` + correctFlag + `"}`
	for i := 0; i < 2; i++ {
		req = httptest.NewRequest(http.MethodPost, "/api/v1/submit", bytes.NewBufferString(subBody))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		SubmitFlagHandler(w, req)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["correct"] == true {
		if score, ok := resp["score_gained"].(float64); ok && score > 0 {
			t.Error("duplicate correct submission should not award points again")
		}
	}
}
