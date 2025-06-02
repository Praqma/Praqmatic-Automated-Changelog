package config

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/Praqma/Praqmatic-Automated-Changelog/src/model"
)

// MergeViperSettings merges configuration from Viper into the settings struct
func MergeViperSettings(settings *model.Settings) error {
	// Unmarshal entire config
	if err := viper.Unmarshal(settings); err != nil {
		return fmt.Errorf("error unmarshaling config: %w", err)
	}
	
	// Handle special cases for complex structures
	if viper.IsSet("templates") {
		var templates []model.Template
		if err := viper.UnmarshalKey("templates", &templates); err != nil {
			return fmt.Errorf("error unmarshaling templates: %w", err)
		}
		settings.Templates = templates
	}
	
	if viper.IsSet("task_systems") {
		var taskSystems []model.TaskSystem
		if err := viper.UnmarshalKey("task_systems", &taskSystems); err != nil {
			return fmt.Errorf("error unmarshaling task systems: %w", err)
		}
		settings.TaskSystems = taskSystems
	}
	
	if viper.IsSet("properties") {
		props := viper.GetStringMap("properties")
		if settings.Properties == nil {
			settings.Properties = make(map[string]any)
		}
		for k, v := range props {
			settings.Properties[k] = v
		}
	}
	
	return nil
}
