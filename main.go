package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/mcnull/qai/api"
)

func formatPrompt(prompt string, platform string) string {
	template := `Using the terminal on operating system: {{.platform}}:\n{{.prompt}}`
	result := strings.ReplaceAll(template, "{{.platform}}", platform)
	result = strings.ReplaceAll(result, "{{.prompt}}", prompt)
	return result
}

func main() {

	args := parseArgs()

	// Parse the URL string to a url.URL object
	baseURL, err := url.Parse(args.Url)

	if err != nil {
		panic("Failed to parse URL: " + err.Error())
	}

	// Create an http.Client instance
	httpClient := &http.Client{}

	// Create an API client
	client := api.NewClient(baseURL, httpClient)

	// Create a context
	ctx := context.Background()

	// Set up the GenerateRequest
	req := &api.GenerateRequest{
		Model:  args.Model,
		Prompt: formatPrompt(args.Prompt, args.Platform),
		System: args.System,
	}

	// Define the GenerateResponseFunc
	responseFunc, err := getResponseRenderer()

	if err != nil {
		fmt.Println("Error getting response renderer:", err)
		return
	}

	// Call the Generate function
	err = client.Generate(ctx, req, responseFunc)
	if err != nil {
		fmt.Println("Error calling Generate:", err)
		return
	}

	fmt.Println()

}
