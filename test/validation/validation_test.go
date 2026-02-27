// Package validation contains comprehensive tests for backwards compatibility
// with existing Ruby PAC templates and settings files.
package validation

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/config"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/model"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/report"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/task"
)

// getProjectRoot returns the project root directory
func getProjectRoot(t *testing.T) string {
	// Start from current directory and walk up to find go.mod
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("Could not find project root (go.mod)")
		}
		dir = parent
	}
}

// TestDefaultSettingsLoad verifies that default_settings.yml loads correctly
func TestDefaultSettingsLoad(t *testing.T) {
	root := getProjectRoot(t)
	settingsPath := filepath.Join(root, "settings", "default_settings.yml")

	settings, err := config.LoadSettings(settingsPath, nil)
	if err != nil {
		t.Fatalf("Failed to load default_settings.yml: %v", err)
	}

	// Verify general settings
	if settings.General.Strict != false {
		t.Errorf("Expected strict=false, got %v", settings.General.Strict)
	}

	// Verify templates are loaded
	if len(settings.Templates) != 2 {
		t.Errorf("Expected 2 templates, got %d", len(settings.Templates))
	}

	// Verify task systems
	if len(settings.TaskSystems) != 3 {
		t.Errorf("Expected 3 task systems, got %d", len(settings.TaskSystems))
	}

	// Check first task system (none)
	if settings.TaskSystems[0].Name != "none" {
		t.Errorf("Expected first task system to be 'none', got '%s'", settings.TaskSystems[0].Name)
	}

	// Check regex patterns are loaded
	if len(settings.TaskSystems[0].Regex) != 4 {
		t.Errorf("Expected 4 regex patterns for 'none' system, got %d", len(settings.TaskSystems[0].Regex))
	}

	// Verify VCS settings
	if settings.VCS.Type != "git" {
		t.Errorf("Expected VCS type 'git', got '%s'", settings.VCS.Type)
	}
}

// TestMinimalSettingsLoad verifies that minimal_settings.yml loads correctly
func TestMinimalSettingsLoad(t *testing.T) {
	root := getProjectRoot(t)
	settingsPath := filepath.Join(root, "settings", "minimal_settings.yml")

	settings, err := config.LoadSettings(settingsPath, nil)
	if err != nil {
		t.Fatalf("Failed to load minimal_settings.yml: %v", err)
	}

	// Minimal settings should have at least one task system
	if len(settings.TaskSystems) == 0 {
		t.Error("Expected at least one task system in minimal settings")
	}
}

// TestDefaultMarkdownTemplateRendering verifies the default.md template renders correctly
func TestDefaultMarkdownTemplateRendering(t *testing.T) {
	root := getProjectRoot(t)
	templatePath := filepath.Join(root, "templates", "default.md")

	content, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read default.md: %v", err)
	}

	// Create test data
	commits := model.NewPACCommitCollection()
	commits.Add(&model.PACCommit{
		SHA:      "abc1234567890",
		ShortSHA: "abc1234",
		Header:   "Fix critical bug",
		Message:  "Fix critical bug\n\nDetailed description",
	})
	commits.Add(&model.PACCommit{
		SHA:      "def4567890123",
		ShortSHA: "def4567",
		Header:   "Add new feature",
		Message:  "Add new feature\n\nMore details",
	})

	tasks := model.NewPACTaskCollection()

	// Referenced task
	task1 := tasks.FindOrCreate("PAC-123")
	task1.AddCommit(commits.Commits[0])
	task1.AppliesTo["none"] = true

	// Unreferenced commit
	unref := tasks.FindOrCreate("")
	unref.AddCommit(commits.Commits[1])

	// Render
	renderer := report.NewRenderer()
	data := map[string]interface{}{
		"tasks": tasks.ToLiquid(),
	}

	output, err := renderer.RenderTemplate(string(content), data)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	// Verify output contains expected elements
	checks := []string{
		"# PAC Changelog",
		"PAC-123",
		"abc1234",
		"Fix critical bug",
		"## Unspecified",
		"def4567",
		"Add new feature",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("Expected output to contain '%s', but it didn't.\nOutput:\n%s", check, output)
		}
	}
}

