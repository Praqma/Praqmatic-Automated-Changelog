# Architecture

This document describes the architecture and design of the Go implementation of PAC.

## High-Level Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                           CLI (Cobra)                           │
│                    cmd/pac, internal/cmd                        │
└─────────────────────────────────────────────────────────────────┘
                                │
                                ▼
┌─────────────────────────────────────────────────────────────────┐
│                        Core Workflow                            │
│                      internal/core                              │
└─────────────────────────────────────────────────────────────────┘
        │                       │                       │
        ▼                       ▼                       ▼
┌───────────────┐   ┌───────────────────┐   ┌─────────────────────┐
│  VCS (Git)    │   │   Task Systems    │   │  Report Generator   │
│ internal/vcs  │   │  internal/task    │   │  internal/report    │
└───────────────┘   └───────────────────┘   └─────────────────────┘
        │                       │                       │
        ▼                       ▼                       ▼
┌─────────────────────────────────────────────────────────────────┐
│                         Model Layer                             │
│                       internal/model                            │
└─────────────────────────────────────────────────────────────────┘
```

## Component Details

### CLI Layer (`cmd/pac`, `internal/cmd`)

The CLI layer uses [Cobra](https://github.com/spf13/cobra) for command-line parsing.

**Key files:**
- `cmd/pac/main.go` - Entry point
- `internal/cmd/root.go` - Root command with global flags
- `internal/cmd/from.go` - `from` subcommand
- `internal/cmd/from_latest_tag.go` - `from-latest-tag` subcommand

**Responsibilities:**
- Parse command-line arguments
- Load configuration
- Invoke core workflow
- Handle errors and exit codes

### Configuration (`internal/config`)

Uses [Viper](https://github.com/spf13/viper) for configuration management.

**Key files:**
- `config.go` - Configuration loading and merging
- `settings.go` - Go structs matching YAML schema

**Features:**
- YAML file parsing
- Environment variable support
- Command-line flag overrides
- Credential injection

**Configuration flow:**
```
YAML File → Viper → Settings Struct
    ↑           ↓
CLI Flags → Overrides
```

### Core Workflow (`internal/core`)

Orchestrates the main PAC workflow.

**Key files:**
- `workflow.go` - Main workflow functions

**Workflow steps:**
1. Initialize VCS (Git)
2. Get commit delta between references
3. Build task collection from commits
4. Apply task systems (fetch external data)
5. Generate reports from templates

```go
func Run(settings, oldestRef, newestRef) error {
    vcs := NewGitVCS(settings.VCS)
    commits := vcs.GetDelta(oldestRef, newestRef)
    tasks := BuildTaskCollection(commits, settings.TaskSystems)
    
    for _, ts := range taskSystems {
        ts.Apply(tasks)
    }
    
    generator := NewGenerator(tasks, commits)
    return generator.Generate(settings)
}
```

### VCS Layer (`internal/vcs`)

Abstracts version control operations using [go-git](https://github.com/go-git/go-git).

**Key files:**
- `interfaces.go` - VCS interface definition
- `git.go` - Git implementation

**Interface:**
```go
type VCS interface {
    GetDelta(oldest, newest string) (*model.PACCommitCollection, error)
    GetLatestTag(pattern string) (string, error)
}
```

**Features:**
- Pure Go implementation (no system Git required)
- Commit range extraction
- Tag pattern matching
- Path filtering

### Model Layer (`internal/model`)

Core data structures used throughout PAC.

**Key files:**
- `commit.go` - PACCommit, PACCommitCollection
- `task.go` - PACTask, PACTaskCollection
- `liquid_drops.go` - Liquid template conversion

**Data model:**
```
PACCommitCollection
    └── []PACCommit
            ├── SHA
            ├── ShortSHA
            ├── Message
            ├── Header
            ├── Body
            └── Date

PACTaskCollection
    └── []PACTask
            ├── TaskID
            ├── Commits → PACCommitCollection
            ├── Labels
            ├── AppliesTo
            └── Data (from external systems)
