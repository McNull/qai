package utils

import (
	"os"
	"testing"
)

func TestEncodeDecodeWithProvidedKey(t *testing.T) {
	// Test data
	originalText := "Hello World"
	key := "test_secret_key"

	// Encode the text
	encoded, err := Encode(originalText, key)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	if encoded == "" {
		t.Fatal("Encoded string is empty")
	}
	if encoded == originalText {
		t.Fatal("Encoded string should not match original text")
	}

	// Decode the text
	decoded, err := Decode(encoded, key)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if decoded != originalText {
		t.Fatalf("Decoded text does not match original. Got: %s, Want: %s", decoded, originalText)
	}
}

func TestEncodeDecodeWithEnvironmentKey(t *testing.T) {
	// Save original env var
	originalEnv := os.Getenv("QAI_SECRET_KEY")
	defer os.Setenv("QAI_SECRET_KEY", originalEnv)

	// Set environment variable for test
	os.Setenv("QAI_SECRET_KEY", "env_secret_key")

	// Test data
	originalText := "Hello World"

	// Encode the text
	encoded, err := Encode(originalText)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	if encoded == "" {
		t.Fatal("Encoded string is empty")
	}
	if encoded == originalText {
		t.Fatal("Encoded string should not match original text")
	}

	// Decode the text
	decoded, err := Decode(encoded)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if decoded != originalText {
		t.Fatalf("Decoded text does not match original. Got: %s, Want: %s", decoded, originalText)
	}
}

func TestEncodeDecodeWithDifferentKeys(t *testing.T) {
	// Test data
	originalText := "Hello World"
	key1 := "secret_key_one"
	key2 := "secret_key_two"

	// Encode with key1
	encoded, err := Encode(originalText, key1)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Try to decode with wrong key
	decoded, err := Decode(encoded, key2)
	if err == nil {
		t.Fatal("Expected error when decoding with wrong key, got nil")
	}
	if err != nil && !contains(err.Error(), "cipher: message authentication failed") {
		t.Fatalf("Expected GCM error, got: %v", err)
	}
	if decoded == originalText {
		t.Fatal("Decoded text should not match original when using wrong key")
	}

	// Decode with correct key
	decoded, err = Decode(encoded, key1)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if decoded != originalText {
		t.Fatalf("Decoded text does not match original. Got: %s, Want: %s", decoded, originalText)
	}
}

func TestLongKeyTruncation(t *testing.T) {
	// Test data
	originalText := "Hello World"
	longKey := "this_is_a_very_long_key_that_will_be_truncated_to_32_bytes_for_aes256"

	// Encode with long key
	encoded, err := Encode(originalText, longKey)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}

	// Decode with same long key
	decoded, err := Decode(encoded, longKey)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if decoded != originalText {
		t.Fatalf("Decoded text does not match original. Got: %s, Want: %s", decoded, originalText)
	}
}

func TestEmptyString(t *testing.T) {
	// Test data
	emptyString := ""
	key := "test_secret_key"

	// Encode empty string
	encoded, err := Encode(emptyString, key)
	if err != nil {
		t.Fatalf("Encode failed: %v", err)
	}
	if encoded == "" {
		t.Fatal("Encoded string is empty")
	}

	// Decode back to empty string
	decoded, err := Decode(encoded, key)
	if err != nil {
		t.Fatalf("Decode failed: %v", err)
	}
	if decoded != emptyString {
		t.Fatalf("Decoded text does not match original. Got: %s, Want: %s", decoded, emptyString)
	}
}

// contains checks if substr is in str
func contains(str, substr string) bool {
	return len(substr) == 0 || (len(str) >= len(substr) && (str == substr || (len(str) > len(substr) && (str[0:len(substr)] == substr || contains(str[1:], substr)))))
}
