package cmd

import (
	"fmt"

	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/core"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/logging"
	"github.com/spf13/cobra"
)

var fromLatestTagCmd = &cobra.Command{
	Use:   "from-latest-tag <pattern>",
	Short: "Generate changelog from the latest tag matching a pattern",
	Long: `Generate a changelog from the latest tag matching the given glob pattern to HEAD.

The pattern uses glob-style matching (e.g., "v*", "release-*").
The command finds the most recent tag matching the pattern and generates
a changelog from that tag to HEAD.

Examples:
  pac from-latest-tag "v*"
  pac from-latest-tag "release-*"
  pac from-latest-tag "v*" --to develop`,
	Args: cobra.ExactArgs(1),
	RunE: runFromLatestTag,
}

func init() {
	fromLatestTagCmd.Flags().StringVar(&toRef, "to", "",
		"Newest reference (default: HEAD)")

	rootCmd.AddCommand(fromLatestTagCmd)
}

func runFromLatestTag(cmd *cobra.Command, args []string) error {
	pattern := args[0]

	// Set up logging
	logging.SetVerbosity(verbosity - quiet)

	// Load configuration
	settings, err := loadSettings()
	if err != nil {
		return err
	}

	// Validate settings
	if validateErr := settings.Validate(); validateErr != nil {
		return fmt.Errorf("invalid settings: %w", validateErr)
	}

	// Run the core workflow with tag pattern
	_, err = core.RunFromLatestTag(settings, pattern, toRef)
	return err
}
