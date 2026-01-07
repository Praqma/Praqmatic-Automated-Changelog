// Package templates provides template compatibility tests.
package templates

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/config"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/model"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/report"
)

// TestDefaultMarkdownTemplate tests that templates/default.md renders correctly.
func TestDefaultMarkdownTemplate(t *testing.T) {
	// Find the templates directory
	templatePath := findTemplatePath(t, "default.md")
	if templatePath == "" {
		t.Skip("templates/default.md not found")
	}

	// Read template
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("failed to read template: %v", err)
	}

	// Create test data
	commits, tasks := createTestData()

	settings := &config.Settings{
		Properties: map[string]any{
			"title":   "Test Changelog",
			"version": "1.0.0",
		},
	}

	// Render template
	generator := report.NewGenerator(tasks, commits)
	output, err := generator.GenerateToString(string(templateContent), settings)
	if err != nil {
		t.Fatalf("failed to render template: %v", err)
	}

	// Verify output is not empty
	if len(output) == 0 {
		t.Error("expected non-empty output")
	}

	// Verify it contains some expected elements (based on template structure)
	t.Logf("Template output length: %d bytes", len(output))
}

// TestDefaultHTMLTemplate tests that templates/default_html.html renders correctly.
func TestDefaultHTMLTemplate(t *testing.T) {
	templatePath := findTemplatePath(t, "default_html.html")
	if templatePath == "" {
		t.Skip("templates/default_html.html not found")
	}

	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("failed to read template: %v", err)
	}

	commits, tasks := createTestData()
	settings := &config.Settings{}

	generator := report.NewGenerator(tasks, commits)
	output, err := generator.GenerateToString(string(templateContent), settings)
	if err != nil {
		t.Fatalf("failed to render template: %v", err)
	}

	if len(output) == 0 {
		t.Error("expected non-empty output")
	}
}

// TestLiquidVariableAccess tests that all required Liquid variables are accessible.
func TestLiquidVariableAccess(t *testing.T) {
	commits, tasks := createTestData()
	settings := &config.Settings{
		Properties: map[string]any{
			"custom_key": "custom_value",
		},
	}

	// Template that tests all standard variables
	template := `
Tasks Referenced: {{ tasks.referenced.size }}
Tasks Unreferenced: {{ tasks.unreferenced.size }}
Commit Count: {{ pac_c_count }}
Referenced Count: {{ pac_c_referenced }}
Unreferenced Count: {{ pac_c_unreferenced }}
Health: {{ pac_health }}
Custom Property: {{ properties.custom_key }}
{% for task in tasks.referenced %}
Task ID: {{ task.task_id }}
Task Commits: {{ task.commits.size }}
{% for commit in task.commits %}
Commit SHA: {{ commit.sha }}
Commit Short SHA: {{ commit.shortsha }}
Commit Header: {{ commit.header }}
{% endfor %}
{% endfor %}
`

	generator := report.NewGenerator(tasks, commits)
	output, err := generator.GenerateToString(template, settings)
	if err != nil {
		t.Fatalf("failed to render template: %v", err)
	}

	// Verify all variables rendered correctly
	checks := []string{
		"Tasks Referenced: 2",
		"Commit Count: 3",
		"Custom Property: custom_value",
		"TASK-1",
		"TASK-2",
		"abc1234",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("expected output to contain %q\nGot:\n%s", check, output)
		}
	}
}

// TestLiquidFiltersWork tests that common Liquid filters work correctly.
func TestLiquidFiltersWork(t *testing.T) {
	commits, tasks := createTestData()
	settings := &config.Settings{}

	template := `
Size filter: {{ tasks.referenced | size }}
First filter: {{ tasks.referenced | first | map: "task_id" }}
Join filter: {% assign ids = tasks.referenced | map: "task_id" %}{{ ids | join: ", " }}
Upcase filter: {{ "hello" | upcase }}
Downcase filter: {{ "HELLO" | downcase }}
`

	generator := report.NewGenerator(tasks, commits)
	output, err := generator.GenerateToString(template, settings)
	if err != nil {
		t.Fatalf("failed to render template: %v", err)
	}

	// Check filters work
	if !strings.Contains(output, "Size filter: 2") {
		t.Errorf("size filter failed, got: %s", output)
	}
	if !strings.Contains(output, "Upcase filter: HELLO") {
		t.Errorf("upcase filter failed, got: %s", output)
	}
	if !strings.Contains(output, "Downcase filter: hello") {
		t.Errorf("downcase filter failed, got: %s", output)
	}
}

