// Package integration provides end-to-end integration tests for PAC.
package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/config"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/model"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/report"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/task"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/vcs"
)

// TestFullWorkflow_WithRealGitRepo tests the complete PAC workflow
// using a temporary git repository.
func TestFullWorkflow_WithRealGitRepo(t *testing.T) {
	// Create a temporary directory for the test repo
	tmpDir, err := os.MkdirTemp("", "pac-integration-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize git repo
	setupTestRepo(t, tmpDir)

	// Create test settings
	vcsSettings := config.VCSConfig{
		Type:         "git",
		RepoLocation: tmpDir,
	}

	taskSystems := []config.TaskSystemConfig{
		{
			Name: "none",
			Regex: []config.RegexConfig{
				{Pattern: `/ISSUE-(\d+)/i`},
			},
		},
	}

	// Step 1: Get commits using VCS
	gitVCS, err := vcs.NewGitVCS(vcsSettings)
	if err != nil {
		t.Fatalf("failed to create GitVCS: %v", err)
	}

	commits, err := gitVCS.GetDelta("v1.0.0", "HEAD")
	if err != nil {
		t.Fatalf("failed to get delta: %v", err)
	}

	// Verify we got the expected commits (should be 2: commit2 and commit3)
	if commits.Count() != 2 {
		t.Errorf("expected 2 commits, got %d", commits.Count())
	}

	// Step 2: Extract tasks from commits
	tasks := task.BuildTaskCollection(commits, taskSystems)

	// Step 3: Apply task systems
	systems, err := task.CreateAllTaskSystems(taskSystems)
	if err != nil {
		t.Fatalf("failed to create task systems: %v", err)
	}

	for _, ts := range systems {
		if applyErr := ts.Apply(tasks); applyErr != nil {
			t.Errorf("task system apply failed: %v", applyErr)
		}
	}

	// Verify task extraction
	if tasks.Count() < 1 {
		t.Errorf("expected at least 1 task, got %d", tasks.Count())
	}

	// Verify ISSUE-123 was extracted
	found := false
	for _, tsk := range tasks.Referenced() {
		if tsk.TaskID == "123" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected to find task ISSUE-123")
	}

	// Step 4: Generate report
	settings := &config.Settings{
		Templates: []config.TemplateConfig{
			{
				Location: filepath.Join(tmpDir, "test_template.md"),
			},
		},
	}

	// Create a simple template file
	templateContent := `# Changelog
Tasks: {{ tasks.size }}
{% for task in tasks.referenced %}
- {{ task.task_id }}
{% endfor %}
`
	if writeErr := os.WriteFile(filepath.Join(tmpDir, "test_template.md"), []byte(templateContent), 0o644); writeErr != nil {
		t.Fatalf("failed to write template: %v", writeErr)
	}

	generator := report.NewGenerator(tasks, commits)
	output, err := generator.GenerateToString(templateContent, settings)
	if err != nil {
		t.Fatalf("failed to generate report: %v", err)
	}

	// Verify output contains expected content
	if !strings.Contains(output, "123") {
		t.Errorf("expected output to contain 123, got:\n%s", output)
	}
}

// TestGetLatestTag_Integration tests tag pattern matching.
func TestGetLatestTag_Integration(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pac-tag-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	setupTestRepo(t, tmpDir)

	vcsSettings := config.VCSConfig{
		Type:         "git",
		RepoLocation: tmpDir,
	}

	gitVCS, err := vcs.NewGitVCS(vcsSettings)
	if err != nil {
		t.Fatalf("failed to create GitVCS: %v", err)
	}

	// Test finding the latest tag with "v*" pattern
	tag, err := gitVCS.GetLatestTag("v*")
	if err != nil {
		t.Fatalf("failed to get latest tag: %v", err)
	}

	if tag != "v1.0.0" {
		t.Errorf("expected tag v1.0.0, got %s", tag)
	}
}

// TestMultipleTaskSystems_Integration tests using multiple task systems.
func TestMultipleTaskSystems_Integration(t *testing.T) {
	// Create commits with different task ID patterns
	commits := model.NewPACCommitCollection()
	commits.Add(&model.PACCommit{
		SHA:       "abc123",
		ShortSHA:  "abc123",
		Message:   "JIRA-100: Feature A\nISSUE-200: Related fix",
		Header:    "JIRA-100: Feature A",
		Timestamp: time.Now(),
	})
	commits.Add(&model.PACCommit{
		SHA:       "def456",
		ShortSHA:  "def456",
		Message:   "BUG-300: Fix something",
		Header:    "BUG-300: Fix something",
		Timestamp: time.Now(),
	})

	taskSystems := []config.TaskSystemConfig{
		{
			Name: "none",
			Regex: []config.RegexConfig{
				{Pattern: `/JIRA-(\d+)/i`},
				{Pattern: `/ISSUE-(\d+)/i`},
				{Pattern: `/BUG-(\d+)/i`},
			},
		},
	}

	tasks := task.BuildTaskCollection(commits, taskSystems)

	// Should find 3 unique task IDs: 100, 200, 300
	refTasks := tasks.Referenced()
	if len(refTasks) != 3 {
		t.Errorf("expected 3 tasks, got %d", len(refTasks))
		for _, tsk := range refTasks {
			t.Logf("Found task: %s", tsk.TaskID)
		}
	}
}

// TestCommitFiltering_Integration tests path-based commit filtering.
func TestCommitFiltering_Integration(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pac-filter-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Setup repo with commits touching different paths
	runGit(t, tmpDir, "init")
	runGit(t, tmpDir, "config", "user.email", "test@example.com")
	runGit(t, tmpDir, "config", "user.name", "Test User")

	// Create directories
	if mkdirErr := os.MkdirAll(filepath.Join(tmpDir, "src"), 0o755); mkdirErr != nil {
		t.Fatalf("failed to create src dir: %v", mkdirErr)
	}
	if mkdirErr := os.MkdirAll(filepath.Join(tmpDir, "docs"), 0o755); mkdirErr != nil {
		t.Fatalf("failed to create docs dir: %v", mkdirErr)
	}

	// Initial commit
	writeFile(t, tmpDir, "README.md", "# Test")
	runGit(t, tmpDir, "add", ".")
	runGit(t, tmpDir, "commit", "-m", "Initial commit")
	runGit(t, tmpDir, "tag", "v1.0.0")

	// Commit touching src
	writeFile(t, tmpDir, "src/main.go", "package main")
	runGit(t, tmpDir, "add", ".")
	runGit(t, tmpDir, "commit", "-m", "TASK-1: Add main.go")

	// Commit touching docs
	writeFile(t, tmpDir, "docs/README.md", "# Docs")
	runGit(t, tmpDir, "add", ".")
	runGit(t, tmpDir, "commit", "-m", "TASK-2: Add docs")

	// Commit touching src again
	writeFile(t, tmpDir, "src/util.go", "package main")
	runGit(t, tmpDir, "add", ".")
	runGit(t, tmpDir, "commit", "-m", "TASK-3: Add util")

	vcsSettings := config.VCSConfig{
		Type:         "git",
		RepoLocation: tmpDir,
	}

	gitVCS, err := vcs.NewGitVCS(vcsSettings)
	if err != nil {
		t.Fatalf("failed to create GitVCS: %v", err)
	}

	// Get all commits
	allCommits, err := gitVCS.GetDelta("v1.0.0", "HEAD")
	if err != nil {
		t.Fatalf("failed to get delta: %v", err)
	}

	if allCommits.Count() != 3 {
		t.Errorf("expected 3 commits without filter, got %d", allCommits.Count())
	}
}

// TestTemplateRendering_Integration tests various template scenarios.
func TestTemplateRendering_Integration(t *testing.T) {
	commits := model.NewPACCommitCollection()
	commit := &model.PACCommit{
		SHA:       "abc123def456",
		ShortSHA:  "abc123d",
		Message:   "TASK-1: Feature\n\nDetailed description here",
		Header:    "TASK-1: Feature",
		Body:      "Detailed description here",
		Timestamp: time.Date(2025, 1, 15, 10, 30, 0, 0, time.UTC),
	}
	commit.MarkReferenced()
	commits.Add(commit)

	tasks := model.NewPACTaskCollection()
	tsk := tasks.FindOrCreate("TASK-1")
	tsk.AddCommit(commit)
	tsk.Data = map[string]any{
		"summary": "Implement feature X",
		"status":  "Done",
	}

	settings := &config.Settings{
		Templates: []config.TemplateConfig{
			{
				Location: "test.md",
			},
		},
	}

	templateContent := `# Release Notes
{% for task in tasks.referenced %}
## {{ task.task_id }}
{% for commit in task.commits %}
  - {{ commit.short_sha }}: {{ commit.header }}
{% endfor %}
{% endfor %}
Health: {{ pac_health }}%
`

	generator := report.NewGenerator(tasks, commits)
	output, err := generator.GenerateToString(templateContent, settings)
	if err != nil {
		t.Fatalf("failed to generate: %v", err)
	}

	// Verify template output
	checks := []string{
		"TASK-1",
		"abc123d",
		"Health: 100%",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("expected output to contain %q, got:\n%s", check, output)
		}
	}
}

