# PAC Go Implementation

This directory contains documentation specific to the Go implementation of PAC (Praqmatic Automated Changelog).

## Documentation Index

| Document | Description |
|----------|-------------|
| [Installation](installation.md) | How to install and run PAC |
| [Migration Guide](migration.md) | Migrating from Ruby to Go version |
| [Development](development.md) | Developer guide for contributing |
| [Architecture](architecture.md) | System architecture and design |

## Quick Start

### Installation

```bash
# Using Go
go install github.com/Praqma/Praqmatic-Automated-Changelog/cmd/pac@latest

# Using Docker
docker pull ghcr.io/praqma/pac:latest

# From source
git clone https://github.com/Praqma/Praqmatic-Automated-Changelog.git
cd Praqmatic-Automated-Changelog
make build
./bin/pac --help
```

### Basic Usage

```bash
# Generate changelog from a tag to HEAD
pac from v1.0.0 --settings pac_settings.yml

# Generate changelog from latest matching tag
pac from-latest-tag "v*" --settings pac_settings.yml

# Increase verbosity
pac from v1.0.0 -vvv

# Override credentials
pac from v1.0.0 -c username password jira
```

## Compatibility

The Go implementation is **100% compatible** with existing:
- Configuration files (YAML settings)
- Liquid templates
- Git repositories

No changes to your existing templates or settings files are required.
