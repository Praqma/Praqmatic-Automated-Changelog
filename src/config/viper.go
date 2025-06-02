package config

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/Praqma/Praqmatic-Automated-Changelog/src/logging"
	"github.com/Praqma/Praqmatic-Automated-Changelog/src/model"
)

// MergeViperSettings merges configuration from Viper into the settings struct
func MergeViperSettings(settings *model.Settings) error {
	logging.Verbose("Merging viper configuration into settings")
	
	// Unmarshal entire config
	if err := viper.Unmarshal(settings); err != nil {
		return fmt.Errorf("error unmarshaling config: %w", err)
	}
	logging.Verbose("Base configuration unmarshaled")
	
	// Handle special cases for complex structures
	if viper.IsSet("templates") {
		logging.Verbose("Found templates configuration")
		var templates []model.Template
		if err := viper.UnmarshalKey("templates", &templates); err != nil {
			return fmt.Errorf("error unmarshaling templates: %w", err)
		}
		settings.Templates = templates
		logging.Verbose("Loaded %d template(s)", len(templates))
	}
	
	if viper.IsSet("task_systems") {
		logging.Verbose("Found task systems configuration")
		var taskSystems []model.TaskSystem
		if err := viper.UnmarshalKey("task_systems", &taskSystems); err != nil {
			return fmt.Errorf("error unmarshaling task systems: %w", err)
		}
		settings.TaskSystems = taskSystems
		logging.Verbose("Loaded %d task system(s)", len(taskSystems))
	}
	
	if viper.IsSet("properties") {
		logging.Verbose("Found properties configuration")
		props := viper.GetStringMap("properties")
		if settings.Properties == nil {
			settings.Properties = make(map[string]any)
		}
		for k, v := range props {
			settings.Properties[k] = v
			logging.Verbose("Set property: %s = %v", k, v)
		}
	}
	
	logging.Verbose("Configuration merge completed")
	return nil
}
