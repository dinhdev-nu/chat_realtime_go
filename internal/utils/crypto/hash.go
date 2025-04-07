package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/sha256"
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

func HashPassword(password string, salt string) string {
	hash := sha256.Sum256([]byte(password + salt))
	return hex.EncodeToString(hash[:])
}

func VerifyPassword(password string, inPassword string, salt string) bool {
	hashPassword := HashPassword(password, salt)
	return hashPassword == inPassword
}

func CreateSalt() (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}
