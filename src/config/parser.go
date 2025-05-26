package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/Praqma/Praqmatic-Automated-Changelog/src/model"
)

// ReadSettingsFile reads the settings file based on command line input
func ReadSettingsFile(settingsPath string) ([]byte, error) {

	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("settings file '%s' does not exist", settingsPath)
	}

	return os.ReadFile(settingsPath)
}

// GenerateSettings creates the final settings based on command line arguments
func GenerateSettings(settingsPath string) (*model.Settings, error) {
	settings := model.NewSettings()

	configData, err := ReadSettingsFile(settingsPath)
	if err != nil {
		return nil, fmt.Errorf("error reading settings file: %w", err)
	}

	err = yaml.Unmarshal(configData, &settings)
	if err != nil {
		return nil, fmt.Errorf("error parsing YAML settings: %w", err)
	}

	// Initialize properties if it's nil
	if settings.Properties == nil {
		settings.Properties = make(map[string]any)
	}
	return settings, nil
}

// 	// Handle username/password overrides
// 	if cFlag, ok := cmdArgs["-c"].(bool); ok && cFlag {
// 		users, _ := cmdArgs["<user>"].([]string)
// 		passwords, _ := cmdArgs["<password>"].([]string)
// 		targets, _ := cmdArgs["<target>"].([]string)

// 		for i := 0; i < len(users) && i < len(passwords) && i < len(targets); i++ {
// 			user := users[i]
// 			password := passwords[i]
// 			target := targets[i]

// 			for j, system := range settings.TaskSystems {
// 				if system.Name == target {
// 					settings.TaskSystems[j].Username = user
// 					settings.TaskSystems[j].Password = password
// 				}
// 			}
// 		}
// 	}

// 	// Handle additional properties
// 	if propsJson, ok := cmdArgs["--properties"].(string); ok && propsJson != "" {
// 		var jsonProps map[string]any
// 		if err := json.Unmarshal([]byte(propsJson), &jsonProps); err != nil {
// 			return nil, fmt.Errorf("error parsing JSON properties: %w", err)
// 		}

// 		// Merge JSON properties into existing properties
// 		for k, v := range jsonProps {
// 			settings.Properties[k] = v
// 		}
// 	}

// 	// Set verbosity level
// 	verbosity := 1
// 	if vFlags, ok := cmdArgs["-v"].(int); ok {
// 		verbosity += vFlags
// 	}
// 	if qFlags, ok := cmdArgs["-q"].(int); ok {
// 		verbosity -= qFlags
// 	}
// 	settings.Verbosity = verbosity

// 	return settings, nil
// }
