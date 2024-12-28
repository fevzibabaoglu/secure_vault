package utils

import (
	"crypto/sha256"
	"io"
	"os"
)

const (
	HashSize = sha256.Size
)

func GenerateDataHash(data []byte) []byte {
	hash := sha256.Sum256(data)
	return hash[:]
}

func GenerateFileHash(file *os.File, startPosOffset, endPosOffset int64, startWhence, endWhence int) ([]byte, error) {
	// Save the current file pointer position
	currentPos, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	// Initialize the hasher
	hasher := sha256.New()

	limitedReader, err := GetLimitedReader(file, startPosOffset, endPosOffset, startWhence, endWhence)
	if err != nil {
		return nil, err
	}

	// Stream the file content into the hasher in chunks
	_, err = io.Copy(hasher, limitedReader)
	if err != nil {
		return nil, err
	}

	// Restore the file pointer to its original position
	_, err = file.Seek(currentPos, io.SeekStart)
	if err != nil {
		return nil, err
	}

	// Return the final hash
	return hasher.Sum(nil), nil
}
