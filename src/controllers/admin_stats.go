package controllers

import (
	"encoding/json"
	"net/http"
	"time"
)

var startTime = time.Now()

func AdminStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	mu.RLock()
	challengeCount := len(challenges)
	mu.RUnlock()

	writeupMu.RLock()
	writeupCount := len(writeups)
	writeupMu.RUnlock()

	resp := map[string]interface{}{
		"users":       0,
		"teams":       0,
		"challenges":  challengeCount,
		"submissions": 0,
		"containers":  0,
		"writeups":    writeupCount,
		"uptime":      time.Since(startTime).String(),
	}

	json.NewEncoder(w).Encode(resp)
}
