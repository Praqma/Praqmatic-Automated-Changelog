package task

import (
	"fmt"

	"github.com/Praqma/Praqmatic-Automated-Changelog/src/model"
)

// TaskSystemProcessor defines the interface that all task systems must implement
// This interface standardizes how task systems process commits and create tasks
type TaskSystemProcessor interface {
	// ProcessCommits processes a collection of commits and returns tasks
	// taskSystem contains the configuration for this specific task system
	// commits contains all the commits to be processed
	ProcessCommits(taskSystem model.TaskSystem, commits *model.PACCommitCollection) (*model.PACTaskCollection, error)

	// GetName returns the name/type of this task system (e.g., "github", "jira", "none")
	GetName() string
}

// TaskSystemFactory creates task system processors based on configuration
type TaskSystemFactory struct{}

// NewTaskSystemFactory creates a new factory instance
func NewTaskSystemFactory() *TaskSystemFactory {
	return &TaskSystemFactory{}
}

// CreateTaskSystem creates and returns a task system processor based on the task system configuration
func (f *TaskSystemFactory) CreateTaskSystem(taskSystem model.TaskSystem, token string) (TaskSystemProcessor, error) {
	switch taskSystem.Name {
	case "github":
		ghTaskSystem, err := NewGHTask(taskSystem.QueryString, token)
		if err != nil {
			return nil, err
		}
		return ghTaskSystem, nil
	case "none":
		return NewNoneTaskSystem(), nil
	case "jira":
		jiraTaskSystem, err := NewJiraTaskSystem(taskSystem.QueryString)
		if err != nil {
			return nil, err
		}
		return jiraTaskSystem, nil
	case "gitlab":
		// TODO: Implement GitLab task system
		return nil, fmt.Errorf("gitlab task system not implemented yet")
	case "bitbucket":
		// TODO: Implement Bitbucket task system
		return nil, fmt.Errorf("bitbucket task system not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported task system: %s", taskSystem.Name)
	}
}

