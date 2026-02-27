// Package cmd provides the CLI commands for PAC.
package cmd

import (
	"github.com/spf13/cobra"
)

// Global flags
var (
	settingsFile string
	properties   string
	verbosity    int
	quiet        int
	credentials  []string
)

// Version is set at build time.
var Version = "dev"

// rootCmd is the base command for PAC.
var rootCmd = &cobra.Command{
	Use:   "pac",
	Short: "Praqmatic Automated Changelog",
	Long: `PAC (Praqmatic Automated Changelog) generates changelogs from git commits
and task tracking systems like Jira.

It extracts task IDs from commit messages using configurable regex patterns,
optionally fetches additional data from external systems, and renders
changelogs using Liquid templates.`,
	Version: Version,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

// SetVersion sets the version string (called from main).
func SetVersion(v string) {
	Version = v
	rootCmd.Version = v
}

func init() {
	// Global flags available to all commands
	rootCmd.PersistentFlags().StringVarP(&settingsFile, "settings", "s", "settings/default_settings.yml",
		"Path to settings YAML file")

	rootCmd.PersistentFlags().StringVar(&properties, "properties", "",
		"Additional properties as JSON (merged with settings)")

	rootCmd.PersistentFlags().CountVarP(&verbosity, "verbose", "v",
		"Increase verbosity (can be repeated: -v, -vv, -vvv)")

	rootCmd.PersistentFlags().CountVarP(&quiet, "quiet", "q",
		"Decrease verbosity (can be repeated)")

	rootCmd.PersistentFlags().StringArrayVarP(&credentials, "credentials", "c", []string{},
		"Override credentials: -c user -c password -c target (in groups of 3) or -c token -c target (in groups of 2)")
}
