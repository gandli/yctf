package models

import (
	"fmt"
	"regexp"
	"strings"
)

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
	Role         string `json:"role"`
	TeamID       string `json:"team_id"`
	IsBanned     bool   `json:"is_banned"`
	LastLogin    string `json:"last_login"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func NewUser(username, email, passwordHash, role string) *User {
	return &User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
	}
}

func (u *User) ValidateRole() error {
	switch u.Role {
	case "admin", "author", "player":
		return nil
	default:
		return fmt.Errorf("invalid role: %s", u.Role)
	}
}

func (u *User) ValidateEmail() error {
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return fmt.Errorf("invalid email format: %s", u.Email)
	}
	return nil
}

func (u *User) ValidateUsername() error {
	if len(u.Username) < 3 || len(u.Username) > 32 {
		return fmt.Errorf("username must be 3-32 characters")
	}
	usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !usernameRegex.MatchString(u.Username) {
		return fmt.Errorf("username can only contain alphanumeric and underscore")
	}
	return nil
}

func (u *User) FullName() string {
	return strings.TrimSpace(u.Username)
}