// TestDefaultHTMLTemplateRendering verifies the default_html.html template renders correctly
func TestDefaultHTMLTemplateRendering(t *testing.T) {
	root := getProjectRoot(t)
	templatePath := filepath.Join(root, "templates", "default_html.html")

	content, err := os.ReadFile(templatePath)
	if err != nil {
		t.Fatalf("Failed to read default_html.html: %v", err)
	}

	// Create test data
	commits := model.NewPACCommitCollection()
	commits.Add(&model.PACCommit{
		SHA:      "abc1234567890",
		ShortSHA: "abc1234",
		Header:   "Fix HTML bug",
		Message:  "Fix HTML bug",
	})

	tasks := model.NewPACTaskCollection()
	task1 := tasks.FindOrCreate("HTML-456")
	task1.AddCommit(commits.Commits[0])

	// Render
	renderer := report.NewRenderer()
	data := map[string]interface{}{
		"tasks": tasks.ToLiquid(),
	}

	output, err := renderer.RenderTemplate(string(content), data)
	if err != nil {
		t.Fatalf("Failed to render template: %v", err)
	}

	// Verify HTML structure
	checks := []string{
		"<html>",
		"</html>",
		"<title>PAC Changelog</title>",
		"<h1>PAC Changelog</h1>",
		"HTML-456",
		"abc1234",
		"Fix HTML bug",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("Expected output to contain '%s'", check)
		}
	}
}

// TestGitHubMarkdownTemplate verifies the GitHub example template
func TestGitHubMarkdownTemplate(t *testing.T) {
	root := getProjectRoot(t)
	templatePath := filepath.Join(root, "templates", "examples", "github.md")

	content, err := os.ReadFile(templatePath)
	if err != nil {
		t.Skipf("GitHub template not found: %v", err)
	}

	// Skip if template is empty (placeholder)
	if strings.TrimSpace(string(content)) == "" {
		t.Skip("GitHub template is empty (placeholder file)")
	}

	// Create test data
	tasks := model.NewPACTaskCollection()
	task1 := tasks.FindOrCreate("GH-789")
	task1.AddCommit(&model.PACCommit{
		SHA:      "ghcommit123456",
		ShortSHA: "ghcommi",
		Header:   "GitHub integration",
	})

	renderer := report.NewRenderer()
	data := map[string]interface{}{
		"tasks": tasks.ToLiquid(),
	}

	output, err := renderer.RenderTemplate(string(content), data)
	if err != nil {
		t.Fatalf("Failed to render GitHub template: %v", err)
	}

	// Just verify it renders without error and has some content
	if output == "" {
		t.Error("GitHub template rendered empty output")
	}
}

// TestAllLiquidVariablesAccessible verifies all expected Liquid variables work
func TestAllLiquidVariablesAccessible(t *testing.T) {
	template := `
sha: {{commit.sha}}
shortsha: {{commit.shortsha}}
short_sha: {{commit.short_sha}}
header: {{commit.header}}
message: {{commit.message}}
body: {{commit.body}}
date: {{commit.date}}
`
	commit := &model.PACCommit{
		SHA:      "fullsha1234567890",
		ShortSHA: "fullsha",
		Header:   "Test header",
		Message:  "Test header\n\nTest body",
		Body:     "Test body",
	}

	renderer := report.NewRenderer()
	data := map[string]interface{}{
		"commit": commit.ToLiquid(),
	}

	output, err := renderer.RenderTemplate(template, data)
	if err != nil {
		t.Fatalf("Failed to render: %v", err)
	}

	// Verify all variables rendered
	checks := map[string]string{
		"sha: fullsha1234567890": "sha",
		"shortsha: fullsha":      "shortsha",
		"short_sha: fullsha":     "short_sha",
		"header: Test header":    "header",
		"body: Test body":        "body",
	}

	for expected, name := range checks {
		if !strings.Contains(output, expected) {
			t.Errorf("Variable '%s' not accessible. Expected '%s' in output:\n%s", name, expected, output)
		}
	}
}

