package model

// PACTask represents a task extracted from commit messages.
type PACTask struct {
	TaskID     string                 // Task identifier (e.g., "JIRA-123")
	Commits    *PACCommitCollection   // Commits associated with this task
	Attributes map[string]any         // Additional attributes from external systems
	AppliesTo  map[string]bool        // Which task systems this task applies to
	Labels     map[string]bool        // Labels applied to this task (from regex matches)
	Data       any                    // Raw data from external system (typically JSON)
}

// PACTaskCollection holds a collection of tasks with an index for quick lookup.
type PACTaskCollection struct {
	Tasks     []*PACTask
	taskIndex map[string]*PACTask
}

// NewPACTask creates a new task with the given ID.
func NewPACTask(taskID string) *PACTask {
	return &PACTask{
		TaskID:     taskID,
		Commits:    NewPACCommitCollection(),
		Attributes: make(map[string]any),
		AppliesTo:  make(map[string]bool),
		Labels:     make(map[string]bool),
	}
}

// NewPACTaskCollection creates a new empty task collection.
func NewPACTaskCollection() *PACTaskCollection {
	return &PACTaskCollection{
		Tasks:     []*PACTask{},
		taskIndex: make(map[string]*PACTask),
	}
}

// AddCommit adds a commit to this task.
func (t *PACTask) AddCommit(commit *PACCommit) {
	t.Commits.Add(commit)
}

// AddLabel adds a label to this task.
func (t *PACTask) AddLabel(label string) {
	t.Labels[label] = true
}

// ClearLabels removes all labels from this task.
func (t *PACTask) ClearLabels() {
	t.Labels = make(map[string]bool)
}

// FindOrCreate returns an existing task or creates a new one.
func (tc *PACTaskCollection) FindOrCreate(taskID string) *PACTask {
	if task, exists := tc.taskIndex[taskID]; exists {
		return task
	}

	task := NewPACTask(taskID)
	tc.Tasks = append(tc.Tasks, task)
	tc.taskIndex[taskID] = task
	return task
}

// Count returns the total number of tasks.
func (tc *PACTaskCollection) Count() int {
	return len(tc.Tasks)
}

// Referenced returns tasks that have a non-empty TaskID.
func (tc *PACTaskCollection) Referenced() []*PACTask {
	result := make([]*PACTask, 0)
	for _, task := range tc.Tasks {
		if task.TaskID != "" {
			result = append(result, task)
		}
	}
	return result
}

// Unreferenced returns the task with empty TaskID (contains commits without task references).
// Returns nil if no unreferenced commits exist.
func (tc *PACTaskCollection) Unreferenced() *PACTask {
	if task, exists := tc.taskIndex[""]; exists {
		return task
	}
	return nil
}

// UnreferencedCommits returns all commits that don't reference any task.
func (tc *PACTaskCollection) UnreferencedCommits() []*PACCommit {
	unref := tc.Unreferenced()
	if unref == nil {
		return []*PACCommit{}
	}
	return unref.Commits.Commits
}

// GetLabels returns all labels as a slice.
func (t *PACTask) GetLabels() []string {
	labels := make([]string, 0, len(t.Labels))
	for label := range t.Labels {
		labels = append(labels, label)
	}
	return labels
}

// HasLabel checks if the task has a specific label.
func (t *PACTask) HasLabel(label string) bool {
	return t.Labels[label]
}
