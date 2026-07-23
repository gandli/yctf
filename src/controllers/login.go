package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gandli/yctf/utils"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid_json","message":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, `{"error":"missing_fields","message":"Username and password are required"}`, http.StatusBadRequest)
		return
	}

	secret := "dev-secret-min-32-bytes-long!!"
	token, err := utils.GenerateToken(req.Username, "player", secret)
	if err != nil {
		http.Error(w, `{"error":"internal","message":"Failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := map[string]interface{}{
		"token": token,
		"user": map[string]string{
			"id":       req.Username,
			"username": req.Username,
			"role":     "player",
		},
	}
	json.NewEncoder(w).Encode(resp)
}
