package task

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/config"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/logging"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/model"
)

// GitHubTaskSystem fetches issue data from the GitHub REST API.
type GitHubTaskSystem struct {
	config     config.TaskSystemConfig
	httpClient *http.Client
}

// NewGitHubTaskSystem creates a new GitHubTaskSystem.
func NewGitHubTaskSystem(cfg *config.TaskSystemConfig) *GitHubTaskSystem {
	return &GitHubTaskSystem{
		config: *cfg,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Name returns the task system name.
func (g *GitHubTaskSystem) Name() string {
	return g.config.Name
}

// Apply fetches data from GitHub for each task that applies to this system.
func (g *GitHubTaskSystem) Apply(tasks *model.PACTaskCollection) error {
	var errors []string

	for _, task := range tasks.Tasks {
		// Skip unreferenced commits (empty task ID)
		if task.TaskID == "" {
			continue
		}

		// Skip tasks that don't apply to this system
		if !task.AppliesTo[g.config.Name] {
			continue
		}

		// Fetch data from GitHub
		data, err := g.fetchIssueData(task.TaskID)
		if err != nil {
			logging.Warn("GitHub error for %s: %v", task.TaskID, err)
			// Mark task as unknown when fetch fails
			task.ClearLabels()
			task.AddLabel("unknown")
			errors = append(errors, fmt.Sprintf("%s: %v", task.TaskID, err))
			continue
		}

		// Populate task with GitHub data
		task.Data = data
		task.Attributes["data"] = data

		logging.Info("Applied GitHub issue data to %s", task.TaskID)
	}

	if len(errors) > 0 {
		return fmt.Errorf("github errors: %s", strings.Join(errors, "; "))
	}

	return nil
}

// fetchIssueData retrieves issue data from the GitHub API.
// Task ID can be in formats: "#123", "123", "owner/repo#123"
func (g *GitHubTaskSystem) fetchIssueData(taskID string) (map[string]any, error) {
	url := g.buildURL(taskID)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add authentication if configured (token-based auth)
	if g.config.Password != "" {
		req.Header.Set("Authorization", "Bearer "+g.config.Password)
	}

	// Set headers
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "PAC-Changelog")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("issue not found: %s", taskID)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}

	var data map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return data, nil
}

// buildURL constructs the GitHub API URL for fetching an issue.
// If query_string is configured, it uses interpolation with #{task_id}.
// Otherwise, it extracts the issue number and uses the default GitHub API.
func (g *GitHubTaskSystem) buildURL(taskID string) string {
	// If query_string is configured, use it with interpolation
	if g.config.QueryString != "" {
		return strings.ReplaceAll(g.config.QueryString, "#{task_id}", g.extractIssueNumber(taskID))
	}

	// Default GitHub API URL format
	// Assumes the issue number is the taskID (without #)
	issueNumber := g.extractIssueNumber(taskID)

	// If we have a full reference (owner/repo#123), parse it
	if match := g.parseFullReference(taskID); match != nil {
		return fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%s",
			match["owner"], match["repo"], match["number"])
	}

	// For just issue numbers, query_string must be configured with the repo
	// Fall back to a placeholder that will fail gracefully
	return fmt.Sprintf("https://api.github.com/repos/OWNER/REPO/issues/%s", issueNumber)
}

// extractIssueNumber removes the # prefix from issue references.
func (g *GitHubTaskSystem) extractIssueNumber(taskID string) string {
	// Remove # prefix if present
	taskID = strings.TrimPrefix(taskID, "#")

	// If it's a full reference like owner/repo#123, extract just the number
	if idx := strings.LastIndex(taskID, "#"); idx != -1 {
		return taskID[idx+1:]
	}

	return taskID
}

// parseFullReference parses a full GitHub issue reference (owner/repo#123).
func (g *GitHubTaskSystem) parseFullReference(taskID string) map[string]string {
	// Pattern: owner/repo#123
	pattern := regexp.MustCompile(`^([^/]+)/([^#]+)#(\d+)$`)
	matches := pattern.FindStringSubmatch(taskID)

	if len(matches) != 4 {
		return nil
	}

	return map[string]string{
		"owner":  matches[1],
		"repo":   matches[2],
		"number": matches[3],
	}
}

// SetHTTPClient allows setting a custom HTTP client (useful for testing).
func (g *GitHubTaskSystem) SetHTTPClient(client *http.Client) {
	g.httpClient = client
}
