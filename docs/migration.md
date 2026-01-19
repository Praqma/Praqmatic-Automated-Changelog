# Migration Guide: Ruby to Go

This guide helps you migrate from the Ruby version of PAC to the Go version.

## Overview

The Go implementation is designed to be a **drop-in replacement** for the Ruby version. Your existing configuration files and templates should work without modification.

## What's the Same

### Configuration Files ✅

Your existing YAML configuration files work without changes:

```yaml
:general:
  :strict: false

:templates:
  - { location: templates/default.md, output: CHANGELOG.md }

:task_systems:
  - :name: none
    :regex:
      - { pattern: '/PAC-(\d+)/i', label: task }

:vcs:
  :type: git
  repo_location: '.'
```

### Templates ✅

All Liquid templates are fully compatible:

```liquid
# Changelog
{% for task in tasks.referenced %}
## {{ task.task_id }}
{% for commit in task.commits %}
- {{ commit.shortsha }}: {{ commit.header }}
{% endfor %}
{% endfor %}
```

### Available Template Variables

| Variable | Description |
|----------|-------------|
| `tasks.referenced` | Array of tasks with IDs |
| `tasks.unreferenced` | Array of commits without task references |
| `task.task_id` | The task identifier |
| `task.commits` | Array of commits for this task |
| `task.label` | Array of labels applied to this task |
| `task.data` | Data fetched from external system (Jira, etc.) |
| `commit.sha` | Full commit SHA |
| `commit.shortsha` / `commit.short_sha` | Short SHA (7 chars) |
| `commit.header` | First line of commit message |
| `commit.message` | Full commit message |
| `commit.body` | Commit message body (without header) |
| `commit.date` | Commit timestamp |
| `pac_c_count` | Total commit count |
| `pac_c_referenced` | Commits with task references |
| `pac_c_unreferenced` | Commits without task references |
| `pac_health` | Percentage of commits with task references |
| `properties` | Custom properties from config/CLI |

## What's Different

### CLI Syntax

The command structure is the same, with minor improvements:

```bash
# Ruby version
ruby pac.rb from v1.0.0 --settings settings.yml

# Go version
pac from v1.0.0 --settings settings.yml
```

### Flag Changes

| Ruby | Go | Notes |
|------|-----|-------|
| `-v` | `-v` | Same - increase verbosity |
| `-vv` | `-vv` | Same - more verbose |
| `-q` | `-q` | Same - quiet mode |
| `--settings` | `--settings` | Same |
| `--to` | `--to` | Same |
| `-c user pw target` | `-c user -c pw -c target` | Credentials as separate flags |

### Credentials Override

```bash
# Ruby version
pac from v1.0.0 -c myuser mypassword jira

# Go version
pac from v1.0.0 -c myuser -c mypassword -c jira
```

### New Features in Go Version

1. **Shell Completion**: Built-in completion for Bash, Zsh, Fish, PowerShell
2. **Cross-Platform Binaries**: Pre-built for Linux, macOS, Windows
3. **Smaller Image**: Docker image ~20MB vs ~200MB+ for Ruby
4. **Faster Execution**: 3-10x faster for large repositories
5. **No Ruby Required**: Single static binary with no dependencies

## Regex Compatibility

Go's regex engine is mostly compatible with Ruby, but there are some differences:

### Fully Supported ✅

```yaml
# Simple patterns
- { pattern: '/PAC-(\d+)/', label: task }

# Case-insensitive
- { pattern: '/jira-(\d+)/i', label: jira }

# Character classes
- { pattern: '/[A-Z]+-\d+/', label: task }

# Alternation
- { pattern: '/(BUG|FEAT|FIX)-(\d+)/', label: task }
```

### Not Supported ❌

Go's regexp package uses RE2, which does not support:

- Backreferences (`\1`, `\2`)
- Lookahead/Lookbehind (`(?=...)`, `(?<=...)`)
- Possessive quantifiers (`++`, `*+`)

If you use these features, you'll need to rewrite your patterns:

```yaml
# Ruby (not supported in Go)
- { pattern: '/(\w+)-\1/', label: task }

# Go alternative - use simpler pattern
- { pattern: '/(\w+)-(\w+)/', label: task }
```

## Docker Migration

### Ruby Docker Usage

```bash
docker run --rm \
  -v $(pwd):/repo \
  praqma/pac:latest \
  from v1.0.0
```

### Go Docker Usage

```bash
docker run --rm \
  -v $(pwd):/repo \
  ghcr.io/praqma/pac:latest \
  from v1.0.0
```

The Go image is available from GitHub Container Registry (`ghcr.io`).

## Troubleshooting Migration

### Template rendering differently

If your template output differs, check for:

1. **Whitespace differences**: Liquid engines may handle whitespace slightly differently
2. **Filter availability**: Ensure you're using standard Liquid filters

### Regex not matching

1. Check for unsupported regex features (see above)
2. Test your regex with Go's regexp tester
3. Ensure the pattern is properly escaped in YAML

### Jira/API authentication failing

The Go version uses the same authentication mechanism:

```yaml
task_systems:
  - name: jira
    query_string: 'https://jira.example.com/rest/api/2/issue/#{task_id}'
    usr: myuser
    pw: mypassword
```

Ensure your credentials are correctly specified. You can also use environment variables or the `-c` flag.

## Getting Help

If you encounter issues migrating:

1. Check the [GitHub Issues](https://github.com/Praqma/Praqmatic-Automated-Changelog/issues)
2. Run with `-vvv` for detailed debug output
3. Compare output with Ruby version on a small commit range
4. Open an issue with your configuration (sanitized) and error message
