package model

import (
	"encoding/json"
	"log"
)

// PACTaskCollection represents a collection of PAC tasks
type PACTaskCollection struct {
	Tasks []*PACTask
}

func (tc *PACTaskCollection) ToJSON() []byte {
	data, err := json.MarshalIndent(tc.Tasks, "", "  ")
	if err != nil {
		log.Printf("PACTaskCollection.ToJSON error: %v", err)
		return nil
	}
	return data
}

// NewPACTaskCollection creates a new empty task collection
func NewPACTaskCollection() *PACTaskCollection {
	return &PACTaskCollection{
		Tasks: make([]*PACTask, 0),
	}
}

// Add a task or multiple tasks to the collection
func (tc *PACTaskCollection) Add(tasks ...*PACTask) {
	for _, task := range tasks {
		found := false
		for i, existingTask := range tc.Tasks {
			if task.Equal(existingTask) {
				found = true
				// Add commits from the new task to the existing one
				for _, commit := range task.Commits {
					tc.Tasks[i].AddCommit(commit)
				}
				break
			}
		}
		if !found {
			tc.Tasks = append(tc.Tasks, task)
		}
	}
}

// GetTaskByID returns a task by its ID
func (tc *PACTaskCollection) GetTaskByID(id string) *PACTask {
	for _, task := range tc.Tasks {
		if task.TaskID == id {
			return task
		}
	}
	return nil
}

// GetReferencedTasks returns all tasks that have a non-nil task ID
func (tc *PACTaskCollection) GetReferencedTasks() []*PACTask {
	var result []*PACTask
	for _, task := range tc.Tasks {
		if task.TaskID != "" {
			result = append(result, task)
		}
	}
	return result
}

// GetUnreferencedCommits returns all commits that don't belong to a specific task
func (tc *PACTaskCollection) GetUnreferencedCommits() []*PACCommit {
	var result []*PACCommit
	for _, task := range tc.Tasks {
		if task.TaskID == "" {
			result = append(result, task.Commits...)
		}
	}
	return result
}

// GetTasksByLabel organizes tasks by their labels
func (tc *PACTaskCollection) GetTasksByLabel() map[string][]*PACTask {
	result := make(map[string][]*PACTask)

	for _, task := range tc.Tasks {
		for _, label := range task.Labels {
			if _, exists := result[label]; !exists {
				result[label] = make([]*PACTask, 0)
			}
			result[label] = append(result[label], task)
		}
	}

	return result
}
