package controllers

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/google/uuid"
)

type Writeup struct {
	ID           string `json:"id"`
	ChallengeID  string `json:"challenge_id"`
	TeamID       string `json:"team_id"`
	UserID       string `json:"user_id"`
	URL          string `json:"url"`
	Content      string `json:"content"`
	IsApproved   bool   `json:"is_approved"`
	ReviewedBy   string `json:"reviewed_by"`
	Score        int    `json:"score"`
}

var (
	writeups   = make(map[string]*Writeup)
	writeupMu  sync.RWMutex
)

type SubmitWriteupRequest struct {
	ChallengeID string `json:"challenge_id"`
	URL         string `json:"url"`
	Content     string `json:"content"`
}

func SubmitWriteupHandler(w http.ResponseWriter, r *http.Request) {
	var req SubmitWriteupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid_json","message":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if req.ChallengeID == "" {
		http.Error(w, `{"error":"missing_challenge_id","message":"challenge_id is required"}`, http.StatusBadRequest)
		return
	}

	if req.URL == "" && req.Content == "" {
		http.Error(w, `{"error":"missing_content","message":"url or content is required"}`, http.StatusBadRequest)
		return
	}

	writeup := &Writeup{
		ID:          uuid.New().String(),
		ChallengeID: req.ChallengeID,
		TeamID:      "team-placeholder",
		UserID:      "user-placeholder",
		URL:         req.URL,
		Content:     req.Content,
		IsApproved:  false,
	}

	writeupMu.Lock()
	writeups[writeup.ID] = writeup
	writeupMu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(writeup)
}

func ListWriteupsHandler(w http.ResponseWriter, r *http.Request) {
	writeupMu.RLock()
	list := make([]*Writeup, 0, len(writeups))
	for _, w := range writeups {
		list = append(list, w)
	}
	writeupMu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"writeups": list,
	})
}

func ApproveWriteupHandler(w http.ResponseWriter, r *http.Request) {
	writeupID := r.URL.Path[len("/api/v1/writeups/"):]
	writeupID = writeupID[:len(writeupID)-len("/approve")]

	writeupMu.Lock()
	writeup, ok := writeups[writeupID]
	if ok {
		writeup.IsApproved = true
		writeup.ReviewedBy = "admin-placeholder"
	}
	writeupMu.Unlock()

	if !ok {
		http.Error(w, `{"error":"not_found","message":"Writeup not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(writeup)
}