```

**Liquid conversion:**
Models implement `ToLiquid()` method for template rendering:
```go
func (c *PACCommit) ToLiquid() map[string]interface{} {
    return map[string]interface{}{
        "sha":      c.SHA,
        "shortsha": c.ShortSHA,
        "header":   c.Header,
        // ...
    }
}
```

### Task Systems (`internal/task`)

Pluggable system for fetching external task data.

**Key files:**
- `system.go` - TaskSystem interface
- `factory.go` - Factory function
- `none.go` - NoneTaskSystem (regex-only)
- `jira.go` - JiraTaskSystem (HTTP/JSON)
- `regex.go` - Regex utilities
- `extractor.go` - Task ID extraction

**Interface:**
```go
type TaskSystem interface {
    Name() string
    Apply(tasks *model.PACTaskCollection) error
}
```

**Task extraction flow:**
```
Commits → Regex Matching → Task IDs → Task Collection
                                            ↓
External APIs (Jira, etc.) → Task Data → Enriched Tasks
```

### Report Generation (`internal/report`)

Generates output using [Liquid](https://github.com/osteele/liquid) templates.

**Key files:**
- `generator.go` - Report orchestration
- `renderer.go` - Template rendering

**Template data structure:**
```go
{
    "tasks": {
        "referenced": [...],    // Tasks with IDs
        "unreferenced": [...]   // Commits without tasks
    },
    "pac_c_count": 100,
    "pac_c_referenced": 85,
    "pac_c_unreferenced": 15,
    "pac_health": 85.0,
    "properties": {...}
}
```

### Logging (`internal/logging`)

Verbosity-based logging system.

**Verbosity levels:**
| Level | Flag | Output |
|-------|------|--------|
| -1 | `-q` | Errors only |
| 0 | (default) | Normal output |
| 1 | `-v` | Verbose |
| 2 | `-vv` | Debug |
| 3 | `-vvv` | Trace |

**Usage:**
```go
logging.Info("Processing %d commits", count)   // Level 0
logging.Debug("Commit SHA: %s", sha)           // Level 2
logging.Trace("HTTP response: %v", resp)       // Level 3
```

## Data Flow

```
                    ┌─────────────────┐
                    │  Git Repository │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │   VCS Layer     │
                    │  (go-git)       │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ PACCommit       │
                    │ Collection      │
                    └────────┬────────┘
                             │
         ┌───────────────────┼───────────────────┐
         │                   │                   │
         ▼                   ▼                   ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│  Regex Extract  │ │  Regex Extract  │ │  Regex Extract  │
│  (none system)  │ │  (jira system)  │ │ (other system)  │
└────────┬────────┘ └────────┬────────┘ └────────┬────────┘
         │                   │                   │
         └───────────────────┼───────────────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ PACTask         │
                    │ Collection      │
                    └────────┬────────┘
                             │
         ┌───────────────────┼───────────────────┐
         │                   │                   │
         ▼                   ▼                   ▼
┌─────────────────┐ ┌─────────────────┐ ┌─────────────────┐
│  NoneTaskSystem │ │ JiraTaskSystem  │ │ Other System    │
│  (no-op)        │ │ (HTTP fetch)    │ │ (HTTP fetch)    │
└────────┬────────┘ └────────┬────────┘ └────────┬────────┘
         │                   │                   │
         └───────────────────┼───────────────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ Enriched Tasks  │
                    │ (with data)     │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ Liquid Template │
                    │ Renderer        │
                    └────────┬────────┘
                             │
                             ▼
                    ┌─────────────────┐
                    │ Output Files    │
                    │ (MD, HTML)      │
                    └─────────────────┘
```

## Extension Points

### Adding a new VCS

Implement the `VCS` interface in `internal/vcs/interfaces.go`:

```go
type VCS interface {
    GetDelta(oldest, newest string) (*model.PACCommitCollection, error)
    GetLatestTag(pattern string) (string, error)
}
```

### Adding a new Task System

1. Implement `TaskSystem` interface
2. Register in `factory.go`
3. Add tests

### Adding Custom Liquid Filters

Register filters in `internal/report/renderer.go`:

```go
engine.RegisterFilter("my_filter", func(input string) string {
    // Transform input
    return transformed
})
```

## Testing Strategy

- **Unit tests**: Each package has `*_test.go` files
- **Integration tests**: `test/integration/` for end-to-end testing
- **Template tests**: `test/templates/` for Liquid compatibility

Run all tests:
```bash
make test
make test-integration
```
