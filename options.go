package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Options struct to hold the command line arguments
type Options struct {
	Model    string `json:"model"`
	Prompt   string `json:"-"`
	Url      string `json:"url"`
	System   string `json:"system"`
	Platform string `json:"-"`
}

func getDefaultOptions() (*Options, error) {

	platform, err := getPlatform()

	if err != nil {
		return nil, err
	}

	return &Options{
		Model:  "llama3",
		Prompt: "",
		Url:    "http://localhost:11434",
		System: strings.TrimSpace(`
		Your answers are always formal, short and to the point.
		Your answers never contain explanations or examples unless explicitly asked for.
		`),
		Platform: platform,
	}, nil
}

func getOptions() (*Options, error) {
	defaults, err := getDefaultOptions()
	if err != nil {
		return nil, err
	}

	configPath, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	configPath = filepath.Join(configPath, ".config", "qai", "config.json")

	// Check if the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create the directory if it does not exist
		if err := os.MkdirAll(filepath.Dir(configPath), 0755); err != nil {
			return nil, err
		}

		// Create the file with default options
		file, err := os.Create(configPath)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		encoder := json.NewEncoder(file)
		encoder.SetIndent("", "    ") // Sets the prefix and indentation levels
		if err := encoder.Encode(defaults); err != nil {
			return nil, err
		}
	}

	// File exists, load and override defaults
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(defaults); err != nil {
		return nil, err
	}

	return defaults, nil
}
