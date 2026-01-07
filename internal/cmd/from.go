package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/config"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/core"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/logging"
	"github.com/spf13/cobra"
)

var toRef string

var fromCmd = &cobra.Command{
	Use:   "from <oldest-ref>",
	Short: "Generate changelog from a specific git reference",
	Long: `Generate a changelog from the specified git reference to HEAD (or --to ref).

The oldest-ref can be a tag name, branch name, or commit SHA.
Commits are collected exclusively (oldest-ref is not included).

Examples:
  pac from v1.0.0
  pac from v1.0.0 --to v1.1.0
  pac from abc123 --settings my_settings.yml`,
	Args: cobra.ExactArgs(1),
	RunE: runFrom,
}

func init() {
	fromCmd.Flags().StringVar(&toRef, "to", "",
		"Newest reference (default: HEAD)")

	rootCmd.AddCommand(fromCmd)
}

func runFrom(cmd *cobra.Command, args []string) error {
	oldestRef := args[0]

	// Set up logging
	logging.SetVerbosity(verbosity - quiet)

	// Load configuration
	settings, err := loadSettings()
	if err != nil {
		return err
	}

	// Validate settings
	if err := settings.Validate(); err != nil {
		return fmt.Errorf("invalid settings: %w", err)
	}

	// Run the core workflow
	_, err = core.Run(settings, oldestRef, toRef)
	return err
}

func loadSettings() (*config.Settings, error) {
	overrides := buildOverrides()

	settings, err := config.LoadSettings(settingsFile, overrides)
	if err != nil {
		return nil, fmt.Errorf("failed to load settings from %s: %w", settingsFile, err)
	}

	// Apply verbosity from CLI
	settings.Verbosity = verbosity - quiet

	return settings, nil
}

func buildOverrides() map[string]any {
	overrides := make(map[string]any)

	// Parse properties JSON if provided
	if properties != "" {
		var props map[string]any
		if err := json.Unmarshal([]byte(properties), &props); err != nil {
			logging.Warn("Failed to parse --properties JSON: %v", err)
		} else {
			overrides[":properties"] = props
		}
	}

	// Parse credentials if provided
	if len(credentials) > 0 {
		creds, err := config.ParseCredentialFlags(credentials)
		if err != nil {
			logging.Warn("Failed to parse credentials: %v", err)
		} else {
			overrides["credential_overrides"] = creds
		}
	}

	return overrides
}
