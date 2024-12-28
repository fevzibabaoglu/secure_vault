package utils

import (
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

const (
	saltLength  = 16        // Salt length in bytes
	timeCost    = 3         // Number of iterations
	memoryCost  = 64 * 1024 // Memory in KiB (64 MiB)
	parallelism = 2         // Number of threads
	keyLength   = 32        // Key length in bytes
)

// GenerateSalt creates a random salt for Argon2.
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

// DeriveKey derives a key using Argon2id.
func DeriveKey(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, timeCost, memoryCost, parallelism, keyLength)
}
