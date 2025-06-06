package task

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Praqma/Praqmatic-Automated-Changelog/src/model"
	"github.com/google/go-github/v72/github"
)

type GHTask struct {
	Github *github.Client
	Owner  string
	Repo   string
	Branch string
	PR     int
	Title  string
	Body   string
	Labels []string
	Assignees []string
	Reviewers []string
}

// NewGHTask creates an instance that connects to a remote GitHub repository based on a URL only. it fills in the
// owner, repo and branch fields based on the URL.
func NewGHTask(url string, token string) (*GHTask, error) {
    client := github.NewClient(nil).WithAuthToken(token)
	
	task := &GHTask{
		Github: client,
	}
	
	fmt.Printf("Initializing GHTask with URL: %s\n", url)
	// Parse the URL to extract owner and repo
	owner, repo, branch := parseGitHubURL(url)
	fmt.Printf("Parsed URL: Owner=%s, Repo=%s, Branch=%s\n", owner, repo, branch)
	task.Owner = owner
	task.Repo = repo
	task.Branch = branch
	
	return task, nil
}

// parseGitHubURL extracts the owner, repository name, and branch from a GitHub URL
func parseGitHubURL(repoURL string) (owner, repo, branch string) {
	// Default branch if not specified
	branch = "main"
	
	// Handle URLs of different formats:
	// - https://github.com/username/repo
	// - https://github.com/username/repo.git
	// - https://github.com/username/repo/tree/branch
	// - git@github.com:username/repo.git
	
	// Format: https://github.com/username/repo[.git][/tree/branch]
	httpRegex := regexp.MustCompile(`github\.com/([^/]+)/([^/\.]+)(?:\.git)?(?:/tree/([^/]+))?`)
	sshRegex := regexp.MustCompile(`git@github\.com:([^/]+)/([^/\.]+)(?:\.git)?`)
	
	// Check for HTTP(S) URL format
	if matches := httpRegex.FindStringSubmatch(repoURL); matches != nil {
		owner = matches[1]
		repo = matches[2]
		if len(matches) > 3 && matches[3] != "" {
			branch = matches[3]
		}
		return
	}
	
	// Check for SSH URL format
	if matches := sshRegex.FindStringSubmatch(repoURL); matches != nil {
		owner = matches[1]
		repo = matches[2]
		return
	}
	
	// If the URL doesn't match expected formats, try to parse it as a simple owner/repo string
	parts := strings.Split(repoURL, "/")
	if len(parts) >= 2 {
		owner = parts[0]
		repo = strings.TrimSuffix(parts[1], ".git")
	}
	
	return
}

func (t *GHTask) GetIssueByNumber(number int) (*github.Issue, error) {
	ctx := context.Background()
	issue, resp, err := t.Github.Issues.Get(ctx, t.Owner, t.Repo, number)
	if err != nil {
		return nil, fmt.Errorf("error fetching issue #%d: %w", number, err)
	}
	
	// Check if the response indicates a 404 Not Found
	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("issue #%d not found", number)
	}
	
	return issue, nil
}

// GetPullRequests retrieves pull requests from the GitHub repository
func (t *GHTask) GetPullRequests(ctx context.Context, state string) ([]*github.PullRequest, error) {
	opts := &github.PullRequestListOptions{
		State: state,
		ListOptions: github.ListOptions{
			PerPage: 100, // Get maximum 100 PRs per page (GitHub API limit)
		},
	}
	
	// Store all pull requests
	var allPRs []*github.PullRequest
	
	// Get PRs page by page
	for {
		prs, resp, err := t.Github.PullRequests.List(ctx, t.Owner, t.Repo, opts)
		if err != nil {
			return nil, fmt.Errorf("error fetching pull requests: %w", err)
		}
		
		// Add PRs from current page to our collection
		allPRs = append(allPRs, prs...)
		
		// If there are no more pages, break
		if resp.NextPage == 0 {
			break
		}
		
		// Otherwise, update page number for next iteration
		opts.Page = resp.NextPage
	}
	
	return allPRs, nil
}

