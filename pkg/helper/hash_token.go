package helper

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashTokenSHA256 แปลง Refresh Token (JWT String) ให้เป็น SHA-256 Hash String
func HashTokenSHA256(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
