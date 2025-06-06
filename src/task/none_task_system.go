package task

import (
	"github.com/Praqma/Praqmatic-Automated-Changelog/src/model"
)

type NoneTaskSystem struct{}

func NewNoneTaskSystem() *NoneTaskSystem {
	return &NoneTaskSystem{}
}

func (t *NoneTaskSystem) ProcessCommits(taskSystem model.TaskSystem, commits *model.PACCommitCollection) (*model.PACTaskCollection, error) {
	tasks := model.NewPACTaskCollection()
	for _, commit := range commits.Commits {
		taskID := extractTaskID(commit, taskSystem.Regex)
		if taskID != "" {
			task := model.NewPACTask(taskID)
			task.AddCommit(commit)
			tasks.Add(task)
		} else {
			task := model.NewPACTask("")
			task.AddCommit(commit)
			tasks.Add(task)
		}
	}
	return tasks, nil
}

// GetName returns the name of this task system
func (t *NoneTaskSystem) GetName() string {
	return "none"
}