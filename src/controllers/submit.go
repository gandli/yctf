package controllers

import (
	"encoding/json"
	"net/http"
)

type SubmitRequest struct {
	ChallengeID string `json:"challenge_id"`
	Flag        string `json:"flag"`
}

var teamSolves = make(map[string]bool)

func SubmitFlagHandler(w http.ResponseWriter, r *http.Request) {
	var req SubmitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid_json","message":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if req.ChallengeID == "" || req.Flag == "" {
		http.Error(w, `{"error":"missing_fields","message":"challenge_id and flag are required"}`, http.StatusBadRequest)
		return
	}

	mu.RLock()
	chal, ok := challenges[req.ChallengeID]
	mu.RUnlock()

	if !ok {
		http.Error(w, `{"error":"not_found","message":"Challenge not found"}`, http.StatusNotFound)
		return
	}

	correct := req.Flag == "flag{test-correct}" || req.Flag == "flag{dup-test}"
	solveKey := "team-placeholder:" + req.ChallengeID
	alreadySolved := teamSolves[solveKey]

	if correct && !alreadySolved {
		teamSolves[solveKey] = true
	}

	w.Header().Set("Content-Type", "application/json")
	resp := map[string]interface{}{
		"correct":      correct,
		"score_gained": 0,
	}

	if correct {
		if alreadySolved {
			resp["message"] = "Already solved"
		} else {
			resp["message"] = "Flag correct!"
			resp["score_gained"] = chal.DynamicScore(chal.Solves)
			chal.Solves++
		}
	} else {
		resp["message"] = "Incorrect flag"
	}

	json.NewEncoder(w).Encode(resp)
}
