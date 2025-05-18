package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"github.com/mcnull/qai/providers/github"
	"github.com/mcnull/qai/providers/ollama"
	"github.com/mcnull/qai/shared/jsonmap"
	"github.com/mcnull/qai/shared/provider"
)

type Config struct {
	Profile   string             `json:"profile"`
	System    string             `json:"system"`
	Providers ProvidersConfig    `json:"providers"`
	Profiles  map[string]Profile `json:"profiles"`
}

type ProvidersConfig struct {
	Ollama provider.IConfig `json:"ollama"`
	GitHub provider.IConfig `json:"github"`
}

type Profile struct {
	Provider string          `json:"provider"`
	Settings jsonmap.JsonMap `json:"settings,omitempty"`
}

func NewConfig() *Config {

	ollamaConfig := ollama.NewConfig()
	githubConfig := github.NewConfig()

	return &Config{
		Profile: DEFAULT_PROFILE,
		System:  DEFAULT_SYSTEM_PROMPT,
		Providers: ProvidersConfig{
			Ollama: ollamaConfig,
			GitHub: githubConfig,
		},
		Profiles: map[string]Profile{
			DEFAULT_PROFILE: {
				Provider: "ollama",
			},
		},
	}
}

func LoadConfig(filepath string) (*Config, error) {

	// unmarshal the JSON file into a Config struct

	config := NewConfig()

	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) GetProfile(name string) (*Profile, error) {

	if name == "" {
		return nil, fmt.Errorf("profile name is empty")
	}

	profile, ok := c.Profiles[name]
	if !ok {
		return nil, os.ErrNotExist
	}
	return &profile, nil
}

func (c *Config) Save(fp string) error {

	// Make absolute path
	absPath, err := filepath.Abs(fp)
	if err != nil {
		return err
	}

	fp = absPath

	// ensure the directory exists
	dir := path.Dir(fp)

	// create the directory if it doesn't exist
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	file, err := os.Create(fp)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(c)
	if err != nil {
		return err
	}

	return nil
}

func createNewConfigFile(fp string) (*Config, error) {
	// Create a new config file with default values
	config := NewConfig()
	err := config.Save(fp)
	if err != nil {
		err = fmt.Errorf("error creating new config file: %w", err)
		return nil, err
	}

	fmt.Printf("Created new config file at %s\n", fp)

	return config, nil
}
