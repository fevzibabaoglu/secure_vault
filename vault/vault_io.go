package vault

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"secure_vault/vault/utils"
	"unsafe"
)

/*
	Vault file has the following structure in disk:
	Vault Metadata Size 	int64 (LittleEndian)
	Vault Metadata 			VaultMetadata
	Files Metadata Size		int64 (LittleEndian)			[Encrypted]
	Files Metadata			[]FileMetadata					[Encrypted]
	Files					[]byte (dumped back to back)	[Encrypted]
	Vault Integriy Hash		SHA256
*/

func writeVaultMetadata(vaultFile *os.File, metadata *VaultMetadata) error {
	// Serialize the vault metadata
	metadataBytes, err := utils.EncodeDataToBytes(*metadata)
	if err != nil {
		return err
	}

	// Write the size of the vault metadata
	err = binary.Write(vaultFile, binary.LittleEndian, int32(len(metadataBytes)))
	if err != nil {
		return err
	}

	// Write the vault metadata
	_, err = vaultFile.Write(metadataBytes)
	return err
}

func writeFilesMetadata(vaultFile *os.File, key []byte, filesMetadata []FileMetadata) error {
	// Serialize the files metadata
	filesMetadataBytes, err := utils.EncodeDataToBytes(filesMetadata)
	if err != nil {
		return err
	}

	// Encrypt the files metadata
	encryptedFilesMetadata, err := utils.Encrypt(filesMetadataBytes, key)
	if err != nil {
		return err
	}

	// Get the size of the file metadata
	encryptedFilesMetadataSize := int32(len(encryptedFilesMetadata))

	// Encode the size of the file metadata
	encryptedFilesMetadataSizeBytes, err := utils.EncodeInt32ToBytes(encryptedFilesMetadataSize)
	if err != nil {
		return err
	}

	// Encrypt the size of the files metadata
	encryptedFilesMetadataEncryptedSize, err := utils.Encrypt(encryptedFilesMetadataSizeBytes, key)
	if err != nil {
		return err
	}

	// Write the size of the files metadata
	_, err = vaultFile.Write(encryptedFilesMetadataEncryptedSize)
	if err != nil {
		return err
	}

	// Write the files metadata
	_, err = vaultFile.Write(encryptedFilesMetadata)
	return err
}

func writeFiles(vaultFile *os.File, files []byte) error {
	// Write the files
	_, err := vaultFile.Write(files)
	return err
}

func writeVaultHash(vaultFile *os.File) error {
	// Generate the vault file hash
	vaultHash, err := utils.GenerateFileHash(vaultFile, 0, 0, io.SeekStart, io.SeekEnd)
	if err != nil {
		return err
	}

	// Write the vault file hash
	_, err = vaultFile.Write(vaultHash)
	return err
}

func readVaultMetadata(vaultFile *os.File) (*VaultMetadata, error) {
	// Read the size of the metadata
	var metadataSize int32
	err := binary.Read(vaultFile, binary.LittleEndian, &metadataSize)
	if err != nil {
		return nil, err
	}

	// Check incorrect size
	if metadataSize < 0 {
		return nil, fmt.Errorf("error while reading vault metadata")
	}

	// Read the unencrypted metadata
	metadataBytes := make([]byte, metadataSize)
	_, err = vaultFile.Read(metadataBytes)
	if err != nil {
		return nil, err
	}

	// Decode the metadata
	var metadata VaultMetadata
	err = utils.DecodeDataFromBytes(metadataBytes, &metadata)
	if err != nil {
		return nil, err
	}

	return &metadata, nil
}

func readFilesMetadata(vaultFile *os.File, key []byte) ([]FileMetadata, error) {
	// Read the size of the files metadata
	encryptedFilesMetadataEncryptedSize := make([]byte, utils.CipherBlockSize+unsafe.Sizeof(int32(0)))
	_, err := vaultFile.Read(encryptedFilesMetadataEncryptedSize)
	if err != nil {
		return nil, err
	}

	// Decrypt the size of the file metadata
	encryptedFilesMetadataSizeBytes, err := utils.Decrypt(encryptedFilesMetadataEncryptedSize, key)
	if err != nil {
		return nil, err
	}

	// Decode the size of the file metadata
	var encryptedFilesMetadataSize int32
	err = utils.DecodeInt32FromBytes(encryptedFilesMetadataSizeBytes, &encryptedFilesMetadataSize)
	if err != nil {
		return nil, err
	}

	// Check incorrect size
	if encryptedFilesMetadataSize < 0 {
		return nil, fmt.Errorf("error while reading vault metadata")
	}

	// Read the encrypted files metadata
	encryptedFilesMetadata := make([]byte, encryptedFilesMetadataSize)
	_, err = vaultFile.Read(encryptedFilesMetadata)
	if err != nil {
		return nil, err
	}

	// Decrypt the files metadata
	filesMetadataBytes, err := utils.Decrypt(encryptedFilesMetadata, key)
	if err != nil {
		return nil, err
	}

	// Decode the files metadata
	var filesMetadata []FileMetadata
	err = utils.DecodeDataFromBytes(filesMetadataBytes, &filesMetadata)
	if err != nil {
		return nil, err
	}

	return filesMetadata, nil
}

func readFiles(vaultFile *os.File) ([]byte, error) {
	// Determine the size of the files data
	filesStartOffset, err := vaultFile.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}
	stat, err := vaultFile.Stat()
	if err != nil {
		return nil, err
	}
	totalVaultSize := stat.Size()
	filesSize := totalVaultSize - filesStartOffset - int64(utils.HashSize)

	// Read the files
	files := make([]byte, filesSize)
	_, err = vaultFile.Read(files)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func readVaultHash(vaultFile *os.File) ([]byte, error) {
	// Save the current file pointer position
	currentPos, err := vaultFile.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	// Determine the size of the data
	stat, err := vaultFile.Stat()
	if err != nil {
		return nil, err
	}
	totalVaultSize := stat.Size()
	dataSize := totalVaultSize - int64(utils.HashSize)

	// Seek to the start of the hash part
	_, err = vaultFile.Seek(dataSize, io.SeekStart)
	if err != nil {
		return nil, err
	}

	// Read the vault hash
	vaultHash := make([]byte, int64(utils.HashSize))
	_, err = vaultFile.Read(vaultHash)
	if err != nil {
		return nil, err
	}

	// Restore the file pointer to its original position
	_, err = vaultFile.Seek(currentPos, io.SeekStart)
	if err != nil {
		return nil, err
	}

	return vaultHash, nil
}
