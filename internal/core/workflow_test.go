package core

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/config"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/model"
)

func TestRun_WithRealRepo(t *testing.T) {
	// Create a temporary directory for the test repo
	tmpDir, err := os.MkdirTemp("", "pac-core-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Initialize git repo
	setupTestRepo(t, tmpDir)

	// Create test settings
	settings := &config.Settings{
		VCS: config.VCSConfig{
			Type:         "git",
			RepoLocation: tmpDir,
		},
		TaskSystems: []config.TaskSystemConfig{
			{
				Name: "none",
				Regex: []config.RegexConfig{
					{Pattern: `/ISSUE-(\d+)/i`},
				},
			},
		},
		Templates: []config.TemplateConfig{
			{
				Location: filepath.Join(tmpDir, "test_template.md"),
				Output:   filepath.Join(tmpDir, "output.md"),
			},
		},
		Verbosity: 0,
	}

	// Create a simple template
	templateContent := `# Changelog
Tasks: {{ tasks.referenced | size }}
`
	if err := os.WriteFile(settings.Templates[0].Location, []byte(templateContent), 0644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	// Run the workflow
	result, err := Run(settings, "v1.0.0", "HEAD")
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Verify results
	if result.Commits.Count() != 2 {
		t.Errorf("expected 2 commits, got %d", result.Commits.Count())
	}

	if len(result.Tasks.Referenced()) < 1 {
		t.Errorf("expected at least 1 referenced task, got %d", len(result.Tasks.Referenced()))
	}

	// Check output file was created
	if _, err := os.Stat(settings.Templates[0].Output); os.IsNotExist(err) {
		t.Error("output file was not created")
	}
}

func TestRunFromLatestTag_WithRealRepo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pac-core-tag-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	setupTestRepo(t, tmpDir)

	settings := &config.Settings{
		VCS: config.VCSConfig{
			Type:         "git",
			RepoLocation: tmpDir,
		},
		TaskSystems: []config.TaskSystemConfig{
			{
				Name: "none",
				Regex: []config.RegexConfig{
					{Pattern: `/ISSUE-(\d+)/i`},
				},
			},
		},
		Templates: []config.TemplateConfig{
			{
				Location: filepath.Join(tmpDir, "test_template.md"),
				Output:   filepath.Join(tmpDir, "output.md"),
			},
		},
		Verbosity: 0,
	}

	// Create template
	if err := os.WriteFile(settings.Templates[0].Location, []byte("# Test"), 0644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	// Run from latest tag
	result, err := RunFromLatestTag(settings, "v*", "HEAD")
	if err != nil {
		t.Fatalf("RunFromLatestTag failed: %v", err)
	}

	if result.Commits.Count() != 2 {
		t.Errorf("expected 2 commits, got %d", result.Commits.Count())
	}
}

func TestRun_StrictMode(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "pac-strict-test")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	setupTestRepo(t, tmpDir)

	settings := &config.Settings{
		General: config.GeneralSettings{
			Strict: true,
		},
		VCS: config.VCSConfig{
			Type:         "git",
			RepoLocation: tmpDir,
		},
		TaskSystems: []config.TaskSystemConfig{
			{
				Name: "jira", // Jira will fail because no URL configured
				Regex: []config.RegexConfig{
					{Pattern: `/ISSUE-(\d+)/i`},
				},
			},
		},
		Templates: []config.TemplateConfig{
			{
				Location: filepath.Join(tmpDir, "test_template.md"),
			},
		},
		Verbosity: 0,
	}

	// Create template
	if err := os.WriteFile(settings.Templates[0].Location, []byte("# Test"), 0644); err != nil {
		t.Fatalf("failed to write template: %v", err)
	}

	// Run should fail in strict mode when task system fails
	_, err = Run(settings, "v1.0.0", "HEAD")
	if err == nil {
		t.Error("expected error in strict mode, got nil")
	}
}

func TestApplyTaskSystems_AllSuccess(t *testing.T) {
	settings := &config.Settings{
		TaskSystems: []config.TaskSystemConfig{
			{
				Name: "none",
				Regex: []config.RegexConfig{
					{Pattern: `/TASK-(\d+)/i`},
				},
			},
		},
	}

	tasks := createTestTaskCollection()

	allOK := applyTaskSystems(settings, tasks)
	if !allOK {
		t.Error("expected allOK to be true")
	}
}

func TestDefaultStr(t *testing.T) {
	tests := []struct {
		s, def, expected string
	}{
		{"value", "default", "value"},
		{"", "default", "default"},
		{"", "", ""},
	}

	for _, tt := range tests {
		result := defaultStr(tt.s, tt.def)
		if result != tt.expected {
			t.Errorf("defaultStr(%q, %q) = %q, want %q", tt.s, tt.def, result, tt.expected)
		}
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
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write file %s: %v", path, err)
	}
}

func createTestTaskCollection() *model.PACTaskCollection {
	tasks := model.NewPACTaskCollection()
	task := tasks.FindOrCreate("TEST-1")
	task.AddCommit(&model.PACCommit{
		SHA:      "abc123",
		ShortSHA: "abc123",
		Message:  "TEST-1: Test commit",
		Header:   "TEST-1: Test commit",
	})
	return tasks
}
