package task

import (
	"fmt"

	"github.com/Praqma/Praqmatic-Automated-Changelog/src/model"
)

// TaskIDList generates a collection of tasks based on the commits found
func TaskIDList(settings model.Settings, commits *model.PACCommitCollection) (*model.PACTaskCollection, error) {
	tasks := model.NewPACTaskCollection()

	taskSystems := settings.TaskSystems

	for _, taskSystem := range taskSystems {
		systemTasks, err := processTaskSystem(settings, taskSystem, commits)
		if err != nil {
			return nil, fmt.Errorf("error processing task system %s: %w", taskSystem.Name, err)
		}

		// Merge results
		for _, task := range systemTasks.Tasks {
			tasks.Add(task)
		}
	}
	return tasks, nil
}

// processTaskSystem handles processing for a specific task system
func processTaskSystem(settings model.Settings, taskSystem model.TaskSystem, commits *model.PACCommitCollection) (*model.PACTaskCollection, error) {
	switch taskSystem.Name {
	case "github":
		ghTaskSystem, err := NewGHTask(taskSystem.QueryString, settings.VCS.Token)
		if err != nil {
			return nil, fmt.Errorf("error initializing GitHub task system: %w", err)
		}
		return ghTaskSystem.ProcessCommits(taskSystem, commits)
	case "jira":
		return nil, fmt.Errorf("jira is not supported yet")
	case "gitlab":
		return nil, fmt.Errorf("gitlab is not supported yet")
	case "bitbucket":
		return nil, fmt.Errorf("bitbucket is not supported yet")
	default:
		return nil, fmt.Errorf("unsupported task system: %s", taskSystem.Name)
	}
}