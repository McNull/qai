package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"os"
)

// _KK_ is set at build time via ldflags
// For example: go build -ldflags="-X github.com/mcnull/github.com/mcnull/qai/shared/utils._KK_=value"
var _KK_ string = "s0m3tH1ng_v3ry_S3cur3!*$#_"

// Encode encrypts a string using the provided key or falls back to QAI_SECRET_KEY from environment.
// The function returns the base64 encoded encrypted string.
func Encode(value string, key ...string) (string, error) {
	secretKey := getSecretKey(key...)
	if secretKey == "" {
		return "", errors.New("no secret key provided and QAI_SECRET_KEY environment variable not set")
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher([]byte(padKey(secretKey)))
	if err != nil {
		return "", err
	}

	// Create a new GCM AEAD
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Create a byte array from the input string
	plaintext := []byte(value)

	// Generate a random nonce
	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Encrypt the data
	// The nonce is prepended to the ciphertext
	ciphertext := aead.Seal(nonce, nonce, plaintext, nil)

	// Return as base64 encoded string
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decode decrypts a base64 encoded string that was encrypted with Encode.
// It uses the provided key or falls back to QAI_SECRET_KEY from environment.
func Decode(encodedValue string, key ...string) (string, error) {
	secretKey := getSecretKey(key...)
	if secretKey == "" {
		return "", errors.New("no secret key provided and QAI_SECRET_KEY environment variable not set")
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher([]byte(padKey(secretKey)))
	if err != nil {
		return "", err
	}

	// Create a new GCM AEAD
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Decode the base64 encoded input
	ciphertext, err := base64.StdEncoding.DecodeString(encodedValue)
	if err != nil {
		return "", err
	}

	// Check if the ciphertext is valid
	nonceSize := aead.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	// Get the nonce from the ciphertext
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the data
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// getSecretKey retrieves the secret key from the provided argument, DefaultSecretKey, or environment variable
func getSecretKey(key ...string) string {

	// First priority: provided key
	if len(key) > 0 && key[0] != "" {
		return key[0]
	}
	// Second priority: build-time defined default key
	if _KK_ != "" {
		return _KK_
	}

	// Third priority: environment variable
	return os.Getenv("QAI_SECRET_KEY")
}

// padKey ensures the key is exactly 32 bytes (256 bits) for AES-256
func padKey(key string) []byte {
	// Target key size for AES-256
	keySize := 32

	// Convert key to bytes
	keyBytes := []byte(key)

	// If key is exactly the right length, return it
	if len(keyBytes) == keySize {
		return keyBytes
	}

	// If key is too short, pad it
	if len(keyBytes) < keySize {
		paddedKey := make([]byte, keySize)
		copy(paddedKey, keyBytes)
		// Repeat the key for padding if necessary
		for i := len(keyBytes); i < keySize; i++ {
			paddedKey[i] = keyBytes[i%len(keyBytes)]
		}
		return paddedKey
	}

	// If key is too long, truncate it
	return keyBytes[:keySize]
}
