package ollama

type GenerateRequest struct {
	Model     string   `json:"model"`
	Prompt    string   `json:"prompt,omitempty"`
	System    string   `json:"system,omitempty"`
	Template  string   `json:"template,omitempty"`
	Context   []int    `json:"context,omitempty"`
	Stream    bool     `json:"stream"`
	Raw       bool     `json:"raw,omitempty"`
	Format    string   `json:"format,omitempty"`
	Options   *Options `json:"options,omitempty"`
	KeepAlive string   `json:"keep_alive,omitempty"`
}

type Options struct {
	NumKeep          *int     `json:"num_keep,omitempty"`
	Seed             *int     `json:"seed,omitempty"`
	NumPredict       *int     `json:"num_predict,omitempty"`
	TopK             *int     `json:"top_k,omitempty"`
	TopP             *float64 `json:"top_p,omitempty"`
	TFSZ             *float64 `json:"tfs_z,omitempty"`
	TypicalP         *float64 `json:"typical_p,omitempty"`
	RepeatLastN      *int     `json:"repeat_last_n,omitempty"`
	Temperature      *float64 `json:"temperature,omitempty"`
	RepeatPenalty    *float64 `json:"repeat_penalty,omitempty"`
	PresencePenalty  *float64 `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float64 `json:"frequency_penalty,omitempty"`
	Mirostat         *int     `json:"mirostat,omitempty"`
	MirostatTau      *float64 `json:"mirostat_tau,omitempty"`
	MirostatEta      *float64 `json:"mirostat_eta,omitempty"`
	PenalizeNewline  *bool    `json:"penalize_newline,omitempty"`
	Stop             []string `json:"stop,omitempty"`
	NumGPU           *int     `json:"num_gpu,omitempty"`
	NumThread        *int     `json:"num_thread,omitempty"`
	NumCtx           *int     `json:"num_ctx,omitempty"`
	LogitsAll        *bool    `json:"logits_all,omitempty"`
	EmbeddingOnly    *bool    `json:"embedding_only,omitempty"`
	F16KV            *bool    `json:"f16_kv,omitempty"`
}

type GenerateResponse struct {
	Model           string `json:"model,omitempty"`
	CreatedAt       string `json:"created_at,omitempty"`
	Response        string `json:"response"`
	Done            bool   `json:"done,omitempty"`
	Context         []int  `json:"context,omitempty"`
	TotalDuration   int64  `json:"total_duration,omitempty"`
	LoadDuration    int64  `json:"load_duration,omitempty"`
	PromptEvalCount int    `json:"prompt_eval_count,omitempty"`
	EvalCount       int    `json:"eval_count,omitempty"`
	EvalDuration    int64  `json:"eval_duration,omitempty"`
}
