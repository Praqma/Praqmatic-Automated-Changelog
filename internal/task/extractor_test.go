package task

import (
	"testing"

	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/config"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/model"
)

func TestBuildTaskCollection(t *testing.T) {
	commits := model.NewPACCommitCollection()
	commits.Add(&model.PACCommit{SHA: "abc123", Message: "Fix #123 issue"})
	commits.Add(&model.PACCommit{SHA: "def456", Message: "PRJ-456: New feature"})
	commits.Add(&model.PACCommit{SHA: "ghi789", Message: "Regular commit"})

	taskSystems := []config.TaskSystemConfig{
		{
			Name: "github",
			Regex: []config.RegexConfig{
				{Pattern: "/(#\\d+)/", Label: "github-issue"},
			},
		},
		{
			Name: "jira",
			Regex: []config.RegexConfig{
				{Pattern: "/(PRJ-\\d+)/", Label: "jira-task"},
			},
		},
	}

	tasks := BuildTaskCollection(commits, taskSystems)

	// Should have 3 tasks: #123, PRJ-456, and unreferenced
	if tasks.Count() != 3 {
		t.Errorf("expected 3 tasks, got %d", tasks.Count())
	}

	// Check referenced tasks
	referenced := tasks.Referenced()
	if len(referenced) != 2 {
		t.Errorf("expected 2 referenced tasks, got %d", len(referenced))
	}

	// Check unreferenced commits
	unreferenced := tasks.UnreferencedCommits()
	if len(unreferenced) != 1 {
		t.Errorf("expected 1 unreferenced commit, got %d", len(unreferenced))
	}

	// Check commit reference flags
	refCount := 0
	for _, c := range commits.Commits {
		if c.Referenced {
			refCount++
		}
	}
	if refCount != 2 {
		t.Errorf("expected 2 referenced commits, got %d", refCount)
	}
}

func TestBuildTaskCollection_TaskAppliesTo(t *testing.T) {
	commits := model.NewPACCommitCollection()
	commits.Add(&model.PACCommit{SHA: "abc123", Message: "PRJ-123: Fix"})

	taskSystems := []config.TaskSystemConfig{
		{
			Name: "jira",
			Regex: []config.RegexConfig{
				{Pattern: "/(PRJ-\\d+)/", Label: "jira"},
			},
		},
	}

	tasks := BuildTaskCollection(commits, taskSystems)

	// Find the PRJ-123 task
	task := tasks.FindOrCreate("PRJ-123")

	if !task.AppliesTo["jira"] {
		t.Error("expected task to apply to 'jira' system")
	}

	if !task.HasLabel("jira") {
		t.Error("expected task to have 'jira' label")
	}
}

func TestCreateTaskSystem(t *testing.T) {
	tests := []struct {
		name         string
		config       config.TaskSystemConfig
		expectedType string
	}{
		{
			name:         "none task system",
			config:       config.TaskSystemConfig{Name: "none"},
			expectedType: "*task.NoneTaskSystem",
		},
		{
			name:         "jira task system",
			config:       config.TaskSystemConfig{Name: "jira"},
			expectedType: "*task.JiraTaskSystem",
		},
		{
			name:         "unknown defaults to none",
			config:       config.TaskSystemConfig{Name: "custom"},
			expectedType: "*task.NoneTaskSystem",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts, err := CreateTaskSystem(tt.config)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Verify the system was created
			if ts == nil {
				t.Fatal("expected non-nil TaskSystem")
			}

			if ts.Name() != tt.config.Name {
				t.Errorf("Name() = %q, want %q", ts.Name(), tt.config.Name)
			}
		})
	}
}

func TestNoneTaskSystem_Apply(t *testing.T) {
	cfg := config.TaskSystemConfig{Name: "none"}
	ts := NewNoneTaskSystem(cfg)

	tasks := model.NewPACTaskCollection()
	task := tasks.FindOrCreate("TASK-1")
	task.AppliesTo["none"] = true

	// Should not return error
	err := ts.Apply(tasks)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
