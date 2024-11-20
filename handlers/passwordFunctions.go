package handlers

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"

	"golang.org/x/crypto/argon2"
)

func hashPassword(password string, saltSize int) (string, string) {
	salt := make([]byte, saltSize)
	_, err := rand.Read(salt)
	if err != nil {
		panic(err)
	}
	hash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return base64.RawStdEncoding.EncodeToString(hash), base64.RawStdEncoding.EncodeToString(salt)
}

func verifyPassword(password, encodedHash, encodedSalt string) bool {
	salt, err := base64.RawStdEncoding.DecodeString(encodedSalt)
	if err != nil {
		return false
	}
	hash, err := base64.RawStdEncoding.DecodeString(encodedHash)
	if err != nil {
		return false
	}
	newHash := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)
	return subtle.ConstantTimeCompare(hash, newHash) == 1
}
