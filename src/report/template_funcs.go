package report

import (
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/Praqma/Praqmatic-Automated-Changelog/src/model"
)

// GetTemplateFuncs returns a FuncMap with custom functions for templates
func GetTemplateFuncs() template.FuncMap {
	return template.FuncMap{
		// Date formatting functions
		"formatDate": func(t *time.Time) string {
			if t == nil {
				return "N/A"
			}
			return t.Format("2006-01-02")
		},
		"formatDateTime": func(t *time.Time) string {
			if t == nil {
				return "N/A"
			}
			return t.Format("2006-01-02 15:04:05")
		},
		"formatDateCustom": func(t *time.Time, format string) string {
			if t == nil {
				return "N/A"
			}
			return t.Format(format)
		},
		
		// Task filtering functions
		"filterByLabel": func(tasks []*model.PACTask, label string) []*model.PACTask {
			var filtered []*model.PACTask
			for _, task := range tasks {
				for _, l := range task.Labels {
					if l == label {
						filtered = append(filtered, task)
						break
					}
				}
			}
			return filtered
		},
		"filterByState": func(tasks []*model.PACTask, state string) []*model.PACTask {
			var filtered []*model.PACTask
			for _, task := range tasks {
				if task.State == state {
					filtered = append(filtered, task)
				}
			}
			return filtered
		},
		"filterByType": func(tasks []*model.PACTask, taskType string) []*model.PACTask {
			var filtered []*model.PACTask
			for _, task := range tasks {
				if task.TaskType == taskType {
					filtered = append(filtered, task)
				}
			}
			return filtered
		},
		"filterByAuthor": func(tasks []*model.PACTask, author string) []*model.PACTask {
			var filtered []*model.PACTask
			for _, task := range tasks {
				if task.Author == author {
					filtered = append(filtered, task)
				}
			}
			return filtered
		},
		
		// Task grouping functions
		"groupByAuthor": func(tasks []*model.PACTask) map[string][]*model.PACTask {
			grouped := make(map[string][]*model.PACTask)
			for _, task := range tasks {
				author := task.Author
				if author == "" {
					author = "Unknown"
				}
				grouped[author] = append(grouped[author], task)
			}
			return grouped
		},
		"groupByMilestone": func(tasks []*model.PACTask) map[string][]*model.PACTask {
			grouped := make(map[string][]*model.PACTask)
			for _, task := range tasks {
				milestone := task.Milestone
				if milestone == "" {
					milestone = "No Milestone"
				}
				grouped[milestone] = append(grouped[milestone], task)
			}
			return grouped
		},
		"groupByState": func(tasks []*model.PACTask) map[string][]*model.PACTask {
			grouped := make(map[string][]*model.PACTask)
			for _, task := range tasks {
				grouped[task.State] = append(grouped[task.State], task)
			}
			return grouped
		},
		
		// Counting functions
		"countTasks": func(tasks []*model.PACTask) int {
			return len(tasks)
		},
		"countCommits": func(tasks []*model.PACTask) int {
			count := 0
			for _, task := range tasks {
				count += len(task.Commits)
			}
			return count
		},
		"countByLabel": func(tasks []*model.PACTask, label string) int {
			count := 0
			for _, task := range tasks {
				for _, l := range task.Labels {
					if l == label {
						count++
						break
					}
				}
			}
			return count
		},
		
		// Utility functions
		"hasLabel": func(task *model.PACTask, label string) bool {
			for _, l := range task.Labels {
				if l == label {
					return true
				}
			}
			return false
		},
		"hasAnyLabel": func(task *model.PACTask, labels ...string) bool {
			for _, label := range labels {
				if hasLabel(task, label) {
					return true
				}
			}
			return false
		},
		"joinLabels": func(labels []string, separator string) string {
			return strings.Join(labels, separator)
		},
		"joinAssignees": func(assignees []string, separator string) string {
			return strings.Join(assignees, separator)
		},
		"truncate": func(s string, length int) string {
			if len(s) <= length {
				return s
			}
			return s[:length] + "..."
		},
		"firstLine": func(s string) string {
			lines := strings.Split(s, "\n")
			if len(lines) > 0 {
				return strings.TrimSpace(lines[0])
			}
			return s
		},
		"contains": func(s, substr string) bool {
			return strings.Contains(s, substr)
		},
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
		"title": strings.Title,
		
		// URL generation helpers
		"issueURL": func(task *model.PACTask) string {
			if task.URL != "" {
				return task.URL
			}
			return fmt.Sprintf("#%d", task.Number)
		},
		
		// Commit helpers
		"commitSHA": func(commit *model.PACCommit, length int) string {
			if len(commit.SHA) <= length {
				return commit.SHA
			}
			return commit.SHA[:length]
		},
		"commitHeader": func(commit *model.PACCommit) string {
			return commit.Header()
		},
		
		// Statistics functions
		"taskStats": func(tasks []*model.PACTask) map[string]interface{} {
			stats := make(map[string]interface{})
			openCount := 0
			closedCount := 0
			issueCount := 0
			prCount := 0
			
			labelCounts := make(map[string]int)
			authorCounts := make(map[string]int)
			
			for _, task := range tasks {
				if task.IsOpen() {
					openCount++
				} else if task.IsClosed() {
					closedCount++
				}
				
				if task.TaskType == "issue" {
					issueCount++
				} else if task.TaskType == "pull_request" {
					prCount++
				}
				
				// Count labels
				for _, label := range task.Labels {
					labelCounts[label]++
				}
				
				// Count authors
				if task.Author != "" {
					authorCounts[task.Author]++
				}
			}
			
			stats["total"] = len(tasks)
			stats["open"] = openCount
			stats["closed"] = closedCount
			stats["issues"] = issueCount
			stats["pull_requests"] = prCount
			stats["labels"] = labelCounts
			stats["authors"] = authorCounts
			
			return stats
		},
		
		// Sorting helpers (returns sorted copy, doesn't modify original)
		"sortTasksByDate": func(tasks []*model.PACTask) []*model.PACTask {
			// Create a copy to avoid modifying the original
			sorted := make([]*model.PACTask, len(tasks))
			copy(sorted, tasks)
			
			// Sort by created date (newest first)
			for i := 0; i < len(sorted)-1; i++ {
				for j := i + 1; j < len(sorted); j++ {
					if sorted[i].CreatedAt != nil && sorted[j].CreatedAt != nil {
						if sorted[i].CreatedAt.Before(*sorted[j].CreatedAt) {
							sorted[i], sorted[j] = sorted[j], sorted[i]
						}
					}
				}
			}
			return sorted
		},
		
		// Task collection functions
		"referencedTasks": func(tc *model.PACTaskCollection) []*model.PACTask {
			return tc.GetReferencedTasks()
		},
		"unreferencedCommits": func(tc *model.PACTaskCollection) []*model.PACCommit {
			return tc.GetUnreferencedCommits()
		},
		"allTasks": func(tc *model.PACTaskCollection) []*model.PACTask {
			return tc.Tasks
		},
		"tasksByLabel": func(tc *model.PACTaskCollection) map[string][]*model.PACTask {
			return tc.GetTasksByLabel()
		},
		"getTaskByID": func(tc *model.PACTaskCollection, id string) *model.PACTask {
			return tc.GetTaskByID(id)
		},
	}
}

// Helper function used internally
func hasLabel(task *model.PACTask, label string) bool {
	for _, l := range task.Labels {
		if l == label {
			return true
		}
	}
	return false
}
