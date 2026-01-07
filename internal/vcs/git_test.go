package vcs

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/config"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// createTestRepo creates a temporary git repository with test commits.
func createTestRepo(t *testing.T) (string, func()) {
	t.Helper()

	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "pac-test-repo-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	// Initialize git repo
	repo, err := git.PlainInit(tmpDir, false)
	if err != nil {
		cleanup()
		t.Fatalf("failed to init repo: %v", err)
	}

	// Get worktree
	worktree, err := repo.Worktree()
	if err != nil {
		cleanup()
		t.Fatalf("failed to get worktree: %v", err)
	}

	// Create initial file
	testFile := filepath.Join(tmpDir, "README.md")
	if err := os.WriteFile(testFile, []byte("# Test\n"), 0644); err != nil {
		cleanup()
		t.Fatalf("failed to write file: %v", err)
	}

	// Stage and commit
	if _, err := worktree.Add("README.md"); err != nil {
		cleanup()
		t.Fatalf("failed to add file: %v", err)
	}

	sig := &object.Signature{
		Name:  "Test User",
		Email: "test@example.com",
		When:  time.Now().Add(-3 * time.Hour),
	}

	_, err = worktree.Commit("Initial commit", &git.CommitOptions{
		Author: sig,
	})
	if err != nil {
		cleanup()
		t.Fatalf("failed to create initial commit: %v", err)
	}

	// Create more commits
	for i := 1; i <= 3; i++ {
		content := []byte("# Test\n\nCommit " + string(rune('0'+i)) + "\n")
		if err := os.WriteFile(testFile, content, 0644); err != nil {
			cleanup()
			t.Fatalf("failed to write file: %v", err)
		}

		if _, err := worktree.Add("README.md"); err != nil {
			cleanup()
			t.Fatalf("failed to add file: %v", err)
		}

		sig := &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  time.Now().Add(-time.Duration(3-i) * time.Hour),
		}

		_, err = worktree.Commit("Commit "+string(rune('0'+i))+"\n\nBody of commit "+string(rune('0'+i)), &git.CommitOptions{
			Author: sig,
		})
		if err != nil {
			cleanup()
			t.Fatalf("failed to create commit: %v", err)
		}
	}

	return tmpDir, cleanup
}

// createTestRepoWithTags creates a test repo with tags.
func createTestRepoWithTags(t *testing.T) (string, func()) {
	t.Helper()

	tmpDir, cleanup := createTestRepo(t)

	repo, err := git.PlainOpen(tmpDir)
	if err != nil {
		cleanup()
		t.Fatalf("failed to open repo: %v", err)
	}

	// Get commit history
	iter, err := repo.Log(&git.LogOptions{})
	if err != nil {
		cleanup()
		t.Fatalf("failed to get log: %v", err)
	}

	var commits []*object.Commit
	iter.ForEach(func(c *object.Commit) error {
		commits = append(commits, c)
		return nil
	})

	// Create tags on different commits
	if len(commits) >= 2 {
		// Create lightweight tag on older commit
		ref := plumbing.NewHashReference(plumbing.NewTagReferenceName("v1.0.0"), commits[len(commits)-2].Hash)
		if err := repo.Storer.SetReference(ref); err != nil {
			cleanup()
			t.Fatalf("failed to create tag: %v", err)
		}

		// Create another tag on newest commit
		ref2 := plumbing.NewHashReference(plumbing.NewTagReferenceName("v1.1.0"), commits[0].Hash)
		if err := repo.Storer.SetReference(ref2); err != nil {
			cleanup()
			t.Fatalf("failed to create tag: %v", err)
		}

		// Create a non-matching tag
		ref3 := plumbing.NewHashReference(plumbing.NewTagReferenceName("release-2.0"), commits[0].Hash)
		if err := repo.Storer.SetReference(ref3); err != nil {
			cleanup()
			t.Fatalf("failed to create tag: %v", err)
		}
	}

	return tmpDir, cleanup
}

func TestNewGitVCS(t *testing.T) {
	tmpDir, cleanup := createTestRepo(t)
	defer cleanup()

	cfg := config.VCSConfig{
		Type:         "git",
		RepoLocation: tmpDir,
	}

	vcs, err := NewGitVCS(cfg)
	if err != nil {
		t.Fatalf("failed to create GitVCS: %v", err)
	}

	if vcs.repo == nil {
		t.Error("expected repo to be initialized")
	}
}

func TestNewGitVCS_InvalidPath(t *testing.T) {
	cfg := config.VCSConfig{
		Type:         "git",
		RepoLocation: "/nonexistent/path",
	}

	_, err := NewGitVCS(cfg)
	if err == nil {
		t.Error("expected error for invalid path")
	}
}

