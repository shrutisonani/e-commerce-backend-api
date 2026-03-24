package utils

import (
	"crypto/rand"
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