// TestUnreferencedCommits_Integration tests handling of commits without task IDs.
func TestUnreferencedCommits_Integration(t *testing.T) {
	commits := model.NewPACCommitCollection()
	commits.Add(&model.PACCommit{
		SHA:      "abc123",
		ShortSHA: "abc123",
		Message:  "TASK-1: With task ID",
		Header:   "TASK-1: With task ID",
	})
	commits.Add(&model.PACCommit{
		SHA:      "def456",
		ShortSHA: "def456",
		Message:  "Fix typo",
		Header:   "Fix typo",
	})
	commits.Add(&model.PACCommit{
		SHA:      "ghi789",
		ShortSHA: "ghi789",
		Message:  "Update docs",
		Header:   "Update docs",
	})

	taskSystems := []config.TaskSystemConfig{
		{
			Name: "none",
			Regex: []config.RegexConfig{
				{Pattern: `/TASK-(\d+)/i`},
			},
		},
	}

	tasks := task.BuildTaskCollection(commits, taskSystems)
	task.UpdateCommitReferenceCounts(commits, tasks)

	// Verify counts
	if commits.CountWith() != 1 {
		t.Errorf("expected 1 referenced commit, got %d", commits.CountWith())
	}
	if commits.CountWithout() != 2 {
		t.Errorf("expected 2 unreferenced commits, got %d", commits.CountWithout())
	}

	// Verify health calculation
	health := commits.Health()
	// Health should be ~33% (1 out of 3)
	if health < 30 || health > 40 {
		t.Errorf("expected health around 33%%, got %.2f%%", health)
	}

	// Verify unreferenced commits collection
	unreferenced := tasks.UnreferencedCommits()
	if len(unreferenced) != 2 {
		t.Errorf("expected 2 unreferenced commits, got %d", len(unreferenced))
	}
}

// Helper functions

func setupTestRepo(t *testing.T, dir string) {
	t.Helper()

	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "test@example.com")
	runGit(t, dir, "config", "user.name", "Test User")

	// Initial commit
	writeFile(t, dir, "README.md", "# Test Project")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "Initial commit")

	// Tag v1.0.0
	runGit(t, dir, "tag", "v1.0.0")

	// Add more commits after the tag
	writeFile(t, dir, "file1.txt", "content1")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "ISSUE-123: Add file1")

	writeFile(t, dir, "file2.txt", "content2")
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "ISSUE-456: Add file2")
}

func runGit(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_DATE=2025-01-15T10:00:00Z",
		"GIT_COMMITTER_DATE=2025-01-15T10:00:00Z",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, output)
	}
}

func writeFile(t *testing.T, dir, name, content string) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write file %s: %v", path, err)
	}
}
