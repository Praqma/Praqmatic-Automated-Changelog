package task

import (
	"fmt"
	"strings"

	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/config"
)

// CreateTaskSystem creates a TaskSystem based on the configuration.
func CreateTaskSystem(cfg config.TaskSystemConfig) (TaskSystem, error) {
	name := strings.ToLower(cfg.Name)

	switch name {
	case "none":
		return NewNoneTaskSystem(cfg), nil
	case "jira":
		return NewJiraTaskSystem(cfg), nil
	default:
		// For unknown task systems, treat them like "none" (regex-only)
		// This maintains compatibility with custom named task systems
		return NewNoneTaskSystem(cfg), nil
	}
}

// CreateAllTaskSystems creates TaskSystem instances for all configured task systems.
func CreateAllTaskSystems(configs []config.TaskSystemConfig) ([]TaskSystem, error) {
	systems := make([]TaskSystem, 0, len(configs))

	for _, cfg := range configs {
		ts, err := CreateTaskSystem(cfg)
		if err != nil {
			return nil, fmt.Errorf("failed to create task system %q: %w", cfg.Name, err)
		}
		systems = append(systems, ts)
	}

	return systems, nil
}
