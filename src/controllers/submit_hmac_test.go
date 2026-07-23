package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gandli/yctf/utils"
)

func TestSubmitFlagHMACCorrect(t *testing.T) {
	body := `{"title": "HMAC Test", "category": "web", "points": 100}`
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
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["correct"] != true {
		t.Errorf("expected correct=true, got %v", resp["correct"])
	}
}

func TestSubmitFlagHMACIncorrect(t *testing.T) {
	body := `{"title": "HMAC Test 2", "category": "web", "points": 100}`
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

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if resp["correct"] == true {
		t.Error("expected incorrect flag to fail")
	}
}
