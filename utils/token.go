package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
)

func GenerateToken() (string, error) {
	str := make([]byte, 32)
	_, err := rand.Read(str)

	if err != nil {
		// Return an empty string
		return "", err
	}

	return hex.EncodeToString(str), nil
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}
