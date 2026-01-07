package vcs

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/config"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/model"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GitVCS implements the VCS interface for Git repositories.
type GitVCS struct {
	repo     *git.Repository
	settings config.VCSConfig
}

// NewGitVCS creates a new GitVCS instance for the given repository.
func NewGitVCS(settings config.VCSConfig) (*GitVCS, error) {
	repoPath := settings.RepoLocation
	if repoPath == "" {
		repoPath = "."
	}

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open git repository at %s: %w", repoPath, err)
	}

	return &GitVCS{
		repo:     repo,
		settings: settings,
	}, nil
}

// GetDelta returns all commits between oldest and newest references.
// If newest is empty, HEAD is used as the default.
// This includes all commits reachable from newest but not reachable from oldest,
// properly handling merge commits and all branches of history.
func (g *GitVCS) GetDelta(oldest, newest string) (*model.PACCommitCollection, error) {
	commits := model.NewPACCommitCollection()

	// Resolve the oldest reference
	oldHash, err := g.resolveRevision(oldest)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve oldest ref %q: %w", oldest, err)
	}

	// Resolve the newest reference (default to HEAD)
	var newHash *plumbing.Hash
	if newest == "" {
		head, err := g.repo.Head()
		if err != nil {
			return nil, fmt.Errorf("failed to get HEAD: %w", err)
		}
		hash := head.Hash()
		newHash = &hash
	} else {
		newHash, err = g.resolveRevision(newest)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve newest ref %q: %w", newest, err)
		}
	}

	// First, collect all commits reachable from the oldest reference.
	// These are the commits we want to EXCLUDE from our result.
	excludeSet := make(map[plumbing.Hash]bool)
	oldCommit, err := g.repo.CommitObject(*oldHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get oldest commit object: %w", err)
	}

	// Add the oldest commit itself to the exclude set
	excludeSet[*oldHash] = true

	// Walk all ancestors of the oldest commit
	oldIter := object.NewCommitIterCTime(oldCommit, nil, nil)
	defer oldIter.Close()
	err = oldIter.ForEach(func(c *object.Commit) error {
		excludeSet[c.Hash] = true
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error walking old commits: %w", err)
	}

	// Now walk all commits from newest and include those NOT in the exclude set
	newCommit, err := g.repo.CommitObject(*newHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get newest commit object: %w", err)
	}

	// Use CommitIterCTime to iterate in chronological order and include all branches
	newIter := object.NewCommitIterCTime(newCommit, nil, nil)
	defer newIter.Close()

	err = newIter.ForEach(func(c *object.Commit) error {
		// Skip if this commit is reachable from oldest
		if excludeSet[c.Hash] {
			return nil
		}

		// Apply path filtering if configured
		if len(g.settings.FilterPaths) > 0 {
			matches, err := g.commitMatchesFilter(c)
			if err != nil {
				return err
			}
			if !matches {
				return nil // Skip this commit
			}
		}

		// Convert to PACCommit
		pacCommit := g.convertCommit(c)
		commits.Add(pacCommit)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking commits: %w", err)
	}

	return commits, nil
}

// GetLatestTag returns the most recent tag matching the given glob pattern.
func (g *GitVCS) GetLatestTag(pattern string) (string, error) {
	tags, err := g.repo.Tags()
	if err != nil {
		return "", fmt.Errorf("failed to get tags: %w", err)
	}
	defer tags.Close()

	var latestTag string
	var latestTime time.Time

	err = tags.ForEach(func(ref *plumbing.Reference) error {
		tagName := ref.Name().Short()

		// Match against pattern
		matched, err := filepath.Match(pattern, tagName)
		if err != nil {
			return fmt.Errorf("invalid pattern %q: %w", pattern, err)
		}
		if !matched {
			return nil
		}

		// Get the commit time for this tag
		commitTime, err := g.getTagTime(ref)
		if err != nil {
			// Skip tags we can't resolve
			return nil
		}

		if latestTag == "" || commitTime.After(latestTime) {
			latestTag = tagName
			latestTime = commitTime
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if latestTag == "" {
		return "", fmt.Errorf("no tags found matching pattern: %s", pattern)
	}

	return latestTag, nil
}

// resolveRevision resolves a revision string to a hash.
// Handles tags, branches, and commit SHAs.
func (g *GitVCS) resolveRevision(rev string) (*plumbing.Hash, error) {
	// Try to resolve as a revision
	hash, err := g.repo.ResolveRevision(plumbing.Revision(rev))
	if err == nil {
		return hash, nil
	}

	// Try as a direct hash
	if len(rev) >= 7 {
		h := plumbing.NewHash(rev)
		if _, err := g.repo.CommitObject(h); err == nil {
			return &h, nil
		}
	}

	return nil, fmt.Errorf("cannot resolve revision: %s", rev)
}

// getTagTime returns the time associated with a tag.
// For annotated tags, returns the tagger time.
// For lightweight tags, returns the commit time.
func (g *GitVCS) getTagTime(ref *plumbing.Reference) (time.Time, error) {
	// Try as annotated tag first
	tagObj, err := g.repo.TagObject(ref.Hash())
	if err == nil {
		return tagObj.Tagger.When, nil
	}

	// Fall back to lightweight tag (points directly to commit)
	commit, err := g.repo.CommitObject(ref.Hash())
	if err != nil {
		return time.Time{}, fmt.Errorf("cannot get commit for tag: %w", err)
	}

	return commit.Committer.When, nil
}

// convertCommit converts a go-git commit to a PACCommit.
func (g *GitVCS) convertCommit(c *object.Commit) *model.PACCommit {
	message := strings.TrimSpace(c.Message)
	lines := strings.SplitN(message, "\n", 2)

	header := lines[0]
	body := ""
	if len(lines) > 1 {
		body = strings.TrimSpace(lines[1])
	}

	shortSHA := c.Hash.String()
	if len(shortSHA) > 7 {
		shortSHA = shortSHA[:7]
	}

	return &model.PACCommit{
		SHA:       c.Hash.String(),
		ShortSHA:  shortSHA,
		Message:   message,
		Header:    header,
		Body:      body,
		Timestamp: c.Committer.When,
	}
}

// commitMatchesFilter checks if a commit touches any of the configured filter paths.
func (g *GitVCS) commitMatchesFilter(c *object.Commit) (bool, error) {
	if len(g.settings.FilterPaths) == 0 {
		return true, nil
	}

	// Get the files changed in this commit
	stats, err := c.Stats()
	if err != nil {
		// If we can't get stats, include the commit to be safe
		return true, nil
	}

	for _, stat := range stats {
		for _, filterPath := range g.settings.FilterPaths {
			// Check if the file path matches the filter
			if matchesPath(stat.Name, filterPath) {
				return true, nil
			}
		}
	}

	return false, nil
}

// matchesPath checks if a file path matches a filter path pattern.
func matchesPath(filePath, filterPath string) bool {
	// Normalize paths
	filePath = filepath.Clean(filePath)
	filterPath = filepath.Clean(filterPath)

	// Check if file is under the filter path
	if strings.HasPrefix(filePath, filterPath) {
		return true
	}

	// Try glob matching
	matched, _ := filepath.Match(filterPath, filePath)
	return matched
}
