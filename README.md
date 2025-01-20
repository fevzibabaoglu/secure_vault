# Secure Vault

Secure Vault is a robust application designed to securely store and manage sensitive files in an encrypted vault. The vault ensures strong protection against unauthorized access, even if the underlying storage medium is compromised.

---

## Features

### Core Functionalities
- **Vault Creation:** Users can create a secure vault protected by a password.
- **File Management:**
  - Add files to the vault with automatic encryption.
  - View a list of stored files (metadata only).
  - Remove or extract files with decryption.
- **Vault Locking and Unlocking:** Lock the vault to prevent unauthorized access and unlock it with the correct password.
- **File and Vault Integrity Checking:** Detect tampering using SHA-256 hashes.

### Security Highlights
- **Strong Encryption:** AES-256 in CTR mode for encrypting file content.
- **Secure Key Derivation:** Argon2 is used for deriving encryption keys from user passwords with a random salt.
- **Data Integrity:** SHA-256 ensures file and vault integrity.
- **Password Security:** Passwords are never stored.

### Vault Structure
The vault file consists of:
1. **Vault Metadata**: Vault-specific details like creation time and salt.
2. **File Metadata \[Encrypted]**: Details about stored files, such as names, offsets, and hashes.
3. **Files \[Encrypted]**: File content stored in a contiguous encrypted format.
4. **Integrity Hash**: A SHA-256 hash of the entire vault to ensure its integrity.

---

## Getting Started

### Prerequisites
- Go 1.21.4 or later installed.

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/your-repo/secure-vault.git
   cd secure-vault
   ```
2. Install dependencies (if applicable):
   ```bash
   go mod tidy
   ```

### Usage
1. Run the application:
   ```bash
   go run .
   ```
2. Follow the UI prompts to create and manage your secure vault.

---

## License
This project is open-source and released under the [MIT License](LICENSE).
