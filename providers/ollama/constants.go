package ollama

const (
	DEFAULT_MODEL    = "llama3.2"
	DEFAULT_URL      = "http://127.0.0.1:11434"
	DEFAULT_SETTINGS = `{
		"model": "` + DEFAULT_MODEL + `",
		"url": "` + DEFAULT_URL + `",
	}`
)
