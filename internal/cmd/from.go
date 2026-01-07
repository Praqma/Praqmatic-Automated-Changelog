package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/config"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/logging"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/report"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/task"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/vcs"
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

	logging.Info("PAC - Praqmatic Automated Changelog")
	logging.Info("From: %s, To: %s", oldestRef, defaultString(toRef, "HEAD"))

	// Load configuration
	settings, err := loadSettings()
	if err != nil {
		return err
	}

	// Validate settings
	if err := settings.Validate(); err != nil {
		return fmt.Errorf("invalid settings: %w", err)
	}

	// Run the workflow
	return runWorkflow(settings, oldestRef, toRef)
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

func runWorkflow(settings *config.Settings, oldestRef, newestRef string) error {
	// Initialize VCS
	gitVCS, err := vcs.NewGitVCS(settings.VCS)
	if err != nil {
		return fmt.Errorf("failed to initialize git: %w", err)
	}

	// Get commit delta
	logging.Info("Collecting commits from %s to %s", oldestRef, defaultString(newestRef, "HEAD"))
	commits, err := gitVCS.GetDelta(oldestRef, newestRef)
	if err != nil {
		return fmt.Errorf("failed to get commits: %w", err)
	}

	logging.Info("Found %d commits", commits.Count())

	// Build task collection from commits
	tasks := task.BuildTaskCollection(commits, settings.TaskSystems)

	logging.Info("Found %d tasks, %d unreferenced commits",
		len(tasks.Referenced()), len(tasks.UnreferencedCommits()))

	// Apply task systems (fetch external data)
	allOK := true
	for _, tsCfg := range settings.TaskSystems {
		ts, err := task.CreateTaskSystem(tsCfg)
		if err != nil {
			logging.Warn("Failed to create task system %s: %v", tsCfg.Name, err)
			allOK = false
			continue
		}

		if err := ts.Apply(tasks); err != nil {
			logging.Warn("Task system %s encountered errors: %v", tsCfg.Name, err)
			allOK = false
		}
	}

	// Generate reports
	generator := report.NewGenerator(tasks, commits)
	if err := generator.Generate(settings); err != nil {
		return fmt.Errorf("failed to generate reports: %w", err)
	}

	// Handle strict mode
	if !allOK && settings.General.Strict {
		return fmt.Errorf("errors encountered in strict mode")
	}

	if !allOK {
		logging.Info("Ignoring encountered errors (strict mode disabled)")
	}

	logging.Info("Done! Health: %.1f%% commits reference tasks", commits.Health())

	return nil
}

func defaultString(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
