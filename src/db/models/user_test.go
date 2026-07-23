package models

import (
	"strings"
	"testing"
)

func TestNewUser(t *testing.T) {
	email := "player@example.com"
	username := "player1"
	password := "hashed-password"
	role := "player"

	user := NewUser(username, email, password, role)

	if user.Username != username {
		t.Errorf("expected username %s, got %s", username, user.Username)
	}
	if user.Email != email {
		t.Errorf("expected email %s, got %s", email, user.Email)
	}
	if user.PasswordHash != password {
		t.Errorf("expected password hash %s, got %s", password, user.PasswordHash)
	}
	if user.Role != role {
		t.Errorf("expected role %s, got %s", role, user.Role)
	}
}

func TestUserValidateRole(t *testing.T) {
	validRoles := []string{"admin", "author", "player"}
	for _, role := range validRoles {
		user := &User{Role: role}
		if err := user.ValidateRole(); err != nil {
			t.Errorf("role %s should be valid: %v", role, err)
		}
	}

	invalidRoles := []string{"superuser", "root", "", "ADMIN"}
	for _, role := range invalidRoles {
		user := &User{Role: role}
		if err := user.ValidateRole(); err == nil {
			t.Errorf("role %q should be invalid", role)
		}
	}
}

func TestUserValidateEmail(t *testing.T) {
	validEmails := []string{"a@b.com", "user@test.org", "name+tag@domain.co"}
	for _, email := range validEmails {
		user := &User{Email: email}
		if err := user.ValidateEmail(); err != nil {
			t.Errorf("email %s should be valid: %v", email, err)
		}
	}

	invalidEmails := []string{"", "not-email", "@nodomain", "spaces in@email.com"}
	for _, email := range invalidEmails {
		user := &User{Email: email}
		if err := user.ValidateEmail(); err == nil {
			t.Errorf("email %q should be invalid", email)
		}
	}
}

func TestUserValidateUsername(t *testing.T) {
	validUsernames := []string{"abc", "player1", "user_name", "a1b2c3d4e5f6g7h8"}
	for _, username := range validUsernames {
		user := &User{Username: username}
		if err := user.ValidateUsername(); err != nil {
			t.Errorf("username %s should be valid: %v", username, err)
		}
	}

	invalidUsernames := []string{"", "ab", "waytoolongusernamethatexceedslimit", "has space", "has!char"}
	for _, username := range invalidUsernames {
		user := &User{Username: username}
		if err := user.ValidateUsername(); err == nil {
			t.Errorf("username %q should be invalid", username)
		}
	}
}

func TestUserFullName(t *testing.T) {
	user := &User{Username: "player1"}
	expected := "player1"
	if got := user.FullName(); got != expected {
		t.Errorf("expected FullName %s, got %s", expected, got)
	}
	_ = strings.TrimSpace
}
