package model

// ToLiquid converts a PACCommit to a map for Liquid template rendering.
// The keys match the Ruby implementation for template compatibility.
func (c *PACCommit) ToLiquid() map[string]any {
	return map[string]any{
		"sha":       c.SHA,
		"shortsha":  c.ShortSHA,
		"short_sha": c.ShortSHA, // Alias for Ruby compatibility
		"header":    c.Header,
		"message":   c.Message,
		"body":      c.Body,
		"time":      c.Timestamp,
	}
}

// ToLiquid converts a PACTask to a map for Liquid template rendering.
// The keys match the Ruby implementation for template compatibility.
func (t *PACTask) ToLiquid() map[string]any {
	commits := make([]map[string]any, 0, len(t.Commits.Commits))
	for _, c := range t.Commits.Commits {
		commits = append(commits, c.ToLiquid())
	}

	return map[string]any{
		"task_id":    t.TaskID,
		"commits":    commits,
		"attributes": t.Attributes,
		"label":      t.GetLabels(),
		"data":       t.Data,
	}
}

// ToLiquid converts a PACTaskCollection to a map for Liquid template rendering.
// Returns a map with "referenced" (tasks with IDs) and "unreferenced" (commits without task refs).
func (tc *PACTaskCollection) ToLiquid() map[string]any {
	referenced := make([]map[string]any, 0)
	unreferenced := make([]map[string]any, 0)

	for _, task := range tc.Tasks {
		if task.TaskID == "" {
			// Unreferenced commits (stored in task with empty ID)
			for _, commit := range task.Commits.Commits {
				unreferenced = append(unreferenced, commit.ToLiquid())
			}
		} else {
			referenced = append(referenced, task.ToLiquid())
		}
	}

	return map[string]any{
		"referenced":   referenced,
		"unreferenced": unreferenced,
	}
}

// ToLiquidCommits converts a PACCommitCollection to a slice for Liquid templates.
func (c *PACCommitCollection) ToLiquidCommits() []map[string]any {
	result := make([]map[string]any, 0, len(c.Commits))
	for _, commit := range c.Commits {
		result = append(result, commit.ToLiquid())
	}
	return result
}
