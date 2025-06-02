package task

import (
	"github.com/Praqma/Praqmatic-Automated-Changelog/src/model"
)

type NoneTaskSystem struct{}

func NewNoneTaskSystem() *NoneTaskSystem {
	return &NoneTaskSystem{}
}

func (t NoneTaskSystem) ProcessCommits(commits *model.PACCommitCollection, taskSystem model.TaskSystem) (*model.PACTaskCollection, error) {
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