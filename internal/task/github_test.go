package task

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/config"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/model"
)

func TestGitHubTaskSystem_Name(t *testing.T) {
	cfg := config.TaskSystemConfig{
		Name: "github-test",
	}
	github := NewGitHubTaskSystem(&cfg)

	if got := github.Name(); got != "github-test" {
		t.Errorf("Name() = %q, want %q", got, "github-test")
	}
}

func TestGitHubTaskSystem_Apply_SkipsUnreferenced(t *testing.T) {
	cfg := config.TaskSystemConfig{
		Name: "github",
	}
	github := NewGitHubTaskSystem(&cfg)

	tasks := model.NewPACTaskCollection()
	// Create unreferenced task (empty ID)
	unref := tasks.FindOrCreate("")
	unref.AppliesTo["github"] = true

	err := github.Apply(tasks)
	if err != nil {
		t.Errorf("Apply() error = %v, want nil", err)
	}
}

func TestGitHubTaskSystem_Apply_SkipsNonApplying(t *testing.T) {
	cfg := config.TaskSystemConfig{
		Name: "github",
	}
	github := NewGitHubTaskSystem(&cfg)

	tasks := model.NewPACTaskCollection()
	task := tasks.FindOrCreate("#123")
	// Don't set AppliesTo for github

	err := github.Apply(tasks)
	if err != nil {
		t.Errorf("Apply() error = %v, want nil", err)
	}

	// Data should not be set
	if task.Data != nil {
		t.Error("expected Data to be nil for non-applying task")
	}
}