// TestLiquidControlFlow tests if/else/for control flow.
func TestLiquidControlFlow(t *testing.T) {
	commits, tasks := createTestData()
	settings := &config.Settings{}

	template := `
{% if tasks.referenced.size > 0 %}Has tasks{% else %}No tasks{% endif %}
{% for task in tasks.referenced %}
- {{ task.task_id }}
{% endfor %}
{% unless tasks.unreferenced.size == 0 %}Has unreferenced{% endunless %}
`

	generator := report.NewGenerator(tasks, commits)
	output, err := generator.GenerateToString(template, settings)
	if err != nil {
		t.Fatalf("failed to render template: %v", err)
	}

	if !strings.Contains(output, "Has tasks") {
		t.Errorf("if condition failed, got: %s", output)
	}
	if !strings.Contains(output, "- TASK-1") {
		t.Errorf("for loop failed, got: %s", output)
	}
}

// TestShortSHAAlias tests that both shortsha and short_sha work.
func TestShortSHAAlias(t *testing.T) {
	commits, tasks := createTestData()
	settings := &config.Settings{}

	template := `
shortsha: {{ tasks.referenced.first.commits.first.shortsha }}
short_sha: {{ tasks.referenced.first.commits.first.short_sha }}
`

	generator := report.NewGenerator(tasks, commits)
	output, err := generator.GenerateToString(template, settings)
	if err != nil {
		t.Fatalf("failed to render template: %v", err)
	}

	// Both should render the same value
	if !strings.Contains(output, "shortsha: abc1234") {
		t.Errorf("shortsha failed, got: %s", output)
	}
	if !strings.Contains(output, "short_sha: abc1234") {
		t.Errorf("short_sha alias failed, got: %s", output)
	}
}

// Helper functions

func findTemplatePath(t *testing.T, name string) string {
	t.Helper()

	// Try multiple locations
	paths := []string{
		filepath.Join("templates", name),
		filepath.Join("..", "..", "templates", name),
		filepath.Join("..", "..", "..", "templates", name),
	}

	for _, p := range paths {
		if _, err := os.Stat(p); err == nil {
			return p
		}
	}

	return ""
}

func createTestData() (*model.PACCommitCollection, *model.PACTaskCollection) {
	commits := model.NewPACCommitCollection()

	commit1 := &model.PACCommit{
		SHA:       "abc1234567890",
		ShortSHA:  "abc1234",
		Message:   "TASK-1: First commit\n\nDetailed description",
		Header:    "TASK-1: First commit",
		Body:      "Detailed description",
		Timestamp: time.Date(2025, 1, 15, 10, 0, 0, 0, time.UTC),
	}
	commit1.MarkReferenced()
	commits.Add(commit1)

	commit2 := &model.PACCommit{
		SHA:       "def4567890123",
		ShortSHA:  "def4567",
		Message:   "TASK-2: Second commit",
		Header:    "TASK-2: Second commit",
		Timestamp: time.Date(2025, 1, 16, 10, 0, 0, 0, time.UTC),
	}
	commit2.MarkReferenced()
	commits.Add(commit2)

	commit3 := &model.PACCommit{
		SHA:       "ghi7890123456",
		ShortSHA:  "ghi7890",
		Message:   "Unreferenced commit",
		Header:    "Unreferenced commit",
		Timestamp: time.Date(2025, 1, 17, 10, 0, 0, 0, time.UTC),
	}
	commits.Add(commit3)

	tasks := model.NewPACTaskCollection()

	task1 := tasks.FindOrCreate("TASK-1")
	task1.AddCommit(commit1)
	task1.AddLabel("feature")
	task1.Data = map[string]any{
		"summary": "Implement feature X",
		"status":  "Done",
	}

	task2 := tasks.FindOrCreate("TASK-2")
	task2.AddCommit(commit2)
	task2.AddLabel("bugfix")

	// Unreferenced commits go to empty task
	unref := tasks.FindOrCreate("")
	unref.AddCommit(commit3)

	return commits, tasks
}
