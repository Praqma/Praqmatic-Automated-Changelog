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
			return fmt.Sprintf("#%s", task.ID)
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

			issueCount := 0
			
			labelCounts := make(map[string]int)
			authorCounts := make(map[string]int)
			
			for _, task := range tasks {
				
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
			stats["issues"] = issueCount
			stats["labels"] = labelCounts
			stats["authors"] = authorCounts
			
			return stats
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
