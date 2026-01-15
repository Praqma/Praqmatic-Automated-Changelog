package task

import (
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/config"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/model"
)

// NoneTaskSystem is a task system that only uses regex matching.
// It doesn't fetch data from any external system.
type NoneTaskSystem struct {
	config config.TaskSystemConfig
}

// NewNoneTaskSystem creates a new NoneTaskSystem.
func NewNoneTaskSystem(cfg *config.TaskSystemConfig) *NoneTaskSystem {
	return &NoneTaskSystem{config: *cfg}
}

// Name returns the task system name.
func (n *NoneTaskSystem) Name() string {
	return n.config.Name
}

// Apply processes tasks that match this system's patterns.
// For NoneTaskSystem, this is a no-op since we don't fetch external data.
// Labels are already applied during task ID extraction.
func (n *NoneTaskSystem) Apply(tasks *model.PACTaskCollection) error {
	// NoneTaskSystem doesn't fetch external data
	// Tasks are already labeled during extraction
	return nil
}
