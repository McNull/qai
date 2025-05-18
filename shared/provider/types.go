package provider

type GenerateRequest struct {
	System string
	Prompt string
}

type GenerateResponse struct {
	Raw      any
	Response string `json:"response"`
	Done     bool   `json:"done,omitempty"`
}

type ConfigFactory func() IConfig
type ProviderFactory func(config IConfig, appCtx *AppContext) (IProvider, error)
