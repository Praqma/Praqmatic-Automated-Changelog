# Developer Guide

This guide is for developers who want to contribute to PAC or extend its functionality.

## Prerequisites

- Go 1.21 or later
- Git
- Make (optional, for using Makefile)
- Docker (optional, for container builds)

## Getting Started

### Clone and Build

```bash
git clone https://github.com/Praqma/Praqmatic-Automated-Changelog.git
cd Praqmatic-Automated-Changelog

# Download dependencies
go mod download

# Build
make build

# Run tests
make test
```

### Project Structure

```
.
├── cmd/
│   └── pac/
│       └── main.go              # Entry point
├── internal/
│   ├── cmd/                     # Cobra CLI commands
│   │   ├── root.go              # Root command and global flags
│   │   ├── from.go              # 'from' subcommand
│   │   ├── from_latest_tag.go   # 'from-latest-tag' subcommand
│   │   └── version.go           # 'version' subcommand
│   ├── config/                  # Configuration loading
│   │   ├── config.go            # Viper configuration
│   │   └── settings.go          # Settings structs
│   ├── core/                    # Core workflow
│   │   └── workflow.go          # Main PAC orchestration
│   ├── logging/                 # Logging utilities
│   │   └── logger.go            # Verbosity-based logger
│   ├── model/                   # Data models
│   │   ├── commit.go            # PACCommit, PACCommitCollection
│   │   ├── task.go              # PACTask, PACTaskCollection
│   │   └── liquid_drops.go      # Liquid template conversion
│   ├── report/                  # Report generation
│   │   ├── generator.go         # Report orchestration
│   │   └── renderer.go          # Liquid template rendering
│   ├── task/                    # Task systems
│   │   ├── system.go            # TaskSystem interface
│   │   ├── factory.go           # TaskSystem factory
│   │   ├── none.go              # NoneTaskSystem
│   │   ├── jira.go              # JiraTaskSystem
│   │   ├── regex.go             # Regex utilities
│   │   └── extractor.go         # Task ID extraction
│   └── vcs/                     # Version control
│       ├── interfaces.go        # VCS interface
│       └── git.go               # Git implementation
├── templates/                   # Default templates
├── settings/                    # Default settings
├── test/                        # Integration tests
├── go.mod
├── go.sum
├── Makefile
└── .goreleaser.yml
```

## Development Workflow

### Running Tests

```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Run tests with race detector
make test-race

# Run specific package tests
go test -v ./internal/model/...

# Run integration tests
make test-integration
```

### Linting

```bash
# Install golangci-lint (first time)
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run linter
make lint
```

### Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Build Docker image
make docker-build
```

## Adding a New Task System

To add support for a new task system (e.g., GitHub Issues):

### 1. Create the Task System File

**File: `internal/task/github.go`**

```go
package task

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strings"

    "github.com/Praqma/Praqmatic-Automated-Changelog/internal/config"
    "github.com/Praqma/Praqmatic-Automated-Changelog/internal/logging"
    "github.com/Praqma/Praqmatic-Automated-Changelog/internal/model"
)

// GitHubTaskSystem fetches issue data from GitHub API
type GitHubTaskSystem struct {
    config     config.TaskSystemConfig
    httpClient *http.Client
}

// NewGitHubTaskSystem creates a new GitHub task system
func NewGitHubTaskSystem(cfg config.TaskSystemConfig) *GitHubTaskSystem {
    return &GitHubTaskSystem{
        config:     cfg,
        httpClient: &http.Client{Timeout: 10 * time.Second},
    }
}

// Name returns the task system name
func (g *GitHubTaskSystem) Name() string {
    return g.config.Name
}

// Apply fetches GitHub issue data for each task
func (g *GitHubTaskSystem) Apply(tasks *model.PACTaskCollection) error {
    for _, task := range tasks.Tasks {
        if !task.AppliesTo[g.config.Name] || task.TaskID == "" {
            continue
        }

        data, err := g.fetchIssue(task.TaskID)
        if err != nil {
            logging.Warn("GitHub error for %s: %v", task.TaskID, err)
            task.ClearLabels()
            task.AddLabel("unknown")
            continue
        }

        task.Data = data
        logging.Debug("Applied GitHub data to %s", task.TaskID)
    }
    return nil
}