func TestGitHubTaskSystem_Apply_WithMockServer(t *testing.T) {
	// Create a mock GitHub API server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify authorization header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Verify GitHub-specific headers
		if r.Header.Get("Accept") != "application/vnd.github+json" {
			t.Error("expected GitHub Accept header")
		}

		// Return mock GitHub issue response
		response := map[string]any{
			"id":     123,
			"number": 42,
			"title":  "Test issue",
			"state":  "open",
			"body":   "This is a test issue",
			"user": map[string]any{
				"login": "testuser",
			},
			"labels": []map[string]any{
				{"name": "bug"},
				{"name": "enhancement"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	cfg := config.TaskSystemConfig{
		Name:        "github",
		QueryString: server.URL + "/repos/owner/repo/issues/#{task_id}",
		Password:    "test-token",
	}
	github := NewGitHubTaskSystem(&cfg)

	tasks := model.NewPACTaskCollection()
	task := tasks.FindOrCreate("#42")
	task.AppliesTo["github"] = true

	err := github.Apply(tasks)
	if err != nil {
		t.Errorf("Apply() error = %v, want nil", err)
	}

	// Verify data was populated
	if task.Data == nil {
		t.Fatal("expected Data to be populated")
	}

	data, ok := task.Data.(map[string]any)
	if !ok {
		t.Fatal("expected Data to be map[string]any")
	}

	if data["title"] != "Test issue" {
		t.Errorf("Data[title] = %v, want 'Test issue'", data["title"])
	}

	if data["state"] != "open" {
		t.Errorf("Data[state] = %v, want 'open'", data["state"])
	}
}

func TestGitHubTaskSystem_Apply_NotFoundError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	cfg := config.TaskSystemConfig{
		Name:        "github",
		QueryString: server.URL + "/repos/owner/repo/issues/#{task_id}",
	}
	github := NewGitHubTaskSystem(&cfg)

	tasks := model.NewPACTaskCollection()
	task := tasks.FindOrCreate("#999")
	task.AppliesTo["github"] = true

	err := github.Apply(tasks)
	if err == nil {
		t.Error("expected error for not found, got nil")
	}

	// Task should be marked as unknown
	if !task.Labels["unknown"] {
		t.Error("expected task to have 'unknown' label after error")
	}
}

func TestGitHubTaskSystem_Apply_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	cfg := config.TaskSystemConfig{
		Name:        "github",
		QueryString: server.URL + "/repos/owner/repo/issues/#{task_id}",
	}
	github := NewGitHubTaskSystem(&cfg)

	tasks := model.NewPACTaskCollection()
	task := tasks.FindOrCreate("#123")
	task.AppliesTo["github"] = true

	err := github.Apply(tasks)
	if err == nil {
		t.Error("expected error for server error, got nil")
	}

	// Task should be marked as unknown
	if !task.Labels["unknown"] {
		t.Error("expected task to have 'unknown' label after error")
	}
}

func TestGitHubTaskSystem_Apply_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, err := w.Write([]byte("not valid json"))
		if err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	cfg := config.TaskSystemConfig{
		Name:        "github",
		QueryString: server.URL + "/repos/owner/repo/issues/#{task_id}",
	}
	github := NewGitHubTaskSystem(&cfg)

	tasks := model.NewPACTaskCollection()
	task := tasks.FindOrCreate("#123")
	task.AppliesTo["github"] = true

	err := github.Apply(tasks)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestGitHubTaskSystem_ExtractIssueNumber(t *testing.T) {
	cfg := config.TaskSystemConfig{Name: "github"}
	github := NewGitHubTaskSystem(&cfg)

	tests := []struct {
		input    string
		expected string
	}{
		{"#123", "123"},
		{"123", "123"},
		{"owner/repo#456", "456"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := github.extractIssueNumber(tt.input)
			if got != tt.expected {
				t.Errorf("extractIssueNumber(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestGitHubTaskSystem_ParseFullReference(t *testing.T) {
	cfg := config.TaskSystemConfig{Name: "github"}
	github := NewGitHubTaskSystem(&cfg)

	tests := []struct {
		input    string
		expected map[string]string
	}{
		{
			"owner/repo#123",
			map[string]string{"owner": "owner", "repo": "repo", "number": "123"},
		},
		{
			"my-org/my-repo#456",
			map[string]string{"owner": "my-org", "repo": "my-repo", "number": "456"},
		},
		{"#123", nil},
		{"123", nil},
		{"invalid", nil},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := github.parseFullReference(tt.input)
			if tt.expected == nil {
				if got != nil {
					t.Errorf("parseFullReference(%q) = %v, want nil", tt.input, got)
				}
				return
			}
			if got == nil {
				t.Errorf("parseFullReference(%q) = nil, want %v", tt.input, tt.expected)
				return
			}
			for k, v := range tt.expected {
				if got[k] != v {
					t.Errorf("parseFullReference(%q)[%s] = %q, want %q", tt.input, k, got[k], v)
				}
			}
		})
	}
}

func TestGitHubTaskSystem_BuildURL_WithQueryString(t *testing.T) {
	cfg := config.TaskSystemConfig{
		Name:        "github",
		QueryString: "https://api.github.com/repos/myorg/myrepo/issues/#{task_id}",
	}
	github := NewGitHubTaskSystem(&cfg)

	url := github.buildURL("#42")
	expected := "https://api.github.com/repos/myorg/myrepo/issues/42"

	if url != expected {
		t.Errorf("buildURL('#42') = %q, want %q", url, expected)
	}
}

func TestGitHubTaskSystem_BuildURL_FullReference(t *testing.T) {
	cfg := config.TaskSystemConfig{
		Name: "github",
	}
	github := NewGitHubTaskSystem(&cfg)

	url := github.buildURL("owner/repo#123")
	expected := "https://api.github.com/repos/owner/repo/issues/123"

	if url != expected {
		t.Errorf("buildURL('owner/repo#123') = %q, want %q", url, expected)
	}
}

func TestGitHubTaskSystem_SetHTTPClient(t *testing.T) {
	cfg := config.TaskSystemConfig{Name: "github"}
	github := NewGitHubTaskSystem(&cfg)

	customClient := &http.Client{}
	github.SetHTTPClient(customClient)

	if github.httpClient != customClient {
		t.Error("SetHTTPClient did not update httpClient")
	}
}
