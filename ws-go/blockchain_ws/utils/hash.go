package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// HashString Return the Hash256 of the target string
func HashString(target string) string {
	sum := sha256.Sum256([]byte(target))
	return hex.EncodeToString(sum[:])
}
