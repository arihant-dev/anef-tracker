package evidence

import (
	"crypto/sha256"
	"encoding/hex"
)

func CalculateHash(data []byte) string {
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

func VerifyPayloadHash(data []byte, expectedHash string) bool {
	if expectedHash == "" {
		return false
	}
	actualHash := CalculateHash(data)
	return actualHash == expectedHash
}
