// Package vcs provides version control system abstractions for PAC.
package vcs

import "github.com/Praqma/Praqmatic-Automated-Changelog/internal/model"

// VCS defines the interface for version control system operations.
// Implementations can provide Git, SVN, or other VCS backends.
type VCS interface {
	// GetDelta returns all commits between oldest and newest references.
	// If newest is empty, HEAD is used as the default.
	GetDelta(oldest, newest string) (*model.PACCommitCollection, error)

	// GetLatestTag returns the most recent tag matching the given pattern.
	// The pattern uses glob-style matching (e.g., "v*", "release-*").
	GetLatestTag(pattern string) (string, error)
}
