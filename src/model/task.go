package model

import "slices"

// PACTask represents a task from a task management system
type PACTask struct {
	TaskID    string
	Commits    []*PACCommit
	Attributes map[string]interface{}
	AppliesTo  []string
	Labels     []string
	Data       interface{}
}

// NewPACTask creates a new task with the given ID
func NewPACTask(TaskID string) *PACTask {
	return &PACTask{
		TaskID:    TaskID,
		Commits:    make([]*PACCommit, 0),
		Attributes: make(map[string]interface{}),
		AppliesTo:  make([]string, 0),
		Labels:     make([]string, 0),
	}
}

// AddCommit adds a commit to this task
func (t *PACTask) AddCommit(commit *PACCommit) {
	// Check if commit already exists to avoid duplicates
	for _, c := range t.Commits {
		if c.SHA == commit.SHA {
			return
		}
	}
	t.Commits = append(t.Commits, commit)
}

// AddAppliesTo adds a system name this task applies to
func (t *PACTask) AddAppliesTo(system string) {
	// Check for duplicates
	if slices.Contains(t.AppliesTo, system) {
			return
		}
	t.AppliesTo = append(t.AppliesTo, system)
}

// AddLabel adds a label to this task
func (t *PACTask) AddLabel(label string) {
	if label == "" {
		return
	}

	// Check for duplicates
	if slices.Contains(t.Labels, label) {
			return
		}
	t.Labels = append(t.Labels, label)
}

// ClearLabels removes all labels from this task
func (t *PACTask) ClearLabels() {
	t.Labels = make([]string, 0)
}

// Equal checks if two tasks are the same based on their ID
func (t *PACTask) Equal(other *PACTask) bool {
	return t.TaskID == other.TaskID
}

// PACTaskCollection represents a collection of PAC tasks
type PACTaskCollection struct {
	Tasks []*PACTask
}

// NewPACTaskCollection creates a new empty task collection
func NewPACTaskCollection() *PACTaskCollection {
	return &PACTaskCollection{
		Tasks: make([]*PACTask, 0),
	}
}

// Add adds a task or multiple tasks to the collection
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
