// Package task provides task system integrations for PAC.package task

package task

import "github.com/Praqma/Praqmatic-Automated-Changelog/internal/model"

// TaskSystem defines the interface for external task tracking systems.
// Implementations include Jira, GitHub Issues, or regex-only matching.
type TaskSystem interface {
	// Apply enriches tasks with data from the external system.
	// Tasks that match this system's patterns will have their
	// Data and Attributes fields populated.
	Apply(tasks *model.PACTaskCollection) error

	// Name returns the identifier for this task system.
	Name() string
}
