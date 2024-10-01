package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

// Generates a random 32-byte (256-bit) AES key for AES-256-GCM
func GenerateAES256Key() ([]byte, error) {
	key := make([]byte, 32) // 32 bytes = 256 bits for AES-256
	_, err := io.ReadFull(rand.Reader, key)
	if err != nil {
		return nil, fmt.Errorf("failed to generate key: %w", err)
	}
	return key, nil
}

// EncryptAESGCM encrypts plaintext using AES-256-GCM with a randomly generated nonce
func EncryptAESGCM(key []byte, plaintext []byte) (string, error) {
	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create a new GCM cipher based on the AES block
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Generate a nonce for AES-GCM (GCM standard nonce size is 12 bytes)
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt the data using AES-GCM
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil) // Seal appends the ciphertext to the nonce
	return hex.EncodeToString(ciphertext), nil
}

// Decrypts the ciphertext using AES-256-GCM
func DecryptAESGCM(key []byte, encrypted string) (string, error) {
	// Decode the hex-encoded ciphertext
	ciphertext, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("failed to decode ciphertext: %w", err)
	}

	// Create AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// Create a new GCM cipher based on the AES block
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Extract the nonce from the ciphertext (Nonce is prefixed to the ciphertext)
	nonceSize := aesGCM.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the data using AES-GCM
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt data: %w", err)
	}

	return string(plaintext), nil
}
