package utils

import (
	"fmt"
	"io"
	"os"
)

// Returns a limited reader between starting and ending positions
func GetLimitedReader(file *os.File, startPosOffset, endPosOffset int64, startWhence, endWhence int) (io.Reader, error) {
	// Save the current file pointer position
	currentPos, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	// Find end position
	endPos, err := file.Seek(endPosOffset, endWhence)
	if err != nil {
		return nil, err
	}

	// Find and seek start position
	startPos, err := file.Seek(startPosOffset, startWhence)
	if err != nil {
		return nil, err
	}

	// Calculate the number of bytes to read
	length := endPos - startPos
	if length <= 0 {
		// Restore the file pointer to its original position
		_, err = file.Seek(currentPos, io.SeekStart)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("endPos must be greater than startPos")
	}

	// Limit the reader to the specified range
	return io.LimitReader(file, length), nil
}
