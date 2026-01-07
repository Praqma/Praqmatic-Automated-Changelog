// Package report handles changelog report generation.
package report

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/config"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/logging"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/model"
)

// Generator coordinates report generation from tasks and commits.
type Generator struct {
	tasks    *model.PACTaskCollection
	commits  *model.PACCommitCollection
	renderer *Renderer
}

// NewGenerator creates a new Generator.
func NewGenerator(tasks *model.PACTaskCollection, commits *model.PACCommitCollection) *Generator {
	return &Generator{
		tasks:    tasks,
		commits:  commits,
		renderer: NewRenderer(),
	}
}

// Generate processes all configured templates and writes output files.
func (g *Generator) Generate(settings *config.Settings) error {
	for _, tmplCfg := range settings.Templates {
		if err := g.generateTemplate(tmplCfg, settings); err != nil {
			logging.Warn("Failed to generate template %s: %v", tmplCfg.Location, err)
			// Continue with other templates
		}
	}
	return nil
}

// generateTemplate processes a single template configuration.
func (g *Generator) generateTemplate(tmplCfg config.TemplateConfig, settings *config.Settings) error {
	// Read template file
	templateContent, err := os.ReadFile(tmplCfg.Location)
	if err != nil {
		return fmt.Errorf("failed to read template file %s: %w", tmplCfg.Location, err)
	}

	logging.Info("Processing template: %s", tmplCfg.Location)

	// Prepare Liquid data bindings
	liquidData := g.prepareLiquidData(settings)

	// Render template
	output, err := g.renderer.RenderBytes(templateContent, liquidData)
	if err != nil {
		return fmt.Errorf("failed to render template: %w", err)
	}

	// Write output
	if tmplCfg.Output != "" {
		// Ensure output directory exists
		outputDir := filepath.Dir(tmplCfg.Output)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}

		if err := os.WriteFile(tmplCfg.Output, output, 0644); err != nil {
			return fmt.Errorf("failed to write output file %s: %w", tmplCfg.Output, err)
		}
		logging.Info("Generated: %s", tmplCfg.Output)
	} else {
		// Print to stdout if no output file specified
		fmt.Print(string(output))
	}

	return nil
}

// prepareLiquidData creates the data bindings for Liquid templates.
// Variable names match the Ruby implementation for template compatibility.
func (g *Generator) prepareLiquidData(settings *config.Settings) map[string]any {
	data := map[string]any{
		// Task data - matches Ruby's tasks.referenced and tasks.unreferenced
		"tasks": g.tasks.ToLiquid(),

		// Commit statistics
		"pac_c_count":        g.commits.Count(),
		"pac_c_referenced":   g.commits.CountWith(),
		"pac_c_unreferenced": g.commits.CountWithout(),
		"pac_health":         g.commits.Health(),

		// Custom properties from settings
		"properties": settings.Properties,
	}

	return data
}

// GenerateToString renders a single template and returns the result as a string.
// Useful for testing or when output should not be written to a file.
func (g *Generator) GenerateToString(templateContent string, settings *config.Settings) (string, error) {
	liquidData := g.prepareLiquidData(settings)
	return g.renderer.RenderTemplate(templateContent, liquidData)
}
