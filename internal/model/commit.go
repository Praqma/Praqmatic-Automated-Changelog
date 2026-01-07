// Package model provides core data structures for PAC.
package model

import "time"

// PACCommit represents a single version control commit.
type PACCommit struct {
	SHA        string    // Full SHA hash of the commit
	ShortSHA   string    // Short SHA (typically 7 characters)
	Message    string    // Full commit message
	Header     string    // First line of commit message
	Body       string    // Rest of commit message after the first line
	Timestamp  time.Time // Commit timestamp
	Referenced bool      // Whether this commit references a task
}

// PACCommitCollection holds a collection of commits.
type PACCommitCollection struct {
	Commits []*PACCommit
}

// NewPACCommitCollection creates a new empty commit collection.
func NewPACCommitCollection() *PACCommitCollection {
	return &PACCommitCollection{
		Commits: []*PACCommit{},
	}
}

// Add appends a commit to the collection.
func (c *PACCommitCollection) Add(commit *PACCommit) {
	c.Commits = append(c.Commits, commit)
}

// Count returns the total number of commits.
func (c *PACCommitCollection) Count() int {
	return len(c.Commits)
}

// CountWith returns the number of commits that reference a task.
func (c *PACCommitCollection) CountWith() int {
	count := 0
	for _, commit := range c.Commits {
		if commit.Referenced {
			count++
		}
	}
	return count
}

// CountWithout returns the number of commits that don't reference a task.
func (c *PACCommitCollection) CountWithout() int {
	return c.Count() - c.CountWith()
}

// Health returns the percentage of commits that reference tasks.
// Returns 100.0 if there are no commits.
func (c *PACCommitCollection) Health() float64 {
	if c.Count() == 0 {
		return 100.0
	}
	return (float64(c.CountWith()) / float64(c.Count())) * 100.0
}

// MarkReferenced marks a commit as referenced by a task.
func (c *PACCommit) MarkReferenced() {
	c.Referenced = true
}
