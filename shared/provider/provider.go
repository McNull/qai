package provider

import "context"

type IProvider interface {
	GetName() string
	Init() error
	Generate(ctx context.Context, request GenerateRequest) (<-chan GenerateResponse, <-chan error)
}

type ProviderBase struct {
	Name       string
	AppContext *AppContext
}

func NewProviderBase(name string, appCtx *AppContext) *ProviderBase {
	p := &ProviderBase{
		Name:       name,
		AppContext: appCtx,
	}

	return p
}

func (p *ProviderBase) GetName() string {
	return p.Name
}

func (p *ProviderBase) Init() error {
	return nil
}

func (p *ProviderBase) Flags() *FlagValues {
	return p.AppContext.Flags
}