// ProcessCommits processes commits for GitHub task system and creates tasks
func (t *GHTask) ProcessCommits(taskSystem model.TaskSystem, commits *model.PACCommitCollection) (*model.PACTaskCollection, error) {
	tasks := model.NewPACTaskCollection()

	fmt.Println("Processing commits for GitHub task system:", taskSystem.Name)

	for _, commit := range commits.Commits {
		matchID := extractTaskID(commit, taskSystem.Regex)
		if matchID == "" {
			// Create a task with empty ID for unmatched commits
			task := model.NewPACTask("")
			task.AddCommit(commit)
			tasks.Add(task)
			continue
		}

		task, err := t.createTaskFromCommit(commit, matchID)
		if err != nil {
			return nil, fmt.Errorf("error creating task for commit %s: %w", commit.SHA, err)
		}

		tasks.Add(task)
	}

	return tasks, nil
}


// createTaskFromCommit creates a PACTask from a commit and GitHub issue
func (t *GHTask) createTaskFromCommit(commit *model.PACCommit, matchID string) (*model.PACTask, error) {
	issueNumber, err := strconv.Atoi(strings.TrimPrefix(matchID, "#"))
	if err != nil {
		return nil, fmt.Errorf("error converting issue number from commit message: %w", err)
	}

	issue, err := t.GetIssueByNumber(issueNumber)
	if err != nil {
		return nil, fmt.Errorf("error fetching issue #%s: %w", matchID, err)
	}

	task := model.NewPACTask(matchID)
	task.AddCommit(commit)
	
	// Populate all the new fields
	task.Number = issueNumber
	task.Title = issue.GetTitle()
	task.Body = issue.GetBody()
	task.State = issue.GetState()
	task.URL = issue.GetHTMLURL()
	task.CreatedAt = issue.CreatedAt.GetTime()
	task.UpdatedAt = issue.UpdatedAt.GetTime()
	task.ClosedAt = issue.ClosedAt.GetTime()
	
	// Set task type
	if issue.IsPullRequest() {
		task.TaskType = "pull_request"
	} else {
		task.TaskType = issue.Type.GetName()
		if task.TaskType == "" {
			task.TaskType = "issue"
		}
	}
	
	// Add author
	if issue.User != nil {
		task.Author = issue.User.GetLogin()
	}
	
	// Add assignees
	for _, assignee := range issue.Assignees {
		if assignee != nil {
			task.AddAssignee(assignee.GetLogin())
		}
	}
	
	// Add milestone
	if issue.Milestone != nil {
		task.Milestone = issue.Milestone.GetTitle()
	}
	
	// Add labels
	for _, label := range issue.Labels {
		task.AddLabel(label.GetName())
	}
	
	// Store the raw issue data for potential future use
	//task.Data = issue

	return task, nil
}

// GetPullRequestByNumber retrieves a pull request by its number
func (t *GHTask) GetPullRequestByNumber(number int) (*github.PullRequest, error) {
	ctx := context.Background()
	pr, resp, err := t.Github.PullRequests.Get(ctx, t.Owner, t.Repo, number)
	if err != nil {
		return nil, fmt.Errorf("error fetching pull request #%d: %w", number, err)
	}
	
	// Check if the response indicates a 404 Not Found
	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("pull request #%d not found", number)
	}
	
	return pr, nil
}

// EnrichTaskWithPRData enriches a task with pull request specific data
func (t *GHTask) EnrichTaskWithPRData(task *model.PACTask) error {
	if task.TaskType != "pull_request" {
		return nil
	}
	
	pr, err := t.GetPullRequestByNumber(task.Number)
	if err != nil {
		return err
	}
	
	// Add PR-specific data
	if pr.MergedAt != nil {
		// You might want to add a MergedAt field to PACTask
		// For now, store it in Data
		//task.Data = pr
	}
	
	// Add reviewers as a special label or in assignees
	for _, reviewer := range pr.RequestedReviewers {
		if reviewer != nil {
			// Could add as assignee or create a separate Reviewers field
			task.AddAssignee(reviewer.GetName())
		}
	}
	
	return nil
}

// GetName returns the name of this task system
func (t *GHTask) GetName() string {
	return "github"
}