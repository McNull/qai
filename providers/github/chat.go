package github

import "encoding/json"

const GITHUB_CHAT_URL = "https://api.githubcopilot.com/chat/completions"

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewChatMessage(role, content string) *ChatMessage {
	return &ChatMessage{
		Role:    role,
		Content: content,
	}
}

type ChatRequest struct {
	Model       string         `json:"model"`
	Temperature float32        `json:"temperature"`
	Top_p       float32        `json:"top_p"`
	N           int            `json:"n"`
	Stream      bool           `json:"stream"`
	Messages    []*ChatMessage `json:"messages"`
}

func NewChatRequest(model string, messages []*ChatMessage) *ChatRequest {
	return &ChatRequest{
		Model:       model,
		Temperature: 0.5,
		Top_p:       1.0,
		N:           1,
		Stream:      false,
		Messages:    messages,
	}
}

func (c *ChatRequest) ToJson() ([]byte, error) {
	return json.Marshal(c)
}