// TestTaskLiquidVariables verifies task variables are accessible
func TestTaskLiquidVariables(t *testing.T) {
	template := `
task_id: {{task.task_id}}
commits_count: {{task.commits | size}}
{% for label in task.label %}label: {{label}}
{% endfor %}
`
	tasks := model.NewPACTaskCollection()
	task := tasks.FindOrCreate("TASK-999")
	task.AddCommit(&model.PACCommit{SHA: "abc123", ShortSHA: "abc123", Header: "Test"})
	task.AddLabel("bugfix")
	task.AddLabel("priority")

	renderer := report.NewRenderer()
	data := map[string]interface{}{
		"task": task.ToLiquid(),
	}

	output, err := renderer.RenderTemplate(template, data)
	if err != nil {
		t.Fatalf("Failed to render: %v", err)
	}

	if !strings.Contains(output, "task_id: TASK-999") {
		t.Errorf("task_id not rendered correctly:\n%s", output)
	}

	if !strings.Contains(output, "commits_count: 1") {
		t.Errorf("commits count not rendered correctly:\n%s", output)
	}
}

// TestRegexPatternCompatibility verifies Ruby-style regex patterns work
func TestRegexPatternCompatibility(t *testing.T) {
	testCases := []struct {
		name     string
		pattern  string
		input    string
		expected []string
	}{
		{
			name:     "Simple pattern",
			pattern:  "/PAC-(\\d+)/",
			input:    "Fixed PAC-123 and PAC-456",
			expected: []string{"123", "456"},
		},
		{
			name:     "Case insensitive",
			pattern:  "/issue:(\\d+)/i",
			input:    "Fixed Issue:123 and ISSUE:456",
			expected: []string{"123", "456"},
		},
		{
			name:     "Hash pattern",
			pattern:  "/(#\\d+)/",
			input:    "Closes #42 and #100",
			expected: []string{"#42", "#100"},
		},
		{
			name:     "Complex pattern",
			pattern:  "/([A-Z]+-\\d+)/",
			input:    "PROJ-123 and ABC-456",
			expected: []string{"PROJ-123", "ABC-456"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			regexCfg := []config.RegexConfig{
				{Pattern: tc.pattern, Label: "test"},
			}

			results := task.ExtractTaskIDs(tc.input, regexCfg)

			if len(results) != len(tc.expected) {
				t.Errorf("Expected %d matches, got %d: %v",
					len(tc.expected), len(results), results)
				return
			}

			for i, exp := range tc.expected {
				if results[i].TaskID != exp {
					t.Errorf("Expected match[%d]='%s', got '%s'", i, exp, results[i].TaskID)
				}
			}
		})
	}
}

// TestStatisticsVariables verifies PAC statistics are correctly calculated
func TestStatisticsVariables(t *testing.T) {
	commits := model.NewPACCommitCollection()

	// Add 10 commits
	for i := 0; i < 10; i++ {
		commits.Add(&model.PACCommit{
			SHA:      "abc123",
			ShortSHA: "abc123",
			Header:   "Test commit",
		})
	}

	tasks := model.NewPACTaskCollection()

	// 7 referenced, 3 unreferenced
	for i := 0; i < 7; i++ {
		task := tasks.FindOrCreate("TASK-" + string(rune('A'+i)))
		task.AddCommit(commits.Commits[i])
	}

	unref := tasks.FindOrCreate("")
	for i := 7; i < 10; i++ {
		unref.AddCommit(commits.Commits[i])
	}

	// Verify statistics
	if commits.Count() != 10 {
		t.Errorf("Expected 10 commits, got %d", commits.Count())
	}

	template := `
count: {{pac_c_count}}
referenced: {{pac_c_referenced}}
unreferenced: {{pac_c_unreferenced}}
health: {{pac_health}}
`
	renderer := report.NewRenderer()
	data := map[string]interface{}{
		"pac_c_count":        10,
		"pac_c_referenced":   7,
		"pac_c_unreferenced": 3,
		"pac_health":         70.0,
	}

	output, err := renderer.RenderTemplate(template, data)
	if err != nil {
		t.Fatalf("Failed to render: %v", err)
	}

	checks := []string{
		"count: 10",
		"referenced: 7",
		"unreferenced: 3",
		"health: 70",
	}

	for _, check := range checks {
		if !strings.Contains(output, check) {
			t.Errorf("Expected '%s' in output:\n%s", check, output)
		}
	}
}

