package task

import (
	"github.com/Praqma/Praqmatic-Automated-Changelog/src/model"
)

// TaskIDList generates a collection of tasks based on the commits found
func TaskIDList(taskSystem []model.TaskSystem, commits *model.PACCommitCollection) *model.PACTaskCollection {
	tasks := model.NewPACTaskCollection()

	for _, commit := range commits.Commits {
		referenced := false

		// Loop over each task system
		for _, taskSystem := range taskSystem {
			// Get the delimiter if defined
			var splitPattern string
			if taskSystem.Delimiter != "" {
				// Convert the Ruby-style delimiter to a Go-compatible one
				// This simplified implementation handles common cases
				if taskSystem.Delimiter == "/,/" {
					splitPattern = ","
				} else if taskSystem.Delimiter == "/\\s+/" {
					splitPattern = "\\s+"
				} else {
					// Strip the leading/trailing slashes for other patterns
					delimLen := len(taskSystem.Delimiter)
					if delimLen > 2 {
						splitPattern = taskSystem.Delimiter[1 : delimLen-1]
					}
				}
			}

			// Convert RegexRules to the format expected by MatchTask
			var patterns []map[string]string
			for _, rule := range taskSystem.Regex {
				patterns = append(patterns, map[string]string{
					"pattern": rule.Pattern,
					"label":   rule.Label,
				})
			}

			// Match tasks against this commit
			matchedTasks := commit.MatchTask(patterns, splitPattern)
			if len(matchedTasks) > 0 {
				referenced = true
				tasks.Add(matchedTasks...)
			}
		}

		// If no task was matched, create an unreferenced task
		if !referenced {
			task := model.NewPACTask("")
			task.AddCommit(commit)
			tasks.Add(task)
		}
	}

	return tasks
}