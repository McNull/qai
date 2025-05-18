package provider

type AppContext struct {
	Flags        *FlagValues
	Provider     IProvider
	SystemPrompt string
}
