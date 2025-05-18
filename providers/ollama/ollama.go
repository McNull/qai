package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/mcnull/qai/shared/provider"
)

type OllamaProvider struct {
	provider.ProviderBase
	config Config
}

func NewOllamaProvider(config provider.IConfig, appCtx *provider.AppContext) (provider.IProvider, error) {

	if config == nil {
		return nil, fmt.Errorf("config is nil")
	}

	cfg, ok := config.(*Config)
	if !ok {
		return nil, fmt.Errorf("type mismatch: expected *Config, got %T", config)
	}

	p := &OllamaProvider{
		config: *cfg,
	}

	p.ProviderBase = *provider.NewProviderBase("ollama", appCtx)

	return p, nil
}

func (p *OllamaProvider) Generate(ctx context.Context, request provider.GenerateRequest) (<-chan provider.GenerateResponse, <-chan error) {
	responseChan := make(chan provider.GenerateResponse)
	errorChan := make(chan error, 1)

	go func() {
		defer close(responseChan)
		defer close(errorChan)

		// Convert to Ollama request
		ollamaReq := GenerateRequest{
			Model:  p.config.Model,
			Prompt: request.Prompt,
			System: request.System,
			Stream: !p.Flags().Color,
			Options: &Options{
				Seed: p.config.Seed,
			},
		}

		// Set up custom HTTP client to get raw response instead of using the library's scanner
		client := &http.Client{}
		jsonData, err := json.Marshal(ollamaReq)
		if err != nil {
			errorChan <- fmt.Errorf("error marshaling request: %w", err)
			return
		}

		// Create request manually
		req, err := http.NewRequestWithContext(ctx, "POST",
			p.config.URL+"/api/generate", bytes.NewBuffer(jsonData))
		if err != nil {
			errorChan <- fmt.Errorf("error creating request: %w", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		// Send request
		resp, err := client.Do(req)
		if err != nil {
			errorChan <- fmt.Errorf("error sending request: %w", err)
			return
		}
		defer resp.Body.Close()

		// Check for non-200 status code
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			errorChan <- fmt.Errorf("error response from server: %s", body)
			return
		}

		decoder := json.NewDecoder(resp.Body)

		for {
			var rawMessage json.RawMessage
			if err := decoder.Decode(&rawMessage); err != nil {
				if err == io.EOF {
					return
				}
				errorChan <- fmt.Errorf("error decoding response: %w", err)
				return
			}

			var ollamaResp GenerateResponse
			if err := json.Unmarshal(rawMessage, &ollamaResp); err != nil {
				errorChan <- fmt.Errorf("error unmarshaling response: %w", err)
				return
			}

			providerResp := provider.GenerateResponse{
				Raw:      rawMessage,
				Response: ollamaResp.Response,
				Done:     ollamaResp.Done,
			}

			select {
			case responseChan <- providerResp:
			case <-ctx.Done():
				errorChan <- ctx.Err()
				return
			}

			if ollamaResp.Done {
				return
			}
		}
	}()

	return responseChan, errorChan
}
