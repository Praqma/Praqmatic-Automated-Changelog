package model

import (
	"slices"
	"time"
)

// PACTask represents a task from a task management system
type PACTask struct {
	TaskID      string
	Title       string
	Body        string
	State       string
	TaskType    string
	Commits     []*PACCommit
	Labels      []string
	Assignees   []string
	Author      string
	URL         string
	Number      int
	Milestone   string
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	ClosedAt    *time.Time

}

func NewPACTask(TaskID string) *PACTask {
	return &PACTask{
		TaskID:     TaskID,
		Commits:    make([]*PACCommit, 0),
		Labels:     make([]string, 0),
		Assignees:  make([]string, 0),
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

// IsOpen returns true if the task is in an open state
func (t *PACTask) IsOpen() bool {
	return t.State == "open"
}

// IsClosed returns true if the task is in a closed state
func (t *PACTask) IsClosed() bool {
	return t.State == "closed"
}

// Equal checks if two tasks are the same based on their ID
func (t *PACTask) Equal(other *PACTask) bool {
	return t.TaskID == other.TaskID
}
