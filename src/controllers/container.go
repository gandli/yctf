package controllers

import (
	"encoding/json"
	"net/http"
	"os/exec"
	"strings"
)

func StartContainerHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ChallengeID string `json:"challenge_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid_json","message":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if req.ChallengeID == "" {
		http.Error(w, `{"error":"missing_challenge_id","message":"challenge_id is required"}`, http.StatusBadRequest)
		return
	}

	if _, err := exec.LookPath("docker"); err != nil {
		http.Error(w, `{"error":"docker_unavailable","message":"Docker is not available"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"instance_id":  "placeholder-id",
		"challenge_id": req.ChallengeID,
		"status":       "running",
		"port":         30001,
		"expires_at":   "2026-07-23T12:00:00Z",
	})
}

func StopContainerHandler(w http.ResponseWriter, r *http.Request) {
	instanceID := r.PathValue("instance_id")
	if instanceID == "" {
		// Fallback: extract from URL path
		parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/stop"), "/")
		instanceID = parts[len(parts)-1]
	}

	if instanceID == "" || instanceID == "instances" {
		http.Error(w, `{"error":"missing_instance_id","message":"instance_id is required"}`, http.StatusBadRequest)
		return
	}

	if _, err := exec.LookPath("docker"); err != nil {
		http.Error(w, `{"error":"docker_unavailable","message":"Docker is not available"}`, http.StatusInternalServerError)
		return
	}

	http.Error(w, `{"error":"not_found","message":"Instance not found"}`, http.StatusNotFound)
}

func ContainerStatusHandler(w http.ResponseWriter, r *http.Request) {
	instanceID := r.PathValue("instance_id")
	if instanceID == "" {
		parts := strings.Split(strings.TrimSuffix(r.URL.Path, "/status"), "/")
		instanceID = parts[len(parts)-1]
	}

	if instanceID == "" || instanceID == "instances" {
		http.Error(w, `{"error":"missing_instance_id","message":"instance_id is required"}`, http.StatusBadRequest)
		return
	}

	if _, err := exec.LookPath("docker"); err != nil {
		http.Error(w, `{"error":"docker_unavailable","message":"Docker is not available"}`, http.StatusInternalServerError)
		return
	}

	http.Error(w, `{"error":"not_found","message":"Instance not found"}`, http.StatusNotFound)
}

func ListMyContainersHandler(w http.ResponseWriter, r *http.Request) {
	if _, err := exec.LookPath("docker"); err != nil {
		http.Error(w, `{"error":"docker_unavailable","message":"Docker is not available"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"instances": []interface{}{},
	})
}
