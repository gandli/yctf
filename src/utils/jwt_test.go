package utils

import (
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	secret := "test-secret-min-32-bytes-long!!"
	userID := "user-123"
	role := "player"

	token, err := GenerateToken(userID, role, secret)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	if token == "" {
		t.Error("token must not be empty")
	}
}

func TestValidateToken(t *testing.T) {
	secret := "test-secret-min-32-bytes-long!!"
	userID := "user-123"
	role := "player"

	token, err := GenerateToken(userID, role, secret)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.UserID != userID {
		t.Errorf("expected userID %s, got %s", userID, claims.UserID)
	}
	if claims.Role != role {
		t.Errorf("expected role %s, got %s", role, claims.Role)
	}
}

func TestValidateTokenWrongSecret(t *testing.T) {
	secret := "test-secret-min-32-bytes-long!!"
	wrongSecret := "wrong-secret-min-32-bytes-long!"

	token, _ := GenerateToken("user-123", "player", secret)

	_, err := ValidateToken(token, wrongSecret)
	if err == nil {
		t.Error("token with wrong secret should fail validation")
	}
}

func TestValidateTokenExpired(t *testing.T) {
	secret := "test-secret-min-32-bytes-long!!"

	// Generate token that's already expired
	token, err := GenerateTokenWithExpiry("user-123", "player", secret, -1*time.Hour)
	if err != nil {
		t.Fatalf("GenerateTokenWithExpiry failed: %v", err)
	}

	_, err = ValidateToken(token, secret)
	if err == nil {
		t.Error("expired token should fail validation")
	}
}

func TestTokenExpiry(t *testing.T) {
	secret := "test-secret-min-32-bytes-long!!"

	token, err := GenerateTokenWithExpiry("user-123", "player", secret, 15*time.Minute)
	if err != nil {
		t.Fatalf("GenerateTokenWithExpiry failed: %v", err)
	}

	claims, err := ValidateToken(token, secret)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}

	if claims.ExpiresAt.Before(time.Now()) {
		t.Error("token should not be expired")
	}
}
