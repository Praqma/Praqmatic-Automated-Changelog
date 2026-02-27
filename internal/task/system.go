// Package task provides task system integrations for PAC.
package task

import (
	"fmt"
	"strings"

	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/logging"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/model"
)

// TaskSystem defines the interface for external task tracking systems.
// Implementations include Jira, GitHub Issues, or regex-only matching.
type System interface {
	// Apply enriches tasks with data from the external system.
	// Tasks that match this system's patterns will have their
	// Data and Attributes fields populated.
	Apply(tasks *model.PACTaskCollection) error

	// Name returns the identifier for this task system.
	Name() string
}

// FetchFunc is a function type for fetching task data from an external system.
type FetchFunc func(taskID string) (map[string]any, error)

// applyTasks is a shared helper that iterates over tasks and applies data from
// an external system. It handles the common logic of skipping irrelevant tasks,
// fetching data, and collecting errors.
func applyTasks(tasks *model.PACTaskCollection, systemName string, fetch FetchFunc) error {
	var errors []string

	for _, task := range tasks.Tasks {
		// Skip unreferenced commits (empty task ID)
		if task.TaskID == "" {
			continue
		}

		// Skip tasks that don't apply to this system
		if !task.AppliesTo[systemName] {
			continue
		}

		// Fetch data from the external system
		data, err := fetch(task.TaskID)
		if err != nil {
			logging.Warn("%s error for %s: %v", systemName, task.TaskID, err)
			// Mark task as unknown when fetch fails
			task.ClearLabels()
			task.AddLabel("unknown")
			errors = append(errors, fmt.Sprintf("%s: %v", task.TaskID, err))
			continue
		}

		// Populate task with fetched data
		task.Data = data
		task.Attributes["data"] = data

		logging.Info("Applied %s data to %s", systemName, task.TaskID)
	}

	if len(errors) > 0 {
		return fmt.Errorf("%s errors: %s", strings.ToLower(systemName), strings.Join(errors, "; "))
	}

	return nil
}
