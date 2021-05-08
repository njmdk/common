package crypt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func hmacSha256(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
