package provider

import (
	"flag"
)

type FlagValues struct {
	ConfigFile   string
	CreateConfig bool
	Profile      string
	Prompt       string
	Debug        bool
	DebugStream  bool
	System       string
	Verbose      bool
	Color        bool
	GithubLogin  bool
	Version      bool
}

func NewFlagValues(configFile, system string) *FlagValues {
	return &FlagValues{
		ConfigFile:   configFile,
		CreateConfig: false,
		Profile:      "",
		Prompt:       "",
		Debug:        false,
		DebugStream:  false,
		System:       system,
		Verbose:      false,
		Color:        true,
		GithubLogin:  false,
		Version:      false,
	}
}

func CreateFlagSet(name string, v *FlagValues, exitOnError bool) *flag.FlagSet {

	exitRule := flag.ExitOnError

	if !exitOnError {
		exitRule = flag.ContinueOnError
	}

	fs := flag.NewFlagSet(name, exitRule)

	fs.StringVar(&v.ConfigFile, "config", v.ConfigFile, "Path to the config file")
	fs.BoolVar(&v.CreateConfig, "create-config", v.CreateConfig, "Create a new config file with default values")
	fs.StringVar(&v.Profile, "profile", v.Profile, "Profile name")
	fs.StringVar(&v.System, "system", v.System, "System prompt")
	fs.BoolVar(&v.Debug, "debug", v.Debug, "Enable debug mode")
	fs.BoolVar(&v.DebugStream, "debug-stream", v.DebugStream, "Enable debug response stream")
	fs.BoolVar(&v.Color, "color", v.Color, "Enable colored output")
	fs.BoolVar(&v.Verbose, "verbose", v.Verbose, "Enable verbose output")
	fs.BoolVar(&v.GithubLogin, "github-login", v.GithubLogin, "Create a new GitHub auth token")
	fs.BoolVar(&v.Version, "version", v.Version, "Show version information")

	return fs
}
