package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// LoadSettings reads configuration from a YAML file and applies overrides.
// The configPath should point to a YAML file in the Ruby symbol format (e.g., :key:).
func LoadSettings(configPath string, overrides map[string]any) (*Settings, error) {
	v := viper.New()

	// Set defaults
	v.SetDefault(":general::strict", false)
	v.SetDefault(":vcs::type", "git")

	// Configure viper for the config file
	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	// Apply command-line overrides
	for key, value := range overrides {
		v.Set(key, value)
	}

	// Unmarshal into Settings struct
	var settings Settings
	if err := v.Unmarshal(&settings); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Apply credential overrides if provided
	if credOverrides, ok := overrides["credential_overrides"].([]CredentialOverride); ok {
		applyCredentialOverrides(&settings, credOverrides)
	}

	return &settings, nil
}

// CredentialOverride holds username/password override for a specific task system.
type CredentialOverride struct {
	Username string
	Password string
	Target   string // Task system name to apply to
}

// applyCredentialOverrides updates task system credentials from CLI flags.
func applyCredentialOverrides(settings *Settings, overrides []CredentialOverride) {
	for _, cred := range overrides {
		for i := range settings.TaskSystems {
			if settings.TaskSystems[i].Name == cred.Target {
				settings.TaskSystems[i].Username = cred.Username
				settings.TaskSystems[i].Password = cred.Password
			}
		}
	}
}

// ParseCredentialFlags parses -c flag arguments into CredentialOverride structs.
// Format: -c user password target
func ParseCredentialFlags(credentials []string) ([]CredentialOverride, error) {
	if len(credentials) == 0 {
		return nil, nil
	}

	if len(credentials)%3 != 0 {
		return nil, fmt.Errorf("credentials must be provided in groups of 3: user password target")
	}

	var overrides []CredentialOverride
	for i := 0; i < len(credentials); i += 3 {
		overrides = append(overrides, CredentialOverride{
			Username: credentials[i],
			Password: credentials[i+1],
			Target:   credentials[i+2],
		})
	}

	return overrides, nil
}

// Validate checks if the settings are valid and returns an error if not.
func (s *Settings) Validate() error {
	if len(s.Templates) == 0 {
		return fmt.Errorf("at least one template must be configured")
	}

	for i, tmpl := range s.Templates {
		if tmpl.Location == "" {
			return fmt.Errorf("template %d: location is required", i)
		}
	}

	if s.VCS.Type == "" {
		s.VCS.Type = "git" // Default to git
	}

	if s.VCS.Type != "git" {
		return fmt.Errorf("unsupported VCS type: %s (only 'git' is supported)", s.VCS.Type)
	}

	return nil
}

// GetTaskSystemByName returns the task system config with the given name.
func (s *Settings) GetTaskSystemByName(name string) *TaskSystemConfig {
	for i := range s.TaskSystems {
		if strings.EqualFold(s.TaskSystems[i].Name, name) {
			return &s.TaskSystems[i]
		}
	}
	return nil
}
