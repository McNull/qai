package github

import (
	"fmt"
	"github.com/mcnull/qai/shared/provider"
)

type Config struct {
	Model string `json:"model"`
	Token string `json:"token"`
}

func NewConfig() provider.IConfig {
	return &Config{
		Model: DEFAULT_MODEL,
		Token: "",
	}
}

func (c *Config) Merge(other provider.IConfig) error {

	if other == nil {
		return fmt.Errorf("other config is nil")
	}

	otherConfig, ok := other.(*Config)
	if !ok {
		return fmt.Errorf("type mismatch: expected *Config, got %T", other)
	}

	if otherConfig.Model != DEFAULT_MODEL {
		c.Model = otherConfig.Model
	}

	if otherConfig.Token != "" {
		c.Token = otherConfig.Token
	}

	return nil
}
