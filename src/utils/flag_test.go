package utils

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
)

func TestGenerateFlag(t *testing.T) {
	secret := "test-secret"
	teamID := "team-123"
	challengeID := "chal-456"

	flag := GenerateFlag(teamID, challengeID, secret)

	if !strings.HasPrefix(flag, "flag{") {
		t.Errorf("flag must start with 'flag{', got %s", flag)
	}
	if !strings.HasSuffix(flag, "}") {
		t.Errorf("flag must end with '}', got %s", flag)
	}
	flag2 := GenerateFlag(teamID, challengeID, secret)
	if flag != flag2 {
		t.Errorf("same team+challenge must produce same flag, got %s vs %s", flag, flag2)
	}
	flag3 := GenerateFlag("team-999", challengeID, secret)
	if flag == flag3 {
		t.Errorf("different teams must produce different flags")
	}
	flag4 := GenerateFlag(teamID, "chal-999", secret)
	if flag == flag4 {
		t.Errorf("different challenges must produce different flags")
	}
}

func TestValidateFlag(t *testing.T) {
	secret := "test-secret"
	teamID := "team-123"
	challengeID := "chal-456"

	flag := GenerateFlag(teamID, challengeID, secret)

	if !ValidateFlag(flag, teamID, challengeID, secret) {
		t.Error("valid flag should pass validation")
	}
	if ValidateFlag("flag{wrong}", teamID, challengeID, secret) {
		t.Error("invalid flag should fail validation")
	}
	if ValidateFlag("", teamID, challengeID, secret) {
		t.Error("empty flag should fail validation")
	}
}

func TestFlagFormat(t *testing.T) {
	secret := "test-secret"
	flag := GenerateFlag("t1", "c1", secret)

	inner := strings.TrimPrefix(flag, "flag{")
	inner = strings.TrimSuffix(inner, "}")

	if _, err := hex.DecodeString(inner); err != nil {
		t.Errorf("flag content must be hex, got %s", inner)
	}
	if len(inner) != 16 {
		t.Errorf("flag content must be 16 hex chars, got %d", len(inner))
	}
}

func TestHMACUniqueness(t *testing.T) {
	secret := "test-secret"
	flags := make(map[string]bool)

	for i := 0; i < 1000; i++ {
		flag := GenerateFlag(fmt.Sprintf("team-%d", i), "challenge-1", secret)
		if flags[flag] {
			t.Errorf("duplicate flag generated: %s", flag)
		}
		flags[flag] = true
	}
}
