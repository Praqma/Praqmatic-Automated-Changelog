package task

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/config"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/logging"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/model"
)

// JiraTaskSystem fetches task data from a Jira REST API.
type JiraTaskSystem struct {
	config     config.TaskSystemConfig
	httpClient *http.Client
}

// NewJiraTaskSystem creates a new JiraTaskSystem.
func NewJiraTaskSystem(cfg *config.TaskSystemConfig) *JiraTaskSystem {
	return &JiraTaskSystem{
		config: *cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Name returns the task system name.
func (j *JiraTaskSystem) Name() string {
	return j.config.Name
}

// Apply fetches data from Jira for each task that applies to this system.
func (j *JiraTaskSystem) Apply(tasks *model.PACTaskCollection) error {
	var errors []string

	for _, task := range tasks.Tasks {
		// Skip unreferenced commits (empty task ID)
		if task.TaskID == "" {
			continue
		}

		// Skip tasks that don't apply to this system
		if !task.AppliesTo[j.config.Name] {
			continue
		}

		// Fetch data from Jira
		data, err := j.fetchTaskData(task.TaskID)
		if err != nil {
			logging.Warn("Jira error for %s: %v", task.TaskID, err)
			// Mark task as unknown when fetch fails
			task.ClearLabels()
			task.AddLabel("unknown")
			errors = append(errors, fmt.Sprintf("%s: %v", task.TaskID, err))
			continue
		}

		// Populate task with Jira data
		task.Data = data
		task.Attributes["data"] = data

		logging.Info("Applied Jira data to %s", task.TaskID)
	}

	if len(errors) > 0 {
		return fmt.Errorf("jira errors: %s", strings.Join(errors, "; "))
	}

	return nil
}

// fetchTaskData retrieves task data from the Jira API.
func (j *JiraTaskSystem) fetchTaskData(taskID string) (map[string]any, error) {
	// Build the URL by interpolating task_id
	url := strings.ReplaceAll(j.config.QueryString, "#{task_id}", taskID)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication if configured
	if j.config.Username != "" {
		req.SetBasicAuth(j.config.Username, j.config.Password)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := j.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}

	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return data, nil
}

// SetHTTPClient allows setting a custom HTTP client (useful for testing).
func (j *JiraTaskSystem) SetHTTPClient(client *http.Client) {
	j.httpClient = client
}
