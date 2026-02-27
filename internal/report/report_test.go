package report

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/config"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/model"
)

func TestRenderer_RenderTemplate(t *testing.T) {
	renderer := NewRenderer()

	template := "Hello, {{ name }}!"
	data := map[string]any{"name": "World"}

	result, err := renderer.RenderTemplate(template, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "Hello, World!" {
		t.Errorf("got %q, want %q", result, "Hello, World!")
	}
}

func TestRenderer_RenderTemplate_ForLoop(t *testing.T) {
	renderer := NewRenderer()

	template := "{% for item in items %}{{ item }} {% endfor %}"
	data := map[string]any{"items": []string{"a", "b", "c"}}

	result, err := renderer.RenderTemplate(template, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "a b c " {
		t.Errorf("got %q, want %q", result, "a b c ")
	}
}

func TestRenderer_RenderTemplate_NestedData(t *testing.T) {
	renderer := NewRenderer()

	template := "{% for task in tasks.referenced %}{{ task.task_id }} {% endfor %}"
	data := map[string]any{
		"tasks": map[string]any{
			"referenced": []map[string]any{
				{"task_id": "TASK-1"},
				{"task_id": "TASK-2"},
			},
		},
	}

	result, err := renderer.RenderTemplate(template, data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result != "TASK-1 TASK-2 " {
		t.Errorf("got %q, want %q", result, "TASK-1 TASK-2 ")
	}
}

func TestGenerator_GenerateToString(t *testing.T) {
	// Create test data
	commits := model.NewPACCommitCollection()
	commits.Add(&model.PACCommit{
		SHA:        "abc123def456",
		ShortSHA:   "abc123d",
		Header:     "Fix bug",
		Referenced: true,
	})
	commits.Add(&model.PACCommit{
		SHA:        "def456ghi789",
		ShortSHA:   "def456g",
		Header:     "Untracked change",
		Referenced: false,
	})

	tasks := model.NewPACTaskCollection()
	task := tasks.FindOrCreate("JIRA-123")
	task.AddCommit(commits.Commits[0])

	unref := tasks.FindOrCreate("")
	unref.AddCommit(commits.Commits[1])

	generator := NewGenerator(tasks, commits)

	template := `# Changelog
{% for task in tasks.referenced %}
## {{ task.task_id }}
{% for commit in task.commits %}
- {{ commit.shortsha }}: {{ commit.header }}
{% endfor %}
{% endfor %}
## Unreferenced
{% for commit in tasks.unreferenced %}
- {{ commit.shortsha }}: {{ commit.header }}
{% endfor %}`

	settings := &config.Settings{}

	result, err := generator.GenerateToString(template, settings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Check for expected content
	if !strings.Contains(result, "JIRA-123") {
		t.Error("expected result to contain 'JIRA-123'")
	}

	if !strings.Contains(result, "abc123d") {
		t.Error("expected result to contain short SHA 'abc123d'")
	}

	if !strings.Contains(result, "Fix bug") {
		t.Error("expected result to contain 'Fix bug'")
	}

	if !strings.Contains(result, "Unreferenced") {
		t.Error("expected result to contain 'Unreferenced' section")
	}

	if !strings.Contains(result, "def456g") {
		t.Error("expected result to contain unreferenced commit SHA")
	}
}

func TestGenerator_PrepareLiquidData(t *testing.T) {
	commits := model.NewPACCommitCollection()
	commits.Add(&model.PACCommit{SHA: "1", Referenced: true})
	commits.Add(&model.PACCommit{SHA: "2", Referenced: true})
	commits.Add(&model.PACCommit{SHA: "3", Referenced: false})

	tasks := model.NewPACTaskCollection()

	generator := NewGenerator(tasks, commits)

	settings := &config.Settings{
		Properties: map[string]any{
			"version": "1.0.0",
		},
	}

	data := generator.prepareLiquidData(settings)

	// Check commit counts
	if data["pac_c_count"] != 3 {
		t.Errorf("pac_c_count: got %v, want 3", data["pac_c_count"])
	}

	if data["pac_c_referenced"] != 2 {
		t.Errorf("pac_c_referenced: got %v, want 2", data["pac_c_referenced"])
	}

	if data["pac_c_unreferenced"] != 1 {
		t.Errorf("pac_c_unreferenced: got %v, want 1", data["pac_c_unreferenced"])
	}

	// Check properties
	props, ok := data["properties"].(map[string]any)
	if !ok {
		t.Fatal("expected properties to be map[string]any")
	}

	if props["version"] != "1.0.0" {
		t.Errorf("properties.version: got %v, want '1.0.0'", props["version"])
	}
}

func TestGenerator_Generate_WritesFile(t *testing.T) {
	// Create temp directory for output
	tmpDir, err := os.MkdirTemp("", "pac-report-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a temp template file
	templateContent := "# Report\nCount: {{ pac_c_count }}"
	templatePath := filepath.Join(tmpDir, "template.md")
	if writeErr := os.WriteFile(templatePath, []byte(templateContent), 0o644); writeErr != nil {
		t.Fatalf("failed to write template: %v", writeErr)
	}

	outputPath := filepath.Join(tmpDir, "output", "report.md")

	commits := model.NewPACCommitCollection()
	commits.Add(&model.PACCommit{SHA: "1"})
	commits.Add(&model.PACCommit{SHA: "2"})

	tasks := model.NewPACTaskCollection()

	generator := NewGenerator(tasks, commits)

	settings := &config.Settings{
		Templates: []config.TemplateConfig{
			{
				Location: templatePath,
				Output:   outputPath,
			},
		},
	}

	err = generator.Generate(settings)
	if err != nil {
		t.Fatalf("Generate failed: %v", err)
	}

	// Verify output file was created
	content, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("failed to read output file: %v", err)
	}

	if !strings.Contains(string(content), "Count: 2") {
		t.Errorf("output should contain 'Count: 2', got: %s", string(content))
	}
}

func TestGenerator_RenderDefaultTemplate(t *testing.T) {
	// Test with the actual default.md template structure
	commits := model.NewPACCommitCollection()
	commits.Add(&model.PACCommit{
		SHA:        "abc123",
		ShortSHA:   "abc123",
		Header:     "JIRA-100: First commit",
		Referenced: true,
	})
	commits.Add(&model.PACCommit{
		SHA:        "def456",
		ShortSHA:   "def456",
		Header:     "Random fix",
		Referenced: false,
	})

	tasks := model.NewPACTaskCollection()

	// Add referenced task
	task := tasks.FindOrCreate("JIRA-100")
	task.AddCommit(commits.Commits[0])

	// Add unreferenced
	unref := tasks.FindOrCreate("")
	unref.AddCommit(commits.Commits[1])

	generator := NewGenerator(tasks, commits)

	// Use the same structure as default.md
	template := `# PAC Changelog
{% for task in tasks.referenced %}
## {{task.task_id}}
{% for commit in task.commits %}
- {{commit.shortsha}}: {{commit.header}}
{% endfor %}
{% endfor %}
## Unspecified
{% for commit in tasks.unreferenced %}
- {{commit.shortsha}}: {{commit.header}}
{% endfor %}`

	settings := &config.Settings{}

	result, err := generator.GenerateToString(template, settings)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify structure
	if !strings.Contains(result, "# PAC Changelog") {
		t.Error("missing header")
	}

	if !strings.Contains(result, "## JIRA-100") {
		t.Error("missing task ID section")
	}

	if !strings.Contains(result, "abc123: JIRA-100: First commit") {
		t.Errorf("missing referenced commit, got: %s", result)
	}

	if !strings.Contains(result, "## Unspecified") {
		t.Error("missing Unspecified section")
	}

	if !strings.Contains(result, "def456: Random fix") {
		t.Errorf("missing unreferenced commit, got: %s", result)
	}
}
