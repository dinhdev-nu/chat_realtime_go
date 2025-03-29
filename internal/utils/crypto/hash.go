package crypto

import (
	"crypto"
	"encoding/hex"
)

func HashEmail(email string) string {
	hash := crypto.SHA256.New()
	hash.Write([]byte(email))
	return hex.EncodeToString(hash.Sum(nil))
}

func Compare(str string, hash string) bool {
	return HashEmail(str) == hash
}
