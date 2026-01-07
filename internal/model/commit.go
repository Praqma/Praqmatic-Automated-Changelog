// Package model provides core data structures for PAC.package model

package model

import "time"

// PACCommit represents a single version control commit.
type PACCommit struct {
	SHA       string    // Full SHA hash of the commit
	ShortSHA  string    // Short SHA (typically 7 characters)
	Message   string    // Full commit message
	Header    string    // First line of commit message
	Body      string    // Rest of commit message after the first line
	Timestamp time.Time // Commit timestamp
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
