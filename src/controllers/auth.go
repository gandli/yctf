package controllers

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strings"

	"github.com/gandli/yctf/utils"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Token string `json:"token"`
	User  struct {
		ID       string `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
	} `json:"user"`
}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid_json","message":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	// Validate username
	if len(req.Username) < 3 || len(req.Username) > 32 {
		http.Error(w, `{"error":"invalid_username","message":"Username must be 3-32 characters"}`, http.StatusBadRequest)
		return
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !usernameRegex.MatchString(req.Username) {
		http.Error(w, `{"error":"invalid_username","message":"Username can only contain alphanumeric and underscore"}`, http.StatusBadRequest)
		return
	}

	// Validate email
	if req.Email == "" {
		http.Error(w, `{"error":"invalid_email","message":"Email is required"}`, http.StatusBadRequest)
		return
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		http.Error(w, `{"error":"invalid_email","message":"Invalid email format"}`, http.StatusBadRequest)
		return
	}

	// Validate password
	if len(req.Password) < 8 {
		http.Error(w, `{"error":"invalid_password","message":"Password must be at least 8 characters"}`, http.StatusBadRequest)
		return
	}

	// Hash password
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		http.Error(w, `{"error":"internal","message":"Failed to process password"}`, http.StatusInternalServerError)
		return
	}

	// Generate JWT token
	secret := "dev-secret-min-32-bytes-long!!"
	token, err := utils.GenerateToken(req.Username, "player", secret)
	if err != nil {
		http.Error(w, `{"error":"internal","message":"Failed to generate token"}`, http.StatusInternalServerError)
		return
	}

	// Build response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	resp := map[string]interface{}{
		"token": token,
		"user": map[string]string{
			"id":       req.Username, // Simplified for now
			"username": req.Username,
			"email":    req.Email,
			"role":     "player",
		},
	}
	json.NewEncoder(w).Encode(resp)
	_ = hash
	_ = strings.TrimSpace
}