func TestGitVCS_GetDelta(t *testing.T) {
	tmpDir, cleanup := createTestRepo(t)
	defer cleanup()

	cfg := config.VCSConfig{
		Type:         "git",
		RepoLocation: tmpDir,
	}

	vcs, err := NewGitVCS(cfg)
	if err != nil {
		t.Fatalf("failed to create GitVCS: %v", err)
	}

	// Get all commits
	repo, _ := git.PlainOpen(tmpDir)
	iter, _ := repo.Log(&git.LogOptions{})
	var allCommits []*object.Commit
	iter.ForEach(func(c *object.Commit) error {
		allCommits = append(allCommits, c)
		return nil
	})

	if len(allCommits) < 2 {
		t.Fatal("expected at least 2 commits in test repo")
	}

	// Get delta between oldest and newest
	oldestSHA := allCommits[len(allCommits)-1].Hash.String()
	newestSHA := allCommits[0].Hash.String()

	commits, err := vcs.GetDelta(oldestSHA, newestSHA)
	if err != nil {
		t.Fatalf("GetDelta failed: %v", err)
	}

	// Should get all commits except the oldest (exclusive)
	expectedCount := len(allCommits) - 1
	if commits.Count() != expectedCount {
		t.Errorf("expected %d commits, got %d", expectedCount, commits.Count())
	}

	// Verify commit structure
	if commits.Count() > 0 {
		first := commits.Commits[0]
		if first.SHA == "" {
			t.Error("expected SHA to be set")
		}
		if first.ShortSHA == "" {
			t.Error("expected ShortSHA to be set")
		}
		if first.Header == "" {
			t.Error("expected Header to be set")
		}
	}
}

func TestGitVCS_GetDelta_DefaultToHead(t *testing.T) {
	tmpDir, cleanup := createTestRepo(t)
	defer cleanup()

	cfg := config.VCSConfig{
		Type:         "git",
		RepoLocation: tmpDir,
	}

	vcs, err := NewGitVCS(cfg)
	if err != nil {
		t.Fatalf("failed to create GitVCS: %v", err)
	}

	// Get oldest commit
	repo, _ := git.PlainOpen(tmpDir)
	iter, _ := repo.Log(&git.LogOptions{})
	var allCommits []*object.Commit
	iter.ForEach(func(c *object.Commit) error {
		allCommits = append(allCommits, c)
		return nil
	})

	oldestSHA := allCommits[len(allCommits)-1].Hash.String()

	// Get delta with empty newest (should use HEAD)
	commits, err := vcs.GetDelta(oldestSHA, "")
	if err != nil {
		t.Fatalf("GetDelta failed: %v", err)
	}

	expectedCount := len(allCommits) - 1
	if commits.Count() != expectedCount {
		t.Errorf("expected %d commits, got %d", expectedCount, commits.Count())
	}
}

func TestGitVCS_GetLatestTag(t *testing.T) {
	tmpDir, cleanup := createTestRepoWithTags(t)
	defer cleanup()

	cfg := config.VCSConfig{
		Type:         "git",
		RepoLocation: tmpDir,
	}

	vcs, err := NewGitVCS(cfg)
	if err != nil {
		t.Fatalf("failed to create GitVCS: %v", err)
	}

	// Find latest v* tag
	tag, err := vcs.GetLatestTag("v*")
	if err != nil {
		t.Fatalf("GetLatestTag failed: %v", err)
	}

	if tag != "v1.1.0" {
		t.Errorf("expected latest tag 'v1.1.0', got %q", tag)
	}
}

func TestGitVCS_GetLatestTag_NoMatch(t *testing.T) {
	tmpDir, cleanup := createTestRepoWithTags(t)
	defer cleanup()

	cfg := config.VCSConfig{
		Type:         "git",
		RepoLocation: tmpDir,
	}

	vcs, err := NewGitVCS(cfg)
	if err != nil {
		t.Fatalf("failed to create GitVCS: %v", err)
	}

	// Try to find a tag that doesn't exist
	_, err = vcs.GetLatestTag("nonexistent-*")
	if err == nil {
		t.Error("expected error for non-matching pattern")
	}
}

func TestMatchesPath(t *testing.T) {
	tests := []struct {
		name       string
		filePath   string
		filterPath string
		expected   bool
	}{
		{
			name:       "exact match",
			filePath:   "src/main.go",
			filterPath: "src/main.go",
			expected:   true,
		},
		{
			name:       "prefix match",
			filePath:   "src/pkg/util.go",
			filterPath: "src",
			expected:   true,
		},
		{
			name:       "no match",
			filePath:   "docs/readme.md",
			filterPath: "src",
			expected:   false,
		},
		{
			name:       "glob match",
			filePath:   "main.go",
			filterPath: "*.go",
			expected:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchesPath(tt.filePath, tt.filterPath)
			if result != tt.expected {
				t.Errorf("matchesPath(%q, %q) = %v, want %v",
					tt.filePath, tt.filterPath, result, tt.expected)
			}
		})
	}
}

func TestConvertCommit(t *testing.T) {
	tmpDir, cleanup := createTestRepo(t)
	defer cleanup()

	cfg := config.VCSConfig{
		Type:         "git",
		RepoLocation: tmpDir,
	}

	vcs, err := NewGitVCS(cfg)
	if err != nil {
		t.Fatalf("failed to create GitVCS: %v", err)
	}

	// Get a commit
	repo, _ := git.PlainOpen(tmpDir)
	head, _ := repo.Head()
	commit, _ := repo.CommitObject(head.Hash())

	pacCommit := vcs.convertCommit(commit)

	if pacCommit.SHA != commit.Hash.String() {
		t.Errorf("SHA mismatch: got %q, want %q", pacCommit.SHA, commit.Hash.String())
	}

	if len(pacCommit.ShortSHA) != 7 {
		t.Errorf("ShortSHA should be 7 chars, got %d", len(pacCommit.ShortSHA))
	}

	if pacCommit.Header == "" {
		t.Error("expected Header to be set")
	}

	if pacCommit.Timestamp.IsZero() {
		t.Error("expected Timestamp to be set")
	}
}
