package model

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
