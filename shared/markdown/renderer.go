package markdown

import (
	"strings"

	"github.com/charmbracelet/glamour"
)

// MarkdownRenderer handles streaming markdown content
type MarkdownRenderer struct {
	renderer *glamour.TermRenderer
	buffer   strings.Builder
	style    string
}

// NewMarkdownRenderer creates a new markdown renderer
func NewMarkdownRenderer(style string) (*MarkdownRenderer, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithStylePath(style),
		glamour.WithWordWrap(80),
	)
	if err != nil {
		return nil, err
	}

	return &MarkdownRenderer{
		renderer: r,
		buffer:   strings.Builder{},
		style:    style,
	}, nil
}

// Render processes a chunk of markdown text
func (m *MarkdownRenderer) Render(chunk string, isComplete bool) (string, error) {
	// Append the new chunk to the buffer
	m.buffer.WriteString(chunk)

	// If this isn't the final chunk, don't render yet
	if !isComplete {
		return "", nil
	}

	// Get the entire content from the buffer
	content := m.buffer.String()

	// Render the markdown content
	rendered, err := m.renderer.Render(content)
	if err != nil {
		return "", err
	}

	// Clear the buffer after rendering
	m.buffer.Reset()

	return rendered, nil
}
