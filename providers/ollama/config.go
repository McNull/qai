package ollama

import (
	"fmt"
	"github.com/mcnull/qai/shared/provider"
)

type Config struct {
	Model string `json:"model"`
	URL   string `json:"url"`
	Seed  *int   `json:"seed,omitempty"`
}

func NewConfig() provider.IConfig {
	return &Config{
		Model: DEFAULT_MODEL,
		URL:   DEFAULT_URL,
	}
}

func (c *Config) Merge(other provider.IConfig) error {
	o, ok := other.(*Config)
	if !ok {
		return fmt.Errorf("type mismatch: expected *Config, got %T", other)
	}

	if o.Model != DEFAULT_MODEL {
		c.Model = o.Model
	}

	if o.URL != DEFAULT_URL {
		c.URL = o.URL
	}

	if o.Seed != nil {
		c.Seed = o.Seed
	}

	return nil
}
