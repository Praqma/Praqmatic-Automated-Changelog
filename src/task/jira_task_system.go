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
	JiraClient *jira.Client
}

// NewJiraTaskSystem creates a new Jira task system instance
func NewJiraTaskSystem(baseURL string) (*JiraTaskSystem, error) {

	tp := jira.BearerAuthTransport{
		Token: os.Getenv("JIRA_TOKEN"),
	}
	client, err := jira.NewClient(tp.Client(), baseURL)
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

	task.Title = issue.Fields.Summary
	task.ID = issue.Key 
	task.URL = "https://" + j.JiraClient.GetBaseURL().Host + "/browse/" + issue.Key
	task.Author = issue.Fields.Reporter.Name 
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

	// Remove nil values from the map
	cleanedData := removeNilValues(issueData)

	task.Data[j.GetName()] = cleanedData


	return task, nil
}

func removeNilValues(data interface{}) interface{} {
	switch v := data.(type) {
	case map[string]interface{}:
		cleaned := make(map[string]interface{})
		for key, value := range v {
			if value != nil {
				if cleanedValue := removeNilValues(value); cleanedValue != nil {
					cleaned[key] = cleanedValue
				}
			}
		}
		if len(cleaned) == 0 {
			return nil
		}
		return cleaned
	case []interface{}:
		var cleaned []interface{}
		for _, item := range v {
			if item != nil {
				if cleanedValue := removeNilValues(item); cleanedValue != nil {
					cleaned = append(cleaned, cleanedValue)
				}
			}
		}
		if len(cleaned) == 0 {
			return nil
		}
		return cleaned
	default:
		return data
	}
}

func (j *JiraTaskSystem) fetchIssueFromJira(taskID string) (jira.Issue, error) {
	issue, resp, err := j.JiraClient.Issue.Get(taskID, nil)
	if err != nil {
		return jira.Issue{}, fmt.Errorf("error fetching Jira issue %s: %d", taskID, resp.StatusCode)
	}
	if resp.StatusCode != 200 {
		return jira.Issue{}, fmt.Errorf("error fetching Jira issue %s: received status code %d", taskID, resp.StatusCode)
	}
	return *issue, nil
}

