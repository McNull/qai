package utils

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// LoadEnvFile loads environment variables from .env file in the project root
// It attempts to find the .env file by walking up from the current directory
// until it finds the file or reaches the root directory
func LoadEnvFile() error {
	// Check if environment loading should be skipped
	if os.Getenv("QAI_NO_ENV_LOAD") == "1" {
		return nil
	}

	// Start with the current directory
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Walk up directory tree until we find .env or reach root
	for {
		// Check if .env exists in the current directory
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return godotenv.Load(envPath)
		}

		// Go up one directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// We've reached the root, stop searching
			break
		}
		dir = parent
	}

	// No .env file found, but this isn't an error condition
	return nil
}
