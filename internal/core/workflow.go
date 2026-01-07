// Package core provides the main PAC workflow orchestration.package core

package core

import (
	"fmt"

	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/config"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/logging"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/model"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/report"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/task"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/vcs"
)

// Result contains the outcome of a PAC workflow run.
type Result struct {
	Commits *model.PACCommitCollection
	Tasks   *model.PACTaskCollection
	AllOK   bool
}

// Run executes the complete PAC workflow:
// 1. Initialize VCS and get commit delta
// 2. Build task collection from commits
// 3. Apply task systems to fetch external data
// 4. Generate reports from templates
func Run(settings *config.Settings, oldestRef, newestRef string) (*Result, error) {
	// Initialize logging
	logging.SetVerbosity(settings.Verbosity)

	logging.Info("PAC - Praqmatic Automated Changelog")
	logging.Info("From: %s, To: %s", oldestRef, defaultStr(newestRef, "HEAD"))

	// Initialize VCS
	gitVCS, err := vcs.NewGitVCS(settings.VCS)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize VCS: %w", err)
	}

	// Get commit delta
	logging.Info("Collecting commits...")
	commits, err := gitVCS.GetDelta(oldestRef, newestRef)
	if err != nil {
		return nil, fmt.Errorf("failed to get commit delta: %w", err)
	}

	logging.Verbosef(1, "Found %d commits", commits.Count())

	// Build task collection from commits
	tasks := task.BuildTaskCollection(commits, settings.TaskSystems)

	logging.Verbosef(1, "Found %d referenced tasks", len(tasks.Referenced()))
	logging.Verbosef(1, "Found %d unreferenced commits", len(tasks.UnreferencedCommits()))

	// Apply task systems (fetch external data like Jira)
	allOK := applyTaskSystems(settings, tasks)

	// Generate reports
	generator := report.NewGenerator(tasks, commits)
	if err := generator.Generate(settings); err != nil {
		return nil, fmt.Errorf("failed to generate reports: %w", err)
	}

	// Handle strict mode
	if !allOK && settings.General.Strict {
		return nil, fmt.Errorf("errors encountered in strict mode")
	}

	if !allOK {
		logging.Verbosef(1, "Ignoring encountered errors. Strict mode is disabled.")
	}

	logging.Info("Done! Health: %.1f%% commits reference tasks", commits.Health())

	return &Result{
		Commits: commits,
		Tasks:   tasks,
		AllOK:   allOK,
	}, nil
}

// RunFromLatestTag finds the latest tag matching the pattern and runs the workflow.
func RunFromLatestTag(settings *config.Settings, pattern, newestRef string) (*Result, error) {
	logging.SetVerbosity(settings.Verbosity)

	// Initialize VCS to find latest tag
	gitVCS, err := vcs.NewGitVCS(settings.VCS)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize VCS: %w", err)
	}

	// Find latest matching tag
	latestTag, err := gitVCS.GetLatestTag(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to find latest tag matching %q: %w", pattern, err)
	}

	logging.Info("Found latest tag: %s", latestTag)

	// Continue with normal workflow
	return Run(settings, latestTag, newestRef)
}

// applyTaskSystems applies all configured task systems to the task collection.
// Returns true if all task systems completed successfully, false otherwise.
func applyTaskSystems(settings *config.Settings, tasks *model.PACTaskCollection) bool {
	allOK := true

	for _, tsCfg := range settings.TaskSystems {
		ts, err := task.CreateTaskSystem(tsCfg)
		if err != nil {
			logging.Warn("Failed to create task system %s: %v", tsCfg.Name, err)
			allOK = false
			continue
		}

		logging.Verbosef(1, "Applying task system: %s", tsCfg.Name)

		if err := ts.Apply(tasks); err != nil {
			logging.Warn("Task system %s encountered errors: %v", tsCfg.Name, err)
			allOK = false
		}
	}

	return allOK
}

func defaultStr(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
