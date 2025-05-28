package report

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Praqma/Praqmatic-Automated-Changelog/src/model"
)

// Generator handles the generation of changelog reports
type Generator struct {
	tasks   *model.PACTaskCollection
}

// NewGenerator creates a new report generator
func NewGenerator(tasks *model.PACTaskCollection) *Generator {

	return &Generator{
		tasks:   tasks,
	}
}

// Generate generates the changelog report based on settings
func (g *Generator) Generate(settings *model.Settings) error {
	// Check if templates exist
	if len(settings.Templates) == 0 {
		return fmt.Errorf("no report templates specified in settings")
	}

	// Process each template format
	for _, template := range settings.Templates {
		templateFile := template.Location
		var outputFile string
		
		// Check if template file is specified
		if templateFile == "" {
			return fmt.Errorf("template file not specified")
		}
		if template.Output == "" { 
			outputFile = strings.Split(templateFile, "/")[len(strings.Split(templateFile, "/"))-1]
		} else {
			outputFile = template.Output
		}
		// Generate the report
		if err := g.generateReport(templateFile, outputFile, settings); err != nil {
			return fmt.Errorf("error generating report: %w", err)
		}
	}

	return nil
}

// generateReport generates a single report based on the template and settings
func (g *Generator) generateReport(templateFile, outputFile string, settings *model.Settings) error {
	// Create the template data
	data := g.createTemplateData(settings.Properties)

	// Look for the template in Praqmatic-Automated-Changelog templates directory first
	templateContent, err := g.findAndReadTemplate(templateFile)
	if err != nil {
		return err
	}

	// Create template with functions
	tmpl := template.New(filepath.Base(templateFile))
	
	// Add custom functions before parsing
	tmpl = tmpl.Funcs(GetTemplateFuncs())

	// Parse the template
	tmpl, err = tmpl.Parse(templateContent)
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputFile)
	if outputDir != "." && outputDir != "" {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("error creating output directory: %w", err)
		}
	}

	// Create the output file
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer file.Close()

	// Execute the template
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	fmt.Printf("[PAC] Generated report: %s\n", outputFile)
	return nil
}

// findAndReadTemplate searches for and reads the template file from possible locations
func (g *Generator) findAndReadTemplate(templateName string) (string, error) {
	content, err := os.ReadFile(templateName)
	if err == nil {
		return string(content), nil
	}
	return "", fmt.Errorf("template file not found: %s", templateName)
}

// createTemplateData prepares data for the template
func (g *Generator) createTemplateData(properties map[string]any) map[string]interface{} {
	data := make(map[string]interface{})

	// Pass the task collection directly
	data["tasks"] = g.tasks
	
	// Basic information
	data["date"] = time.Now()

	// Add properties from settings
	if properties != nil {
		data["properties"] = properties
	} else {
		// Default properties
		data["properties"] = map[string]any{
			"title": "PAC Changelog",
		}
	}

	return data
}
