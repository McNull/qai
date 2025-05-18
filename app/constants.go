package app

import (
	"os"
	"path/filepath"
)

var DEFAULT_CONFIG_FILEPATH string

const (
	APP_NAME              = "qai"
	APP_VERSION           = "0.5.0"
	DEFAULT_PROFILE       = "default"
	DEFAULT_SYSTEM_PROMPT = "The user is running a terminal in the following environment: {{.Platform}}.\nYour responses are {{.Verbose}}."
)

func init() {
	// Resolve "$HOME/.config/github.com/mcnull/qai" to the config directory (windows/linux/mac)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	dir := filepath.Join(homeDir, ".config", "qai")
	DEFAULT_CONFIG_FILEPATH = filepath.Join(dir, "config.json")
}
