package model

import (
	"regexp"
	"strings"
	"time"
)

// PACCommit represents a single commit in the version control system
type PACCommit struct {
	SHA        string
	Message    string
	Date       time.Time
	Referenced bool
}

// NewPACCommit creates a new commit instance
func NewPACCommit(sha, message string, date time.Time) *PACCommit {
	return &PACCommit{
		SHA:        sha,
		Message:    message,
		Date:       date,
		Referenced: false,
	}
}

// Header returns the first line of the commit message
func (c *PACCommit) Header() string {
	if len(c.Message) == 0 {
		return ""
	}
	return strings.Split(c.Message, "\n")[0]
}

// ShortSHA returns an abbreviated SHA
func (c *PACCommit) ShortSHA() string {
	n := 7 // Default length for short SHA
	if len(c.SHA) <= n {
		return c.SHA
	}
	return c.SHA[0:n]
}

// MatchTask matches tasks against this commit
// Returns an array of matched tasks
func (c *PACCommit) MatchTask(patterns []map[string]string, splitPattern string) []*PACTask {
	var tasks []*PACTask

	for _, pattern := range patterns {
		// In Go, we use the regexp package instead of Ruby's eval
		re, err := compilePattern(pattern["pattern"])
		if err != nil {
			continue
		}

		// Check the entire message for matches, not just the first line
		matches := re.FindAllStringSubmatch(c.Message, -1)
		for _, match := range matches {
			if len(match) > 1 {
				if splitPattern != "" {
					// Split the task ID by the provided pattern
					taskIDs := strings.Split(match[1], splitPattern)
					for _, id := range taskIDs {
						id = strings.TrimSpace(id)
						if id == "" {
							continue
						}
						task := NewPACTask(id)
						task.AddCommit(c)
						task.AddLabel(pattern["label"])
						c.Referenced = true
						tasks = append(tasks, task)
					}
				} else {
					id := strings.TrimSpace(match[1])
					if id == "" {
						continue
					}
					task := NewPACTask(id)
					task.AddCommit(c)
					task.AddLabel(pattern["label"])
					c.Referenced = true
					tasks = append(tasks, task)
				}
			}
		}
	}

	return tasks
}

// Helper function to compile regex pattern
func compilePattern(pattern string) (*regexp.Regexp, error) {
	// Handle pattern string - in Go we need to convert from Ruby's regex format
	// This is a simplified conversion that handles common regex formats
	patternStr := pattern

	// If the pattern is enclosed in /.../ syntax (Ruby style), strip it
	if strings.HasPrefix(pattern, "/") && strings.HasSuffix(pattern, "/") {
		patternStr = strings.TrimPrefix(strings.TrimSuffix(pattern, "/"), "/")
	}

	// Ruby regex patterns often need some conversion for Go
	// For example, "#(\\d+)" in Ruby becomes "#(\\d+)" in Go
	// (No changes needed in this case since we're just matching issue numbers)

	return regexp.Compile(patternStr)
}

// PACCommitCollection represents a collection of commits
type PACCommitCollection struct {
	Commits []*PACCommit
}

// NewPACCommitCollection creates a new empty commit collection
func NewPACCommitCollection() *PACCommitCollection {
	return &PACCommitCollection{
		Commits: make([]*PACCommit, 0),
	}
}

// Add adds a commit or commits to the collection
func (cc *PACCommitCollection) Add(commits ...*PACCommit) {
	cc.Commits = append(cc.Commits, commits...)
}

// Count returns the total number of commits
func (cc *PACCommitCollection) Count() int {
	return len(cc.Commits)
}

// CountReferenced returns the number of commits that are referenced by tasks
func (cc *PACCommitCollection) CountReferenced() int {
	count := 0
	for _, commit := range cc.Commits {
		if commit.Referenced {
			count++
		}
	}
	return count
}

// CountUnreferenced returns the number of commits that are not referenced by tasks
func (cc *PACCommitCollection) CountUnreferenced() int {
	return cc.Count() - cc.CountReferenced()
}

// Health returns the percentage of commits that are referenced
func (cc *PACCommitCollection) Health() float64 {
	if cc.Count() == 0 {
		return 0
	}
	return float64(cc.CountReferenced()) / float64(cc.Count())
}
