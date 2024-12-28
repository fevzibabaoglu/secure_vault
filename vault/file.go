package vault

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"secure_vault/vault/utils"
	"time"
)

func AddFileToVault(v *Vault, key []byte, filePath string, deleteFile bool) error {
	// Open the file to be added
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	// Get file info
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	// Read all the file's content into a []byte
	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	// Close the file
	file.Close()

	// Encrypt the file content
	encryptedData, err := utils.Encrypt(data, key)
	if err != nil {
		return err
	}

	// Create file metadata
	fileMetadata := FileMetadata{
		Name:          stat.Name(),
		Index:         int64(len(v.FilesMetadata)),
		Offset:        int64(len(v.Files)),
		IntegrityHash: utils.GenerateDataHash(encryptedData),
		AddedAt:       time.Now().Truncate(0),
	}

	// Update the vault structure
	v.Files = append(v.Files, encryptedData...)
	v.FilesMetadata = append(v.FilesMetadata, fileMetadata)

	// Delete the original file if requested
	if deleteFile {
		err = os.Remove(filePath)
		if err != nil {
			return err
		}
	}

	return nil
}

func RemoveFileFromVault(v *Vault, fileIndex int64) error {
	// If the file doesn't exist, return an error
	if fileIndex < 0 || fileIndex >= int64(len(v.FilesMetadata)) {
		return fmt.Errorf("file index not found: %d", fileIndex)
	}

	// Get file offsets
	fileStartOffset, fileEndOffset, _ := getFileOffsets(v, fileIndex)
	removedFileSize := fileEndOffset - fileStartOffset

	// Remove the file content from the Files list
	v.Files = append(v.Files[:fileStartOffset], v.Files[fileEndOffset:]...)

	// Remove the file metadata from the FilesMetadata list
	v.FilesMetadata = append(v.FilesMetadata[:fileIndex], v.FilesMetadata[fileIndex+1:]...)

	// Recalculate the indices for the remaining files
	for i := fileIndex; i < int64(len(v.FilesMetadata)); i++ {
		v.FilesMetadata[i].Index = int64(i)
	}

	// Recalculate the offsets for the remaining files
	for i := fileIndex; i < int64(len(v.FilesMetadata)); i++ {
		v.FilesMetadata[i].Offset -= removedFileSize
	}

	return nil
}

func ExtractFileFromVault(v *Vault, key []byte, fileIndex int64, extractFolderPath string) error {
	// If the file doesn't exist, return an error
	if fileIndex < 0 || fileIndex >= int64(len(v.FilesMetadata)) {
		return fmt.Errorf("file index not found: %d", fileIndex)
	}

	// Get the file content from the vault
	fileData, _ := getFile(v, fileIndex)

	// Decrypt the file content
	file, err := utils.Decrypt(fileData, key)
	if err != nil {
		return err
	}

	// Combine the extractPath with the file name to get the full path
	fileName := v.FilesMetadata[fileIndex].Name
	outputFilePath := filepath.Join(extractFolderPath, fileName)

	// Open the file at the specified extraction path
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Write the decrypted file data to the extraction path
	_, err = outputFile.Write(file)
	if err != nil {
		return err
	}

	// Remove the file from the vault
	err = RemoveFileFromVault(v, fileIndex)
	if err != nil {
		return err
	}

	return nil
}

func CheckFileIntegrity(v *Vault, fileIndex int64) (bool, error) {
	// If the file doesn't exist, return an error
	if fileIndex < 0 || fileIndex >= int64(len(v.FilesMetadata)) {
		return false, fmt.Errorf("file index not found: %d", fileIndex)
	}

	// Get the expected hash
	fileMetadata := v.FilesMetadata[fileIndex]
	expectedHash := fileMetadata.IntegrityHash

	// Get the file
	fileData, _ := getFile(v, fileIndex)

	// Compute the real hash
	hash := utils.GenerateDataHash(fileData)

	return bytes.Equal(hash, expectedHash), nil
}

func getFileOffsets(v *Vault, fileIndex int64) (int64, int64, error) {
	// If the file doesn't exist, return an error
	if fileIndex < 0 || fileIndex >= int64(len(v.FilesMetadata)) {
		return 0, 0, fmt.Errorf("file index not found: %d", fileIndex)
	}

	if fileIndex == int64(len(v.FilesMetadata)-1) {
		return int64(v.FilesMetadata[fileIndex].Offset), int64(len(v.Files)), nil
	} else {
		return int64(v.FilesMetadata[fileIndex].Offset), int64(v.FilesMetadata[fileIndex+1].Offset), nil
	}
}

func getFile(v *Vault, fileIndex int64) ([]byte, error) {
	// If the file doesn't exist, return an error
	if fileIndex < 0 || fileIndex >= int64(len(v.FilesMetadata)) {
		return nil, fmt.Errorf("file index not found: %d", fileIndex)
	}

	fileStartOffset, fileEndOffset, _ := getFileOffsets(v, fileIndex)
	return v.Files[fileStartOffset:fileEndOffset], nil
}
