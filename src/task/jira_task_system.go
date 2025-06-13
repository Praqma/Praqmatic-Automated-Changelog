package task

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Praqma/Praqmatic-Automated-Changelog/src/model"
	jira "github.com/andygrunwald/go-jira"
)

// JiraTaskSystem handles Jira task integration
type JiraTaskSystem struct {
	JiraClient *jira.Client // Assuming JiraClient is a struct that handles Jira API interactions
}

// NewJiraTaskSystem creates a new Jira task system instance
func NewJiraTaskSystem(baseURL string) (*JiraTaskSystem, error) {

	jt := jira.BasicAuthTransport{
		Username: os.Getenv("JIRA_USER"),
		Password: os.Getenv("JIRA_TOKEN"),
	}
	client, err := jira.NewClient(jt.Client(), baseURL)
	if err != nil {

		return &JiraTaskSystem{
			JiraClient: nil,
		}, fmt.Errorf("error creating Jira client: %w", err)
	}
	return &JiraTaskSystem{
		JiraClient: client,
	}, nil
}

// ProcessCommits processes commits for Jira task system and creates tasks
func (j *JiraTaskSystem) ProcessCommits(taskSystem model.TaskSystem, commits *model.PACCommitCollection) (*model.PACTaskCollection, error) {
	tasks := model.NewPACTaskCollection()

	fmt.Println("Processing commits for Jira task system:", taskSystem.Name)

	for _, commit := range commits.Commits {
		taskID := extractTaskID(commit, taskSystem.Regex)
		if taskID == "" {
			// Create a task with empty ID for unmatched commits
			task := model.NewPACTask("")
			task.AddCommit(commit)
			tasks.Add(task)
			continue
		}

		task, err := j.createTaskFromCommit(commit, taskID)
		if err != nil {
			// If we can't fetch from Jira, create a basic task
			fmt.Printf("Warning: Could not fetch Jira issue %s: %v\n", taskID, err)
			task = model.NewPACTask(taskID)
			task.AddCommit(commit)
		}

		tasks.Add(task)
	}

	return tasks, nil
}

// GetName returns the name of this task system
func (j *JiraTaskSystem) GetName() string {
	return "jira"
}

// createTaskFromCommit creates a PACTask from a commit and Jira issue
func (j *JiraTaskSystem) createTaskFromCommit(commit *model.PACCommit, taskID string) (*model.PACTask, error) {
	issue, err := j.fetchIssueFromJira(taskID)
	if err != nil {
		return nil, fmt.Errorf("error fetching Jira issue %s: %w", taskID, err)
	}

	task := model.NewPACTask(taskID)
	task.AddCommit(commit)

	// Populate task fields from Jira issue
	task.Title = issue.Fields.Summary
	task.ID = issue.ID // Assuming ID is the unique identifier for the task
	task.URL = j.JiraClient.GetBaseURL().RawPath
	task.Author = issue.Fields.Reporter.Name // Assuming Assignee is not nil
	task.AddAssignee(issue.Fields.Assignee.Name)

	for _, label := range issue.Fields.Labels {
		if label != "" {
			task.AddLabel(label)
		}
	}

	// Store the raw issue data for potential future use
	// Convert issue struct to JSON then to map
	issueJSON, err := json.Marshal(issue)
	if err != nil {
		return nil, fmt.Errorf("error marshaling issue to JSON: %w", err)
	}

	var issueData map[string]interface{}
	if err := json.Unmarshal(issueJSON, &issueData); err != nil {
		return nil, fmt.Errorf("error unmarshaling issue JSON to map: %w", err)
	}

	task.Data[j.GetName()] = issueData


	return task, nil
}

func (j *JiraTaskSystem) fetchIssueFromJira(taskID string) (jira.Issue, error) {
	issue, resp, err := j.JiraClient.Issue.Get(taskID, nil)
	if err != nil {
		return jira.Issue{}, fmt.Errorf("error fetching Jira issue %s: %w", taskID, err)
	}
	if resp.StatusCode != 200 {
		return jira.Issue{}, fmt.Errorf("error fetching Jira issue %s: received status code %d", taskID, resp.StatusCode)
	}
	return *issue, nil
}

