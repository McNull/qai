package app

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"os"
	"strings"
	"time"

	"github.com/mcnull/qai/providers/github"
	"github.com/mcnull/qai/providers/ollama"
	"github.com/mcnull/qai/shared/markdown"
	"github.com/mcnull/qai/shared/platform"
	"github.com/mcnull/qai/shared/provider"
	"github.com/mcnull/qai/shared/throbber"
	"github.com/mcnull/qai/shared/utils"
)

type App struct {
	provider.AppContext
	Config *Config
}

func NewApp() *App {
	return &App{
		AppContext: provider.AppContext{
			Flags: provider.NewFlagValues(
				DEFAULT_CONFIG_FILEPATH,
				DEFAULT_SYSTEM_PROMPT,
			),
			Provider: nil,
		},
		Config: nil, // Config will be initialized later
	}
}

func (app *App) parseArgs(args []string) error {

	// Parse arguments

	flags, err := parseFlags(args[1:], true, app.Flags)

	if err != nil {
		err = fmt.Errorf("error parsing flags: %w", err)
		return err
	}

	app.Flags = flags

	return nil
}

func (app *App) initConfig() (bool, error) {

	flags := app.Flags

	// Check if we need to create a new config file

	if flags.CreateConfig {
		_, err := createNewConfigFile(flags.ConfigFile)

		if err != nil {
			return false, err
		}

		return false, nil
	}

	// Load config
	config, err := LoadConfig(flags.ConfigFile)

	if err != nil {

		// If the config file does not exist and the config file is set to the default:
		if os.IsNotExist(err) && flags.ConfigFile == DEFAULT_CONFIG_FILEPATH {
			// create a new config file
			fmt.Println("Config file does not exist, creating a new one...")
			config, err = createNewConfigFile(flags.ConfigFile)

			if err != nil {
				return false, err
			}
		} else {
			err = fmt.Errorf("error loading config: %w", err)
			return false, err
		}
	}

	app.Config = config

	// Check if we need to login to GitHub
	if flags.GithubLogin {
		token, err := github.Login(flags.Debug)

		if err != nil {
			err = fmt.Errorf("error logging in to GitHub: %w", err)
			return false, err
		}

		// Store the token in the config
		ghConfig, ok := app.Config.Providers.GitHub.(*github.Config)

		if !ok {
			err = fmt.Errorf("error casting config to GitHub config")
			return false, err
		}

		ghConfig.Token = token

		err = app.Config.Save(flags.ConfigFile)

		if err != nil {
			err = fmt.Errorf("error saving config: %w", err)
			return false, err
		}
	}

	// Ensure we have a profile name
	if flags.Profile == "" {
		flags.Profile = app.Config.Profile
	}

	return true, nil
}

func (app *App) initProvider() error {

	// Get the profile from the config

	profile, err := app.Config.GetProfile(app.Flags.Profile)

	if err != nil {
		err = fmt.Errorf("error getting profile \"%s\" from config: %w", app.Flags.Profile, err)
		return err
	}

	if app.Flags.Debug {
		fmt.Println("Profile:")
		utils.Dump(profile)
	}

	var pConfig provider.IConfig              // default from config
	var pConfigFactory provider.ConfigFactory // config factory
	var pFactory provider.ProviderFactory     // provider factory

	switch profile.Provider {

	case "ollama":
		pConfig = app.Config.Providers.Ollama
		pConfigFactory = ollama.NewConfig
		pFactory = ollama.NewOllamaProvider
		break

	case "github":
		pConfig = app.Config.Providers.GitHub
		pConfigFactory = github.NewConfig
		pFactory = github.NewGitHubProvider
		break
	}

	pConfig, err = provider.InitConfig(
		pConfig,
		pConfigFactory,
		profile.Settings,
	)

	if err != nil {
		err = fmt.Errorf("error initializing provider config: %w", err)
		return err
	}

	if app.Flags.Debug {
		fmt.Println("Provider config:")
		utils.Dump(pConfig)
	}

	var p provider.IProvider
	p, err = pFactory(pConfig, &app.AppContext)

	if err != nil {
		err = fmt.Errorf("error creating provider: %w", err)
	}

	err = p.Init()

	if err != nil {
		err = fmt.Errorf("error initializing provider: %w", err)
		return err
	}

	app.Provider = p

	return nil
}