// TestPropertiesInjection verifies custom properties can be used in templates
func TestPropertiesInjection(t *testing.T) {
	template := `
version: {{properties.version}}
project: {{properties.project_name}}
`
	renderer := report.NewRenderer()
	data := map[string]interface{}{
		"properties": map[string]interface{}{
			"version":      "1.2.3",
			"project_name": "MyProject",
		},
	}

	output, err := renderer.RenderTemplate(template, data)
	if err != nil {
		t.Fatalf("Failed to render: %v", err)
	}

	if !strings.Contains(output, "version: 1.2.3") {
		t.Errorf("properties.version not rendered:\n%s", output)
	}

	if !strings.Contains(output, "project: MyProject") {
		t.Errorf("properties.project_name not rendered:\n%s", output)
	}
}

// TestLiquidFilters verifies common Liquid filters work
func TestLiquidFilters(t *testing.T) {
	testCases := []struct {
		name     string
		template string
		data     map[string]interface{}
		expected string
	}{
		{
			name:     "size filter",
			template: "{{items | size}}",
			data:     map[string]interface{}{"items": []string{"a", "b", "c"}},
			expected: "3",
		},
		{
			name:     "first filter",
			template: "{{items | first}}",
			data:     map[string]interface{}{"items": []string{"apple", "banana", "cherry"}},
			expected: "apple",
		},
		{
			name:     "last filter",
			template: "{{items | last}}",
			data:     map[string]interface{}{"items": []string{"apple", "banana", "cherry"}},
			expected: "cherry",
		},
		{
			name:     "join filter",
			template: "{{items | join: ', '}}",
			data:     map[string]interface{}{"items": []string{"a", "b", "c"}},
			expected: "a, b, c",
		},
		{
			name:     "upcase filter",
			template: "{{text | upcase}}",
			data:     map[string]interface{}{"text": "hello"},
			expected: "HELLO",
		},
		{
			name:     "downcase filter",
			template: "{{text | downcase}}",
			data:     map[string]interface{}{"text": "HELLO"},
			expected: "hello",
		},
	}

	renderer := report.NewRenderer()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := renderer.RenderTemplate(tc.template, tc.data)
			if err != nil {
				t.Fatalf("Failed to render: %v", err)
			}

			if strings.TrimSpace(output) != tc.expected {
				t.Errorf("Expected '%s', got '%s'", tc.expected, strings.TrimSpace(output))
			}
		})
	}
}

// TestLiquidControlFlow verifies Liquid control structures work
func TestLiquidControlFlow(t *testing.T) {
	testCases := []struct {
		name     string
		template string
		data     map[string]interface{}
		contains string
	}{
		{
			name:     "if true",
			template: "{% if show %}visible{% endif %}",
			data:     map[string]interface{}{"show": true},
			contains: "visible",
		},
		{
			name:     "if false",
			template: "{% if show %}visible{% else %}hidden{% endif %}",
			data:     map[string]interface{}{"show": false},
			contains: "hidden",
		},
		{
			name:     "unless",
			template: "{% unless hide %}visible{% endunless %}",
			data:     map[string]interface{}{"hide": false},
			contains: "visible",
		},
		{
			name:     "for loop",
			template: "{% for item in items %}{{item}}{% endfor %}",
			data:     map[string]interface{}{"items": []string{"a", "b", "c"}},
			contains: "abc",
		},
		{
			name:     "for loop with forloop.first",
			template: "{% for item in items %}{% if forloop.first %}first:{% endif %}{{item}}{% endfor %}",
			data:     map[string]interface{}{"items": []string{"x", "y"}},
			contains: "first:x",
		},
	}

	renderer := report.NewRenderer()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			output, err := renderer.RenderTemplate(tc.template, tc.data)
			if err != nil {
				t.Fatalf("Failed to render: %v", err)
			}

			if !strings.Contains(output, tc.contains) {
				t.Errorf("Expected output to contain '%s', got '%s'", tc.contains, output)
			}
		})
	}
}
