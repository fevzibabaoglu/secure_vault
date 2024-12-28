package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

const (
	CipherBlockSize = aes.BlockSize
)

func Encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Initialize a nonce (IV) for CTR mode (16 bytes for AES)
	nonce := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// Create the CTR stream cipher
	stream := cipher.NewCTR(block, nonce)

	// Encrypt the data
	ciphertext := make([]byte, len(data))
	stream.XORKeyStream(ciphertext, data)

	// Prepend the nonce to the ciphertext for use during decryption
	return append(nonce, ciphertext...), nil
}

func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Extract the nonce from the beginning of the ciphertext
	nonce, ciphertext := ciphertext[:aes.BlockSize], ciphertext[aes.BlockSize:]

	// Create the CTR stream cipher
	stream := cipher.NewCTR(block, nonce)

	// Decrypt the data
	plaintext := make([]byte, len(ciphertext))
	stream.XORKeyStream(plaintext, ciphertext)

	return plaintext, nil
}
