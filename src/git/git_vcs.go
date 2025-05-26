package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"

	"github.com/Praqma/Praqmatic-Automated-Changelog/src/model"
)

// GitVCS implements version control operations for Git
type GitVCS struct {
	config model.VCS
	repo   *git.Repository
}

// NewGitVCS creates a new Git VCS instance
func NewGitVCS(config model.VCS) (*GitVCS, error) {
	// Open the repository
	cwd := config.Repo

	// Check if the path is absolute, relative or a URL
	if strings.HasPrefix(cwd, "http://") || strings.HasPrefix(cwd, "https://") {
		// if it is a URL, read the repo from the URL without cloning
		fmt.Println("Opening git repository from URL:", cwd)
		dirName := filepath.Base(cwd)
		repo, err := git.PlainClone(dirName, false, &git.CloneOptions{
			URL: cwd,
			Progress: os.Stdout,
			SingleBranch: true,
		})
		if err != nil && err.Error() == "repository already exists" {
			// If the repository already exists, open it and fetch the latest changes
			existingRepo, err := git.PlainOpen(dirName)
			if err == nil {
				fmt.Println("Repository already exists. Fetching latest changes...")
				err = existingRepo.Fetch(&git.FetchOptions{
					RemoteName: "origin",
				})
				if err != nil && err != git.NoErrAlreadyUpToDate {
					return nil, fmt.Errorf("error fetching updates for existing repository: %w", err)
				}
				return &GitVCS{
					config: config,
					repo:   existingRepo,
				}, nil
			}
		}
		if err != nil {
			return nil, fmt.Errorf("error cloning git repository from URL: %w", err)
		}
		return &GitVCS{
			config: config,
			repo:   repo,
		}, nil
	} else if !strings.HasPrefix(cwd, "/") {
		// If it's not absolute, make it relative to the current working directory
		absPath, err := filepath.Abs(cwd)
		if err != nil {
			return nil, fmt.Errorf("error getting absolute path: %w", err)
		}
		cwd = absPath
	}
	
	repo, err := git.PlainOpen(cwd)
	if err != nil {
		return nil, fmt.Errorf("error opening git repository (use 'from -h' for help): %w", err)
	}

	return &GitVCS{
		config: config,
		repo:   repo,
	}, nil
}

// GetCommitsBetween retrieves all commits between two references
func (g *GitVCS) GetCommitsBetween(from, to string) (*model.PACCommitCollection, error) {
	commits := model.NewPACCommitCollection()
	fromHash, err := g.repo.ResolveRevision(plumbing.Revision(from))
	if err != nil {
		return nil, fmt.Errorf("error resolving 'from' reference: %w", err)
	}
	// Resolve the 'to' reference
	toHash, err := g.repo.ResolveRevision(plumbing.Revision(to))
	if err != nil {
		return nil, fmt.Errorf("error resolving 'to' reference: %w", err)
	}
	// Create a commit iterator for the range
	commitIter, err := g.repo.Log(&git.LogOptions{
		From:  *toHash,
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating commit iterator: %w", err)
	}

	defer commitIter.Close()

	err = commitIter.ForEach(func(commit *object.Commit) error {
		if commit.Hash == *fromHash {
			return fmt.Errorf("break")
		}
		commits.Add(model.NewPACCommit(commit.Hash.String(), commit.Message, commit.Author.When))
		return nil
	})

	if err != nil && err.Error() != "break" {
		return nil, fmt.Errorf("error iterating commits: %w", err)
	}

	return commits, nil
}
