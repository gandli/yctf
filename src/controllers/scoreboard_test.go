package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestScoreboardHandler(t *testing.T) {
	// Clear any previous state
	scoreboardTeams = nil

	// Add test teams
	scoreboardTeams = append(scoreboardTeams, &TeamScore{
		TeamID:   "team-1",
		TeamName: "Team A",
		Score:    500,
		Solves:   3,
	})
	scoreboardTeams = append(scoreboardTeams, &TeamScore{
		TeamID:   "team-2",
		TeamName: "Team B",
		Score:    300,
		Solves:   2,
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/scoreboard", nil)
	w := httptest.NewRecorder()

	ScoreboardHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	teams, ok := resp["teams"].([]interface{})
	if !ok {
		t.Error("expected teams array")
		return
	}

	if len(teams) != 2 {
		t.Errorf("expected 2 teams, got %d", len(teams))
	}

	// Verify sorted by score (descending)
	first := teams[0].(map[string]interface{})
	if first["team_name"] != "Team A" {
		t.Errorf("expected Team A first, got %v", first["team_name"])
	}
}

func TestScoreboardHandlerEmpty(t *testing.T) {
	scoreboardTeams = nil

	req := httptest.NewRequest(http.MethodGet, "/api/v1/scoreboard", nil)
	w := httptest.NewRecorder()

	ScoreboardHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	teams, ok := resp["teams"].([]interface{})
	if !ok || len(teams) != 0 {
		t.Error("expected empty teams array")
	}
}

func TestScoreboardTimelineHandler(t *testing.T) {
	scoreboardEvents = nil
	now := "2026-07-23T10:00:00Z"

	scoreboardEvents = append(scoreboardEvents, &ScoreEvent{
		TeamID:        "team-1",
		TeamName:      "Team A",
		ChallengeName: "Web 101",
		Points:        100,
		Timestamp:     now,
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/scoreboard/timeline", nil)
	w := httptest.NewRecorder()

	ScoreboardTimelineHandler(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	events, ok := resp["events"].([]interface{})
	if !ok || len(events) != 1 {
		t.Errorf("expected 1 event, got %v", resp)
	}
}

func TestScoreboardTeamRank(t *testing.T) {
	scoreboardTeams = nil

	scoreboardTeams = append(scoreboardTeams, &TeamScore{
		TeamID:   "team-1",
		TeamName: "Team A",
		Score:    500,
	})
	scoreboardTeams = append(scoreboardTeams, &TeamScore{
		TeamID:   "team-2",
		TeamName: "Team B",
		Score:    300,
	})
	scoreboardTeams = append(scoreboardTeams, &TeamScore{
		TeamID:   "team-3",
		TeamName: "Team C",
		Score:    700,
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/scoreboard", nil)
	w := httptest.NewRecorder()

	ScoreboardHandler(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	teams := resp["teams"].([]interface{})

	// Should be sorted: Team C (700), Team A (500), Team B (300)
	if teams[0].(map[string]interface{})["team_name"] != "Team C" {
		t.Errorf("expected Team C first, got %v", teams[0])
	}
	if teams[1].(map[string]interface{})["team_name"] != "Team A" {
		t.Errorf("expected Team A second, got %v", teams[1])
	}
	if teams[2].(map[string]interface{})["team_name"] != "Team B" {
		t.Errorf("expected Team B third, got %v", teams[2])
	}

	// Check ranks
	if teams[0].(map[string]interface{})["rank"].(float64) != 1 {
		t.Error("Team C should be rank 1")
	}
	if teams[1].(map[string]interface{})["rank"].(float64) != 2 {
		t.Error("Team A should be rank 2")
	}
	if teams[2].(map[string]interface{})["rank"].(float64) != 3 {
		t.Error("Team B should be rank 3")
	}
}
