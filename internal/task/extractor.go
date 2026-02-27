package task

import (
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/config"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/model"
)

// BuildTaskCollection processes commits and extracts tasks using configured task systems.
// Returns a PACTaskCollection with tasks grouped by task ID and commits marked as referenced.
func BuildTaskCollection(commits *model.PACCommitCollection, taskSystems []config.TaskSystemConfig) *model.PACTaskCollection {
	tasks := model.NewPACTaskCollection()

	for _, commit := range commits.Commits {
		matched := false

		// Try each task system's regex patterns
		for _, tsCfg := range taskSystems {
			matches := ExtractTaskIDs(commit.Message, tsCfg.Regex)

			for _, match := range matches {
				matched = true

				// Find or create task
				task := tasks.FindOrCreate(match.TaskID)
				task.AddCommit(commit)
				task.AddLabel(match.Label)
				task.AppliesTo[tsCfg.Name] = true
			}
		}

		// Mark commit as referenced if any task was found
		if matched {
			commit.MarkReferenced()
		} else {
			// Add to unreferenced task (empty task ID)
			unrefTask := tasks.FindOrCreate("")
			unrefTask.AddCommit(commit)
		}
	}

	return tasks
}

// UpdateCommitReferenceCounts updates the commit collection's reference counts
// based on the task collection. This should be called after task extraction.
func UpdateCommitReferenceCounts(commits *model.PACCommitCollection, tasks *model.PACTaskCollection) {
	// Create a set of referenced commit SHAs
	referenced := make(map[string]bool)

	for _, task := range tasks.Tasks {
		if task.TaskID != "" { // Skip unreferenced task
			for _, commit := range task.Commits.Commits {
				referenced[commit.SHA] = true
			}
		}
	}

	// Update commits
	for _, commit := range commits.Commits {
		commit.Referenced = referenced[commit.SHA]
	}
}
