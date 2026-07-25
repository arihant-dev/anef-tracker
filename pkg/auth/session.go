package auth

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// DefaultEncryptionKey is used for local config encryption (can be overridden by ANEF_ENCRYPTION_KEY).
var DefaultEncryptionKey = []byte("anef-tracker-secure-key-32bytes!")

func GetSessionFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(homeDir, ".anef")
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return filepath.Join(dir, "session.json"), nil
}

// SaveSession encrypts and writes the session to ~/.anef/session.json
func SaveSession(session *CurlSession) error {
	path, err := GetSessionFilePath()
	if err != nil {
		return err
	}

	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	encrypted, err := encrypt(data, DefaultEncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to encrypt session: %w", err)
	}

	return os.WriteFile(path, []byte(hex.EncodeToString(encrypted)), 0600)
}

// LoadSession loads and decrypts the session from ~/.anef/session.json
func LoadSession() (*CurlSession, error) {
	path, err := GetSessionFilePath()
	if err != nil {
		return nil, err
	}

	hexData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	encrypted, err := hex.DecodeString(string(hexData))
	if err != nil {
		return nil, fmt.Errorf("invalid hex session file: %w", err)
	}

	data, err := decrypt(encrypted, DefaultEncryptionKey)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt session: %w", err)
	}

	var session CurlSession
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

func encrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
