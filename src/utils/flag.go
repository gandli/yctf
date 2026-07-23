package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateFlag(teamID, challengeID, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(teamID + ":" + challengeID))
	hash := hex.EncodeToString(h.Sum(nil))[:16]
	return "flag{" + hash + "}"
}

func ValidateFlag(flag, teamID, challengeID, secret string) bool {
	expected := GenerateFlag(teamID, challengeID, secret)
	return hmac.Equal([]byte(flag), []byte(expected))
}
