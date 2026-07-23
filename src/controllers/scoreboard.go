package controllers

import (
	"encoding/json"
	"net/http"
	"sort"
	"sync"
	"time"
)

type TeamScore struct {
	TeamID   string `json:"team_id"`
	TeamName string `json:"team_name"`
	Score    int    `json:"score"`
	Solves   int    `json:"solves"`
	Rank     int    `json:"rank"`
}

type ScoreEvent struct {
	TeamID        string `json:"team_id"`
	TeamName      string `json:"team_name"`
	ChallengeName string `json:"challenge_name"`
	Points        int    `json:"points"`
	Timestamp     string `json:"timestamp"`
}

var (
	scoreboardTeams   []*TeamScore
	scoreboardEvents  []*ScoreEvent
	scoreboardMu      sync.RWMutex
)

func ScoreboardHandler(w http.ResponseWriter, r *http.Request) {
	scoreboardMu.RLock()
	teams := make([]*TeamScore, len(scoreboardTeams))
	copy(teams, scoreboardTeams)
	scoreboardMu.RUnlock()

	// Sort by score descending
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Score > teams[j].Score
	})

	// Assign ranks
	for i, team := range teams {
		team.Rank = i + 1
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"teams":      teams,
		"updated_at": time.Now().UTC().Format(time.RFC3339),
	})
}

func ScoreboardTimelineHandler(w http.ResponseWriter, r *http.Request) {
	scoreboardMu.RLock()
	events := make([]*ScoreEvent, len(scoreboardEvents))
	copy(events, scoreboardEvents)
	scoreboardMu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"events": events,
	})
}

func UpdateScoreboard(teamID, teamName, challengeName string, points int) {
	scoreboardMu.Lock()
	defer scoreboardMu.Unlock()

	// Find or create team
	var team *TeamScore
	for _, t := range scoreboardTeams {
		if t.TeamID == teamID {
			team = t
			break
		}
	}
	if team == nil {
		team = &TeamScore{
			TeamID:   teamID,
			TeamName: teamName,
		}
		scoreboardTeams = append(scoreboardTeams, team)
	}

	team.Score += points
	team.Solves++

	// Add event
	scoreboardEvents = append(scoreboardEvents, &ScoreEvent{
		TeamID:        teamID,
		TeamName:      teamName,
		ChallengeName: challengeName,
		Points:        points,
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
	})
}
