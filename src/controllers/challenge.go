package controllers

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gandli/yctf/db/models"
	"github.com/google/uuid"
)

var (
	challenges = make(map[string]*models.Challenge)
	mu         sync.RWMutex
)

type CreateChallengeRequest struct {
	Title          string                 `json:"title"`
	Description    string                 `json:"description"`
	Category       string                 `json:"category"`
	Points         int                    `json:"points"`
	ContainerImage string                 `json:"container_image"`
	ContainerConfig map[string]interface{} `json:"container_config"`
}

func CreateChallengeHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateChallengeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid_json","message":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if req.Title == "" {
		http.Error(w, `{"error":"missing_title","message":"Title is required"}`, http.StatusBadRequest)
		return
	}

	chal := models.NewChallenge(req.Title, req.Category, req.Points)
	chal.ID = uuid.New().String()
	if err := chal.ValidateCategory(); err != nil {
		http.Error(w, `{"error":"invalid_category","message":"`+err.Error()+`"}`, http.StatusBadRequest)
		return
	}

	if req.Points <= 0 {
		http.Error(w, `{"error":"invalid_points","message":"Points must be positive"}`, http.StatusBadRequest)
		return
	}

	chal.Description = req.Description
	chal.ContainerImage = req.ContainerImage

	mu.Lock()
	challenges[chal.ID] = chal
	mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chal)
}

func ListChallengesHandler(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")

	mu.RLock()
	list := make([]*models.Challenge, 0, len(challenges))
	for _, chal := range challenges {
		if category == "" || chal.Category == category {
			list = append(list, chal)
		}
	}
	mu.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"challenges": list,
		"total":      len(list),
	})
}

func GetChallengeHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[len("/api/v1/challenges/"):]

	mu.RLock()
	chal, ok := challenges[id]
	mu.RUnlock()

	if !ok {
		http.Error(w, `{"error":"not_found","message":"Challenge not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chal)
}
