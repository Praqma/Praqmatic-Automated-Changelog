package cmd

import (
	"fmt"

	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/logging"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/vcs"
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

	logging.Info("PAC - Praqmatic Automated Changelog")
	logging.Info("Finding latest tag matching: %s", pattern)

	// Load configuration
	settings, err := loadSettings()
	if err != nil {
		return err
	}

	// Validate settings
	if err := settings.Validate(); err != nil {
		return fmt.Errorf("invalid settings: %w", err)
	}

	// Initialize VCS to find the latest tag
	gitVCS, err := vcs.NewGitVCS(settings.VCS)
	if err != nil {
		return fmt.Errorf("failed to initialize git: %w", err)
	}

	// Find the latest matching tag
	latestTag, err := gitVCS.GetLatestTag(pattern)
	if err != nil {
		return fmt.Errorf("failed to find tag matching %q: %w", pattern, err)
	}

	logging.Info("Found latest tag: %s", latestTag)

	// Run the workflow from the found tag
	return runWorkflow(settings, latestTag, toRef)
}
