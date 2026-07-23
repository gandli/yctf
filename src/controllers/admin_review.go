package controllers

import (
	"encoding/json"
	"net/http"
)

type Submission struct {
	ID            string `json:"id"`
	UserID        string `json:"user_id"`
	TeamID        string `json:"team_id"`
	ChallengeID   string `json:"challenge_id"`
	FlagSubmitted string `json:"flag_submitted"`
	IsCorrect     bool   `json:"is_correct"`
	IPAddress     string `json:"ip_address"`
	SubmittedAt   string `json:"submitted_at"`
}

var submissions = make(map[string]*Submission)

func AdminSubmissionsAuditLogHandler(w http.ResponseWriter, r *http.Request) {
	subList := make([]*Submission, 0, len(submissions))
	for _, s := range submissions {
		subList = append(subList, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"submissions": subList,
	})
}

func AdminWriteupReviewQueueHandler(w http.ResponseWriter, r *http.Request) {
	var pending []*Writeup
	writeupMu.RLock()
	for _, w := range writeups {
		if !w.IsApproved && w.ReviewedBy == "" {
			pending = append(pending, w)
		}
	}
	writeupMu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"pending": pending,
	})
}

func AdminScoreWriteupHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		WriteupID string `json:"writeup_id"`
		Score     int    `json:"score"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid_json"}`, http.StatusBadRequest)
		return
	}

	writeupMu.Lock()
	writeup, ok := writeups[req.WriteupID]
	if ok {
		writeup.Score = req.Score
		writeup.IsApproved = true
		writeup.ReviewedBy = "admin"
	}
	writeupMu.Unlock()

	if !ok {
		http.Error(w, `{"error":"not_found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(writeup)
}
