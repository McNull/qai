package provider

import (
	"fmt"
	"github.com/mcnull/qai/shared/jsonmap"
)

type IProviderRegistry interface {
	Register(name string, provider IProvider) error
	Get(name string) IProvider
	List() []string
	GetAll() []IProvider
	GetProviderDefaults() (map[string]jsonmap.JsonMap, error)
}

type ProviderRegistry struct {
	providers map[string]IProvider
}

func NewRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		providers: make(map[string]IProvider),
	}
}

func (r *ProviderRegistry) Register(name string, provider IProvider) error {
	if _, exists := r.providers[name]; exists {
		return fmt.Errorf("provider already registered: %s", name)
	}
	r.providers[name] = provider

	return nil
}

func (r *ProviderRegistry) Get(name string) IProvider {
	if provider, exists := r.providers[name]; exists {
		return provider
	}
	panic("Provider not found: " + name)
}

func (r *ProviderRegistry) List() []string {
	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	return names
}

func (r *ProviderRegistry) GetAll() []IProvider {
	providers := make([]IProvider, 0, len(r.providers))
	for _, provider := range r.providers {
		providers = append(providers, provider)
	}
	return providers
}
