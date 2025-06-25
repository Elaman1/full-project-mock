package service

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateTraceID() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", fmt.Errorf("error generating traceid: %v", err)
	}

	return hex.EncodeToString(b), nil
}
