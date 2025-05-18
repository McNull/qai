package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"github.com/mcnull/qai/shared/jsonmap"
	"github.com/mcnull/qai/shared/provider"
)

type GitHubProvider struct {
	provider.ProviderBase
	config Config
}

func NewGitHubProvider(config provider.IConfig, appCtx *provider.AppContext) (provider.IProvider, error) {

	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}

	cfg, ok := config.(*Config)
	if !ok {
		return nil, fmt.Errorf("type mismatch: expected *Config, got %T", config)
	}

	p := &GitHubProvider{
		config: *cfg,
	}

	p.ProviderBase = *provider.NewProviderBase("github", appCtx)

	return p, nil
}

func (p *GitHubProvider) Init() error {
	if p.config.Token == "" {
		return fmt.Errorf("missing github token.\n\nUse --github-login to create a new token.")
	}

	return nil
}

func (p *GitHubProvider) Generate(ctx context.Context, request provider.GenerateRequest) (<-chan provider.GenerateResponse, <-chan error) {

	responseChan := make(chan provider.GenerateResponse)
	errorChan := make(chan error, 1)

	go func() {
		defer close(responseChan)
		defer close(errorChan)

		apiToken, err := requestApiToken(p.config.Token)

		if err != nil {
			errorChan <- fmt.Errorf("failed to request API token: %w", err)
			return
		}

		chatMessages := []*ChatMessage{
			NewChatMessage("system", request.System),
			NewChatMessage("user", request.Prompt),
		}

		chatReq := NewChatRequest(p.config.Model, chatMessages)

		chatReq.Stream = !p.Flags().Color

		jsonBody, err := chatReq.ToJson()

		if err != nil {
			errorChan <- fmt.Errorf("error marshaling request: %w", err)
			return
		}

		/*
			POST https://api.githubcopilot.com/chat/completions
			User-Agent: github.com/mcnull/qai
			Authorization: Bearer {{api_token}}
			Editor-Version: github.com/mcnull/qai/0.1.0
			Content-Type: application/json
			Copilot-Integration-Id: vscode-chat

			{
				"model": "gpt-4",
				"temperature": 0.5,
				"top_p": 1.0,
				"n": 1,
				"stream": false,
				"messages": [
					{
						"role": "system",
						"content": "You are a helpful assistant."
					},
					{
						"role": "user",
						"content": "How many fingers am I holding up?"
					}
				]
			}
		*/

		client := &http.Client{}
		req, err := http.NewRequestWithContext(ctx, "POST", GITHUB_CHAT_URL, bytes.NewBuffer(jsonBody))

		if err != nil {
			errorChan <- fmt.Errorf("error creating request: %w", err)
			return
		}

		req.Header.Set("User-Agent", "github.com/mcnull/qai")
		req.Header.Set("Authorization", "Bearer "+apiToken)
		req.Header.Set("Editor-Version", "github.com/mcnull/qai/0.1.0")
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Copilot-Integration-Id", "vscode-chat")
		req.Header.Set("Accept", "application/json")

		client.Do(req)
		resp, err := client.Do(req)
		if err != nil {
			errorChan <- fmt.Errorf("error sending request: %w", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			errorChan <- fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			return
		}

		decoder := json.NewDecoder(resp.Body)

		for {
			var response jsonmap.JsonMap

			if err := decoder.Decode(&response); err != nil {
				if err == io.EOF {
					return
				}
				errorChan <- fmt.Errorf("error decoding response: %w", err)
				return
			}

			choices, ok := response["choices"].([]any)

			if !ok || len(choices) == 0 {
				errorChan <- fmt.Errorf("invalid response format (missing choices): %v", response)
				return
			}

			choice, ok := choices[0].(map[string]any)
			if !ok {
				errorChan <- fmt.Errorf("invalid choice format: %v", choices[0])
				return
			}

			messageObj, ok := choice["message"].(map[string]any)
			if !ok {
				errorChan <- fmt.Errorf("invalid response format (missing message object): %v", choice)
				return
			}

			content, ok := messageObj["content"].(string)
			if !ok {
				errorChan <- fmt.Errorf("invalid response format (missing content in message): %v", messageObj)
				return
			}

			providerResp := provider.GenerateResponse{
				Raw:      response,
				Response: content,
				Done:     true,
			}

			select {
			case responseChan <- providerResp:
			case <-ctx.Done():
				errorChan <- ctx.Err()
				return
			}

			return
		}

	}()

	return responseChan, errorChan
}
