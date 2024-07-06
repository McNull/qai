package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/pflag"
)

// parseArgs parses the command line arguments and returns an Args object
func parseArgs() *Options {
	// Custom usage function
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		pflag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nDefault values can be set in the config file at ~/.config/qai/config.json")
	}

	defaultOptions, err := getOptions()

	if err != nil {
		fmt.Fprintln(os.Stderr, "Error getting default arguments:", err)
		os.Exit(1) // Exit code 1 for general error
	}

	// Define flags
	model := pflag.StringP("model", "m", defaultOptions.Model, "Model to use for generation")
	url := pflag.StringP("url", "u", defaultOptions.Url, "URL of the ollama service")
	system := pflag.StringP("system", "s", "", "System message to use for generation")
	platform := pflag.StringP("os", "o", defaultOptions.Platform, "Operating system")
	version := pflag.BoolP("version", "v", false, "Print the version and exit")

	// Parse the flags
	pflag.Parse()

	// If the version flag is set, print the version and exit
	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	// After parsing, if system is not set by the user, use the default
	if *system == "" {
		*system = defaultOptions.System
	}

	// Combine the remaining arguments into the Prompt
	prompt := strings.Join(pflag.Args(), " ")

	if prompt == "" {
		fmt.Fprintln(os.Stderr, "Error: no prompt provided")
		fmt.Println()
		pflag.Usage()
		os.Exit(2) // Exit code 2 for command line usage error
	}

	// Return the Args object
	return &Options{
		Model:    *model,
		Prompt:   prompt,
		Url:      *url,
		System:   *system,
		Platform: *platform,
	}
}