func (app *App) Init(args []string) (bool, error) {

	c := true

	// Parse arguments

	err := app.parseArgs(args)

	if err != nil {
		return false, err
	}

	utils.DumpInColor = app.Flags.Color

	if app.Flags.Version {
		fmt.Printf("%s %s (%s)\n", APP_NAME, APP_VERSION, "https://github.com/mcnull/qai")
		fmt.Printf("Config file: %s\n", app.Flags.ConfigFile)
		return false, nil
	}

	// Load config

	c, err = app.initConfig()

	if err != nil {
		return false, err
	}

	if c != true {
		return false, nil
	}

	// Initialize provider

	err = app.initProvider()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (app *App) getSystemPrompt() (string, error) {

	system := app.Flags.System
	if system == "" {
		system = app.SystemPrompt
	}

	platformInfo, err := platform.GetInfo()

	if err != nil {
		err = fmt.Errorf("error getting platform info: %w", err)
		return "", err
	}

	verbose := "brief and concise. Don't give explanations or details."

	if app.Flags.Verbose {
		verbose = "verbose and detailed. Explain everything and provide examples."
	}

	m := map[string]string{
		"Platform": platformInfo,
		"Verbose":  verbose,
	}

	tmpl, err := template.New("system").Parse(system)
	if err != nil {
		return "", fmt.Errorf("error parsing system prompt template: %w", err)
	}

	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, m)
	if err != nil {
		return "", fmt.Errorf("error executing system prompt template: %w", err)
	}

	return tpl.String(), nil
}

func (app *App) Run() error {

	// Check if prompt is empty
	if app.Flags.Prompt == "" {
		fmt.Println("No prompt provided")
		fmt.Println("Use -h or --help for more information")
		return nil
	}

	system, err := app.getSystemPrompt()
	if err != nil {
		err = fmt.Errorf("error getting system prompt: %w", err)
	}

	request := &provider.GenerateRequest{
		System: system,
		Prompt: app.Flags.Prompt,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	throbber := throbber.NewThrobber().
		WithMessage("Generating response...").
		WithThrob(throbber.ThrobByName("binary")).
		Start()

	mdRenderer, err := markdown.NewMarkdownRenderer("dark")
	if err != nil {
		return fmt.Errorf("failed to initialize markdown renderer: %w", err)
	}

	// Don't defer stop - we'll stop it explicitly to ensure proper sequence

	responseChan, errorChan := app.Provider.Generate(ctx, *request)

	for {
		select {
		case response, ok := <-responseChan:
			if !ok {
				if throbber.IsRunning() {
					throbber.Stop()
				}
				return nil // Channel closed
			}

			if throbber.IsRunning() {
				throbber.Stop()
			}

			if app.Flags.DebugStream {
				utils.Dump(response)
			} else {

				rendered := response.Response

				if app.Flags.Color {
					rendered, err = mdRenderer.Render(response.Response, response.Done)
					if err != nil {
						return fmt.Errorf("error rendering markdown: %w", err)
					}
				}

				if rendered != "" {
					fmt.Print(rendered)
				}

				if response.Done && app.Flags.Color {
					// Flush any remaining content
					remaining, _ := mdRenderer.Render("", true)
					remaining = strings.Trim(remaining, "\n")
					if remaining != "" {
						fmt.Print(remaining)
					}
				}
			}

			if response.Done {
				// Markdown adds a trailing newline
				if !app.Flags.Color {
					fmt.Println()
				}

				return nil
			}

		case err, ok := <-errorChan:
			if !ok {
				if throbber.IsRunning() {
					throbber.Stop()
				}
				return nil // Channel closed
			}

			if throbber.IsRunning() {
				throbber.Stop()
			}

			return err // Return the error

		case <-ctx.Done():
			if throbber.IsRunning() {
				throbber.Stop()
			}
			fmt.Println("Operation timed out")
			return ctx.Err()
		}
	}

}
