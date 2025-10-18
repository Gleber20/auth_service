package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateHash(str string) (string, error) {
	hasher := sha256.New()
	if _, err := hasher.Write([]byte(str)); err != nil {
		return "", err
	}
	hashbytes := hasher.Sum(nil)
	return hex.EncodeToString(hashbytes), nil
}
