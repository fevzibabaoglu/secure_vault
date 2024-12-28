package vault

import (
	"bytes"
	"io"
	"os"
	"secure_vault/vault/utils"
	"time"
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

type Vault struct {
	Metadata      VaultMetadata  // Metadata of the vault
	FilesMetadata []FileMetadata // Metadata for files
	Files         []byte         // File content
}

type VaultMetadata struct {
	Salt      []byte    // Salt for key derivation
	CreatedAt time.Time // Creation timestamp
}

type FileMetadata struct {
	Name          string    // Filename
	Index         int64     // Index in the vault
	Offset        int64     // Offset in the vault
	IntegrityHash []byte    // Integrity hash is computed after the encryption
	AddedAt       time.Time // Timestamp when the file was added
}

func CreateVault(password string) (*Vault, error) {
	// Generate a random salt
	salt, err := utils.GenerateSalt()
	if err != nil {
		return nil, err
	}

	// Create an empty vault
	v := &Vault{
		Metadata: VaultMetadata{
			Salt:      salt,
			CreatedAt: time.Now().Truncate(0),
		},
		FilesMetadata: []FileMetadata{},
		Files:         []byte{},
	}

	return v, nil
}

func SaveVault(v *Vault, key []byte, vaultPath string) error {
	// Create the vault file
	vaultFile, err := os.Create(vaultPath)
	if err != nil {
		return err
	}
	defer vaultFile.Close()

	// Save the metadata
	err = writeVaultMetadata(vaultFile, &v.Metadata)
	if err != nil {
		return err
	}

	// Save the files metadata
	err = writeFilesMetadata(vaultFile, key, v.FilesMetadata)
	if err != nil {
		return err
	}

	// Save the files
	err = writeFiles(vaultFile, v.Files)
	if err != nil {
		return err
	}

	// Save the vault hash
	err = writeVaultHash(vaultFile)
	return err
}

func LoadVault(password, vaultPath string) (*Vault, error) {
	// Open the vault file
	vaultFile, err := os.Open(vaultPath)
	if err != nil {
		return nil, err
	}
	defer vaultFile.Close()

	// Load the metadata
	metadata, err := readVaultMetadata(vaultFile)
	if err != nil {
		return nil, err
	}

	// Derive the key using the password and salt
	key := utils.DeriveKey(password, metadata.Salt)

	// Load the files metadata
	filesMetadata, err := readFilesMetadata(vaultFile, key)
	if err != nil {
		return nil, err
	}

	// Load the files
	files, err := readFiles(vaultFile)
	if err != nil {
		return nil, err
	}

	// Reconstruct the vault structure
	v := &Vault{
		Metadata:      *metadata,
		FilesMetadata: filesMetadata,
		Files:         files,
	}

	return v, nil
}

func CheckVaultIntegrity(vaultPath string) (bool, error) {
	// Open the vault file
	vaultFile, err := os.Open(vaultPath)
	if err != nil {
		return false, err
	}
	defer vaultFile.Close()

	// Load the vault hash
	expectedVaultHash, err := readVaultHash(vaultFile)
	if err != nil {
		return false, err
	}

	// Verify the data integrity by recomputing the hash
	vaultHash, err := utils.GenerateFileHash(vaultFile, 0, -utils.HashSize, io.SeekStart, io.SeekEnd)
	if err != nil {
		return false, err
	}

	return bytes.Equal(vaultHash, expectedVaultHash), nil
}
