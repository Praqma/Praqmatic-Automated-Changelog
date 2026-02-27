package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

const (
	// DefaultVCSType is the default version control system type
	DefaultVCSType = "git"
)

// LoadSettings reads configuration from a YAML file and applies overrides.
// Supports both legacy Ruby format (with colon prefixes like :key) and
// modern format (without prefixes). Legacy keys are automatically normalized.
func LoadSettings(configPath string, overrides map[string]any) (*Settings, error) {
	v := viper.New()

	v.SetDefault("general.strict", false)
	v.SetDefault("vcs.type", DefaultVCSType)

	// Read and normalize the YAML file to handle legacy colon-prefixed keys
	normalizedConfig, err := readAndNormalizeYAML(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	// Merge the normalized config into viper
	if err := v.MergeConfigMap(normalizedConfig); err != nil {
		return nil, fmt.Errorf("failed to merge config: %w", err)
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

// readAndNormalizeYAML reads a YAML file and normalizes keys by removing colon prefixes.
// This allows supporting both legacy Ruby-style (:key) and modern (key) formats.
func readAndNormalizeYAML(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var raw map[string]any
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return nil, err
	}

	return normalizeKeys(raw), nil
}

// normalizeKeys recursively removes colon prefixes from map keys.
func normalizeKeys(m map[string]any) map[string]any {
	result := make(map[string]any, len(m))

	for key, value := range m {
		// Remove leading colon if present (legacy Ruby symbol format)
		normalizedKey := strings.TrimPrefix(key, ":")

		switch v := value.(type) {
		case map[string]any:
			result[normalizedKey] = normalizeKeys(v)
		case []any:
			result[normalizedKey] = normalizeSlice(v)
		default:
			result[normalizedKey] = value
		}
	}

	return result
}

// normalizeSlice recursively normalizes keys in slice elements.
func normalizeSlice(s []any) []any {
	result := make([]any, len(s))

	for i, item := range s {
		switch v := item.(type) {
		case map[string]any:
			result[i] = normalizeKeys(v)
		case []any:
			result[i] = normalizeSlice(v)
		default:
			result[i] = item
		}
	}

	return result
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
// Format: -c user -c password -c target
// or: -c token -c target (for token-based auth)
func ParseCredentialFlags(credentials []string) ([]CredentialOverride, error) {
	if len(credentials) == 0 {
		return nil, nil
	}

	if len(credentials)%3 != 0 && len(credentials)%2 != 0 {
		return nil, fmt.Errorf("credentials must be provided in groups of 3: '-c user -c password -c target' or groups of 2: '-c token -c target'")
	}

	var overrides []CredentialOverride

	if len(credentials)%3 == 0 {
		for i := 0; i < len(credentials); i += 3 {
			overrides = append(overrides, CredentialOverride{
				Username: credentials[i],
				Password: credentials[i+1],
				Target:   credentials[i+2],
			})
		}
	} else {
		for i := 0; i < len(credentials); i += 2 {
			overrides = append(overrides, CredentialOverride{
				Password: credentials[i],
				Target:   credentials[i+1],
			})
		}
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
		s.VCS.Type = DefaultVCSType // Default to git
	}

	if s.VCS.Type != DefaultVCSType {
		return fmt.Errorf("unsupported VCS type: %s (only %q is supported)", s.VCS.Type, DefaultVCSType)
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
