package commands

import (
	"encoding/json"
	"strings"

	"github.com/Praqma/Praqmatic-Automated-Changelog/src/model"
	"github.com/spf13/cobra"
)

// SettingsBuilder handles building settings from command flags
type SettingsBuilder struct {
	cmd *cobra.Command
}

// NewSettingsBuilder creates a new settings builder
func NewSettingsBuilder(cmd *cobra.Command) *SettingsBuilder {
	return &SettingsBuilder{cmd: cmd}
}

// BuildFromFlags builds settings from command flags
func (sb *SettingsBuilder) BuildFromFlags() *model.Settings {
	settings := model.NewSettings()
	
	// VCS settings
	if sb.hasFlag("repo") {
		settings.VCS.Repo = sb.getString("repo")
	}
	
	if sb.hasFlag("gh-token") {
		settings.VCS.Token = sb.getString("gh-token")
	}
	
	// General settings
	if sb.hasFlag("strict") {
		settings.General.Strict = sb.getBool("strict")
	}
	
	// Templates
	if sb.hasFlag("template") {
		settings.Templates = sb.buildTemplates()
	}
	
	// Properties
	settings.Properties = sb.buildProperties()
	
	// Add default task system if none specified
	if len(settings.TaskSystems) == 0 {
		settings.TaskSystems = []model.TaskSystem{
			{
				Name: "none",
				Regex: []model.RegexRule{
					{Pattern: `(#\d+)`},
				},
			},
		}
	}
	
	return settings
}

// ApplyFlagsToSettings applies command flags to existing settings
func (sb *SettingsBuilder) ApplyFlagsToSettings(settings *model.Settings) {
	// VCS settings
	if sb.hasFlag("repo") {
		settings.VCS.Repo = sb.getString("repo")
	}
	
	if sb.hasFlag("gh-token") {
		settings.VCS.Token = sb.getString("gh-token")
	}
	
	// General settings
	if sb.hasFlag("strict") {
		settings.General.Strict = sb.getBool("strict")
	}
	
	// Templates (replace if specified)
	if sb.hasFlag("template") {
		settings.Templates = sb.buildTemplates()
	}
	
	// Properties (merge)
	if settings.Properties == nil {
		settings.Properties = make(map[string]any)
	}
	
	for k, v := range sb.buildProperties() {
		settings.Properties[k] = v
	}
}

// buildTemplates builds template configurations from flags
func (sb *SettingsBuilder) buildTemplates() []model.Template {
	templates := sb.getStringSlice("template")
	outputs := sb.getStringSlice("template-output")
	
	result := []model.Template{}
	for i, tmpl := range templates {
		output := ""
		if i < len(outputs) {
			output = outputs[i]
		}
		result = append(result, model.Template{
			Location: tmpl,
			Output:   output,
		})
	}
	
	return result
}

// buildProperties builds properties from multiple flag sources
func (sb *SettingsBuilder) buildProperties() map[string]any {
	props := make(map[string]any)
	
	// 1. Properties from JSON string
	if sb.hasFlag("properties-json") {
		jsonStr := sb.getString("properties-json")
		var jsonProps map[string]any
		if err := json.Unmarshal([]byte(jsonStr), &jsonProps); err == nil {
			for k, v := range jsonProps {
				props[k] = v
			}
		}
	}
	
	// 2. Properties from --property flag
	if sb.hasFlag("property") {
		for k, v := range sb.getStringMap("property") {
			props[k] = parseValue(v)
		}
	}
	
	// 3. Properties from --set flag (alias)
	if sb.hasFlag("set") {
		for k, v := range sb.getStringMap("set") {
			setNestedProperty(props, k, v)
		}
	}
	
	return props
}

// Helper methods
func (sb *SettingsBuilder) hasFlag(name string) bool {
	return sb.cmd.Flags().Changed(name)
}

func (sb *SettingsBuilder) getString(name string) string {
	val, _ := sb.cmd.Flags().GetString(name)
	return val
}

func (sb *SettingsBuilder) getBool(name string) bool {
	val, _ := sb.cmd.Flags().GetBool(name)
	return val
}

func (sb *SettingsBuilder) getStringSlice(name string) []string {
	val, _ := sb.cmd.Flags().GetStringSlice(name)
	return val
}

func (sb *SettingsBuilder) getStringMap(name string) map[string]string {
	val, _ := sb.cmd.Flags().GetStringToString(name)
	return val
}

// parseValue attempts to parse a string value as JSON, falling back to string
func parseValue(value string) any {
	var parsed any
	if err := json.Unmarshal([]byte(value), &parsed); err == nil {
		return parsed
	}
	return value
}

// setNestedProperty sets a nested property using dot notation
func setNestedProperty(props map[string]any, key string, value string) {
	keys := strings.Split(key, ".")
	current := props
	
	for i, k := range keys {
		if i == len(keys)-1 {
			// Last key, set the value
			current[k] = parseValue(value)
		} else {
			// Intermediate key, create map if doesn't exist
			if _, ok := current[k]; !ok {
				current[k] = make(map[string]any)
			}
			if next, ok := current[k].(map[string]any); ok {
				current = next
			} else {
				// Can't traverse further, overwrite with new map
				current[k] = make(map[string]any)
				current = current[k].(map[string]any)
			}
		}
	}
}
