package hasher

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256Hex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
