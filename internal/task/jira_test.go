package task

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/config"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/model"
)

func TestJiraTaskSystem_Name(t *testing.T) {
	cfg := config.TaskSystemConfig{
		Name: "jira-test",
	}
	jira := NewJiraTaskSystem(&cfg)

	if got := jira.Name(); got != "jira-test" {
		t.Errorf("Name() = %q, want %q", got, "jira-test")
	}
}

func TestJiraTaskSystem_Apply_SkipsUnreferenced(t *testing.T) {
	cfg := config.TaskSystemConfig{
		Name: "jira",
	}
	jira := NewJiraTaskSystem(&cfg)

	tasks := model.NewPACTaskCollection()
	// Create unreferenced task (empty ID)
	unref := tasks.FindOrCreate("")
	unref.AppliesTo["jira"] = true

	err := jira.Apply(tasks)
	if err != nil {
		t.Errorf("Apply() error = %v, want nil", err)
	}
}

func TestJiraTaskSystem_Apply_SkipsNonApplying(t *testing.T) {
	cfg := config.TaskSystemConfig{
		Name: "jira",
	}
	jira := NewJiraTaskSystem(&cfg)

	tasks := model.NewPACTaskCollection()
	task := tasks.FindOrCreate("TASK-1")
	// Don't set AppliesTo for jira

	err := jira.Apply(tasks)
	if err != nil {
		t.Errorf("Apply() error = %v, want nil", err)
	}

	// Data should not be set
	if task.Data != nil {
		t.Error("expected Data to be nil for non-applying task")
	}
}

func TestJiraTaskSystem_Apply_WithMockServer(t *testing.T) {
	// Create a mock Jira server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify basic auth
		user, pass, ok := r.BasicAuth()
		if !ok || user != "testuser" || pass != "testpass" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// Return mock Jira response
		response := map[string]any{
			"key": "TASK-123",
			"fields": map[string]any{
				"summary": "Test issue",
				"status": map[string]any{
					"name": "Done",
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	cfg := config.TaskSystemConfig{
		Name:        "jira",
		QueryString: server.URL + "/rest/api/2/issue/#{task_id}",
		Username:    "testuser",
		Password:    "testpass",
	}
	jira := NewJiraTaskSystem(&cfg)

	tasks := model.NewPACTaskCollection()
	task := tasks.FindOrCreate("TASK-123")
	task.AppliesTo["jira"] = true

	err := jira.Apply(tasks)
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

	if data["key"] != "TASK-123" {
		t.Errorf("Data[key] = %v, want TASK-123", data["key"])
	}
}

func TestJiraTaskSystem_Apply_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	cfg := config.TaskSystemConfig{
		Name:        "jira",
		QueryString: server.URL + "/rest/api/2/issue/#{task_id}",
	}
	jira := NewJiraTaskSystem(&cfg)

	tasks := model.NewPACTaskCollection()
	task := tasks.FindOrCreate("TASK-123")
	task.AppliesTo["jira"] = true

	err := jira.Apply(tasks)
	if err == nil {
		t.Error("expected error for server error, got nil")
	}

	// Task should be marked as unknown
	if !task.Labels["unknown"] {
		t.Error("expected task to have 'unknown' label after error")
	}
}

func TestJiraTaskSystem_Apply_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write([]byte("not valid json")); err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	cfg := config.TaskSystemConfig{
		Name:        "jira",
		QueryString: server.URL + "/rest/api/2/issue/#{task_id}",
	}
	jira := NewJiraTaskSystem(&cfg)

	tasks := model.NewPACTaskCollection()
	task := tasks.FindOrCreate("TASK-123")
	task.AppliesTo["jira"] = true

	err := jira.Apply(tasks)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestNoneTaskSystem_Name(t *testing.T) {
	cfg := config.TaskSystemConfig{
		Name: "none-test",
	}
	none := NewNoneTaskSystem(&cfg)

	if got := none.Name(); got != "none-test" {
		t.Errorf("Name() = %q, want %q", got, "none-test")
	}
}

func TestNoneTaskSystem_Apply_PreservesLabels(t *testing.T) {
	cfg := config.TaskSystemConfig{
		Name: "none",
	}
	none := NewNoneTaskSystem(&cfg)

	tasks := model.NewPACTaskCollection()
	task := tasks.FindOrCreate("TASK-1")
	task.AppliesTo["none"] = true
	task.AddLabel("feature")

	err := none.Apply(tasks)
	if err != nil {
		t.Errorf("Apply() error = %v, want nil", err)
	}

	// Labels should remain unchanged
	if !task.Labels["feature"] {
		t.Error("expected 'feature' label to remain")
	}
}

func TestCreateTaskSystem_None(t *testing.T) {
	cfg := config.TaskSystemConfig{
		Name: "none",
	}

	ts, err := CreateTaskSystem(&cfg)
	if err != nil {
		t.Fatalf("CreateTaskSystem() error = %v", err)
	}

	if _, ok := ts.(*NoneTaskSystem); !ok {
		t.Error("expected NoneTaskSystem type")
	}
}

func TestCreateTaskSystem_Jira(t *testing.T) {
	cfg := config.TaskSystemConfig{
		Name: "jira",
	}

	ts, err := CreateTaskSystem(&cfg)
	if err != nil {
		t.Fatalf("CreateTaskSystem() error = %v", err)
	}

	if _, ok := ts.(*JiraTaskSystem); !ok {
		t.Error("expected JiraTaskSystem type")
	}
}

func TestCreateTaskSystem_Unknown(t *testing.T) {
	cfg := config.TaskSystemConfig{
		Name: "unknown-system",
	}

	ts, err := CreateTaskSystem(&cfg)
	if err != nil {
		t.Fatalf("CreateTaskSystem() error = %v", err)
	}

	// Unknown systems default to NoneTaskSystem
	if _, ok := ts.(*NoneTaskSystem); !ok {
		t.Error("expected unknown system to default to NoneTaskSystem")
	}
}

func TestCreateAllTaskSystems(t *testing.T) {
	configs := []config.TaskSystemConfig{
		{Name: "none"},
		{Name: "jira"},
	}

	systems, err := CreateAllTaskSystems(configs)
	if err != nil {
		t.Fatalf("CreateAllTaskSystems() error = %v", err)
	}

	if len(systems) != 2 {
		t.Errorf("expected 2 systems, got %d", len(systems))
	}
}
