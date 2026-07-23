package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebSocketBroadcastScoreUpdate(t *testing.T) {
	// Clear previous state
	scoreboardMu.Lock()
	scoreboardTeams = nil
	scoreboardMu.Unlock()

	body := `{"title": "WS Test", "category": "web", "points": 100}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/challenges", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	CreateChallengeHandler(w, req)

	UpdateScoreboard("team-1", "Team A", "WS Test", 100)

	req = httptest.NewRequest(http.MethodGet, "/api/v1/scoreboard", nil)
	w = httptest.NewRecorder()
	ScoreboardHandler(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	teams := resp["teams"].([]interface{})
	if len(teams) != 1 {
		t.Errorf("expected 1 team, got %d", len(teams))
	}
}

func TestWebSocketBroadcastChallengeSolved(t *testing.T) {
	scoreboardMu.Lock()
	scoreboardTeams = nil
	scoreboardEvents = nil
	scoreboardMu.Unlock()

	body := `{"title": "WS Test 2", "category": "web", "points": 200}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/challenges", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	CreateChallengeHandler(w, req)

	UpdateScoreboard("team-2", "Team B", "WS Test 2", 200)

	req = httptest.NewRequest(http.MethodGet, "/api/v1/scoreboard/timeline", nil)
	w = httptest.NewRecorder()
	ScoreboardTimelineHandler(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	events := resp["events"].([]interface{})
	if len(events) == 0 {
		t.Error("expected at least 1 event in timeline")
	}
}

func TestWebSocketMultipleTeams(t *testing.T) {
	scoreboardMu.Lock()
	scoreboardTeams = nil
	scoreboardEvents = nil
	scoreboardMu.Unlock()

	UpdateScoreboard("team-3", "Team C", "Chal A", 150)
	UpdateScoreboard("team-4", "Team D", "Chal B", 300)
	UpdateScoreboard("team-5", "Team E", "Chal C", 50)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/scoreboard", nil)
	w := httptest.NewRecorder()
	ScoreboardHandler(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	teams := resp["teams"].([]interface{})
	if len(teams) < 3 {
		t.Errorf("expected at least 3 teams, got %d", len(teams))
	}

	first := teams[0].(map[string]interface{})
	if first["score"].(float64) < teams[len(teams)-1].(map[string]interface{})["score"].(float64) {
		t.Error("scoreboard should be sorted by score descending")
	}
}
