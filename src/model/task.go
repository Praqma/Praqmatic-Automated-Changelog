package model

import (
	"slices"
)

// PACTask represents a task from a task management system
type PACTask struct {
	ID      	string
	Title       string
	Commits     []*PACCommit
	Labels      []string
	Assignees   []string
	Author      string
	URL         string
	Data 	 	map[string]interface{}
}

func NewPACTask(ID string) *PACTask {
	return &PACTask{
		ID:     ID,
		Commits:    make([]*PACCommit, 0),
		Labels:     make([]string, 0),
		Assignees:  make([]string, 0),
		Data:       make(map[string]interface{}),
	}
}

// AddCommit adds a commit to this task
func (t *PACTask) AddCommit(commit *PACCommit) {
	for _, c := range t.Commits {
		if c.SHA == commit.SHA {
			return
		}
	}
	t.Commits = append(t.Commits, commit)
}

// AddLabel adds a label to this task
func (t *PACTask) AddLabel(label string) {
	if label == "" {
		return
	}

	if slices.Contains(t.Labels, label) {
		return
	}
	t.Labels = append(t.Labels, label)
}

// AddAssignee adds an assignee to this task
func (t *PACTask) AddAssignee(assignee string) {
	if assignee == "" {
		return
	}

	if slices.Contains(t.Assignees, assignee) {
		return
	}
	t.Assignees = append(t.Assignees, assignee)
}

// ClearLabels removes all labels from this task
func (t *PACTask) ClearLabels() {
	t.Labels = make([]string, 0)
}

// Equal checks if two tasks are the same based on their ID
func (t *PACTask) Equal(other *PACTask) bool {
	return t.ID == other.ID
}
