package report

import (
	"fmt"

	"github.com/osteele/liquid"
)

// Renderer handles Liquid template rendering.
type Renderer struct {
	engine *liquid.Engine
}

// NewRenderer creates a new Renderer with a configured Liquid engine.
func NewRenderer() *Renderer {
	engine := liquid.NewEngine()
	return &Renderer{engine: engine}
}

// RenderTemplate renders a Liquid template string with the provided data bindings.
func (r *Renderer) RenderTemplate(templateContent string, data map[string]any) (string, error) {
	tmpl, err := r.engine.ParseString(templateContent)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	out, err := tmpl.RenderString(data)
	if err != nil {
		return "", fmt.Errorf("failed to render template: %w", err)
	}

	return out, nil
}

// RenderBytes is like RenderTemplate but returns bytes.
func (r *Renderer) RenderBytes(templateContent []byte, data map[string]any) ([]byte, error) {
	out, err := r.engine.ParseAndRender(templateContent, data)
	if err != nil {
		return nil, fmt.Errorf("failed to render template: %w", err)
	}
	return out, nil
}