func (g *GitHubTaskSystem) fetchIssue(issueID string) (map[string]interface{}, error) {
    url := strings.ReplaceAll(g.config.QueryString, "#{task_id}", issueID)
    
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Accept", "application/vnd.github.v3+json")
    if g.config.Password != "" {
        req.Header.Set("Authorization", "token "+g.config.Password)
    }

    resp, err := g.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var data map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
        return nil, err
    }
    
    return data, nil
}
```

### 2. Register in Factory

**File: `internal/task/factory.go`**

```go
func CreateTaskSystem(cfg config.TaskSystemConfig) (TaskSystem, error) {
    switch cfg.Name {
    case "none":
        return NewNoneTaskSystem(cfg), nil
    case "jira":
        return NewJiraTaskSystem(cfg), nil
    case "github":
        return NewGitHubTaskSystem(cfg), nil  // Add this
    default:
        return nil, fmt.Errorf("unknown task system: %s", cfg.Name)
    }
}
```

### 3. Add Tests

**File: `internal/task/github_test.go`**

```go
func TestGitHubTaskSystem_Apply(t *testing.T) {
    // Create mock HTTP server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(map[string]interface{}{
            "title": "Test Issue",
            "state": "open",
        })
    }))
    defer server.Close()

    cfg := config.TaskSystemConfig{
        Name:        "github",
        QueryString: server.URL + "/issues/#{task_id}",
    }

    ts := NewGitHubTaskSystem(cfg)
    tasks := model.NewPACTaskCollection()
    task := tasks.FindOrCreate("123")
    task.AppliesTo["github"] = true

    err := ts.Apply(tasks)
    assert.NoError(t, err)
    assert.NotNil(t, task.Data)
}
```

### 4. Document Usage

Add usage example to documentation:

```yaml
task_systems:
  - name: github
    query_string: 'https://api.github.com/repos/owner/repo/issues/#{task_id}'
    pw: 'ghp_your_token_here'
    regex:
      - { pattern: '/#(\d+)/', label: issue }
```

## Adding Custom Liquid Filters

To add custom Liquid filters:

**File: `internal/report/filters.go`**

```go
package report

import (
    "strings"
    "github.com/osteele/liquid"
)

func registerCustomFilters(engine *liquid.Engine) {
    engine.RegisterFilter("capitalize_words", func(input string) string {
        return strings.Title(input)
    })
    
    engine.RegisterFilter("truncate_sha", func(sha string) string {
        if len(sha) > 7 {
            return sha[:7]
        }
        return sha
    })
}
```

Then in `renderer.go`:

```go
func NewRenderer() *Renderer {
    engine := liquid.NewEngine()
    registerCustomFilters(engine)
    return &Renderer{engine: engine}
}
```

## Code Style

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` for formatting
- Use `golangci-lint` for static analysis
- Write table-driven tests
- Document exported functions and types

## Pull Request Guidelines

1. **Branch naming**: `feature/description` or `fix/description`
2. **Commits**: Use conventional commits (`feat:`, `fix:`, `docs:`, etc.)
3. **Tests**: Add tests for new functionality
4. **Documentation**: Update docs for user-facing changes
5. **CI**: Ensure all checks pass

## Debugging

### Enable verbose logging

```bash
pac from v1.0.0 -vvv
```

### Debug with Delve

```bash
# Install delve
go install github.com/go-delve/delve/cmd/dlv@latest

# Debug
dlv debug ./cmd/pac -- from v1.0.0
```

### Profile performance

```bash
# CPU profile
go test -cpuprofile=cpu.prof -bench=. ./...
go tool pprof cpu.prof

# Memory profile
go test -memprofile=mem.prof -bench=. ./...
go tool pprof mem.prof
```

## Release Process

Releases are automated via GoReleaser when a tag is pushed:

```bash
# Create and push a tag
git tag -a v4.1.0 -m "Release v4.1.0"
git push origin v4.1.0
```

This triggers GitHub Actions to:
1. Run tests on all platforms
2. Build binaries for all platforms
3. Create Docker images
4. Publish GitHub release with assets
