package report

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Praqma/Praqmatic-Automated-Changelog/src/logging"
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

	logging.Verbose("Starting report generation with %d template(s)", len(settings.Templates))

	// Process each template format
	for i, template := range settings.Templates {
		logging.Verbose("Processing template %d/%d", i+1, len(settings.Templates))
		
		templateFile := template.Location
		var outputFile string
		
		// Check if template file is specified
		if templateFile == "" {
			return fmt.Errorf("template file not specified")
		}
		logging.Verbose("Template location: %s", templateFile)
		
		if template.Output == "" { 
			outputFile = strings.Split(templateFile, "/")[len(strings.Split(templateFile, "/"))-1]
			logging.Verbose("No output specified, using template filename: %s", outputFile)
		} else {
			outputFile = template.Output
			logging.Verbose("Output file: %s", outputFile)
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
	logging.Verbose("Generating report from template: %s", templateFile)
	
	// Create the template data
	data := g.createTemplateData(settings.Properties)
	logging.Verbose("Template data prepared with %d properties", len(settings.Properties))

	// Look for the template in Praqmatic-Automated-Changelog templates directory first
	logging.Verbose("Reading template file: %s", templateFile)
	templateContent, err := g.findAndReadTemplate(templateFile)
	if err != nil {
		return err
	}
	logging.Verbose("Template file read successfully (%d bytes)", len(templateContent))

	// Create template with functions
	tmpl := template.New(filepath.Base(templateFile))
	
	// Add custom functions before parsing
	tmpl = tmpl.Funcs(GetTemplateFuncs())
	logging.Verbose("Added custom template functions")

	// Parse the template
	logging.Verbose("Parsing template")
	tmpl, err = tmpl.Parse(templateContent)
	if err != nil {
		return fmt.Errorf("error parsing template: %w", err)
	}
	logging.Verbose("Template parsed successfully")

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputFile)
	if outputDir != "." && outputDir != "" {
		logging.Verbose("Creating output directory: %s", outputDir)
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return fmt.Errorf("error creating output directory: %w", err)
		}
	}

	// Create the output file
	logging.Verbose("Creating output file: %s", outputFile)
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating output file: %w", err)
	}
	defer file.Close()

	// Execute the template
	logging.Verbose("Executing template")
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("error executing template: %w", err)
	}

	fmt.Printf("[PAC] Generated report: %s\n", outputFile)
	logging.Verbose("Report generated successfully: %s", outputFile)
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
