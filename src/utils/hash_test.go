package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "SecurePass123!"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash == "" {
		t.Error("hash must not be empty")
	}

	if hash == password {
		t.Error("hash must not equal plaintext password")
	}
}

func TestVerifyPassword(t *testing.T) {
	password := "SecurePass123!"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if !VerifyPassword(password, hash) {
		t.Error("correct password should verify")
	}

	if VerifyPassword("wrong-password", hash) {
		t.Error("wrong password should not verify")
	}
}

func TestPasswordHashUniqueness(t *testing.T) {
	password := "SecurePass123!"

	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)

	if hash1 == hash2 {
		t.Error("same password should produce different hashes (salt)")
	}
}
