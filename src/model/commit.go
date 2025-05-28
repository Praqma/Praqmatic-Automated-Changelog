package model

import (
	"strings"
	"time"
)

// PACCommit represents a single commit in the version control system
type PACCommit struct {
	SHA        string
	Message    string
	Date       time.Time
	Author     string
}

// NewPACCommit creates a new commit instance
func NewPACCommit(sha, message string, date time.Time, author string) *PACCommit {
	return &PACCommit{
		SHA:        sha,
		Message:    message,
		Date:       date,
		Author:     author,
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