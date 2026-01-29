# Release Process

This document describes how the PAC release process works using GoReleaser.

## Overview

PAC uses [GoReleaser](https://goreleaser.com/) to automate the entire build and release process. This includes:

- Multi-platform binary builds (Linux, macOS, Windows)
- Archive creation (tar.gz, zip)
- Linux packages (deb, rpm, apk)
- Docker images (multi-arch)
- Package manager integrations (Homebrew, Chocolatey, Winget)
- GitHub releases with changelogs

## Release Types

### 1. Tagged Releases (Production)

When a tag matching `v*` is pushed, the release workflow automatically:

1. Builds binaries for all platforms
2. Creates archives (tar.gz for Unix, zip for Windows)
3. Generates Linux packages (deb, rpm, apk)
4. Builds and pushes multi-arch Docker images
5. Creates a **draft** GitHub release
6. Updates package manager repositories (Homebrew tap, etc.)

**To create a release:**

```bash
# Tag the release
git tag -a v4.0.0 -m "Release v4.0.0"
git push origin v4.0.0

# The GitHub Action will automatically create a draft release
# Review the draft and publish when ready
```

### 2. Dev Builds (Development)

When code is pushed to the `main` branch, the dev workflow automatically:

1. Builds snapshot versions of all binaries
2. Updates the `dev` release with latest artifacts
3. Pushes Docker images with the `dev` tag

This allows users to test the latest development version.

### 3. Prereleases

Tags containing `-rc`, `-beta`, `-alpha` are automatically detected as prereleases:

```bash
git tag -a v4.0.0-rc1 -m "Release candidate 1"
git push origin v4.0.0-rc1
```

## Installation Methods

### Homebrew (macOS/Linux)

```bash
# Add the tap (first time only)
brew tap praqma/tap

# Install latest stable
brew install pac

# Install specific version
brew install pac@4.0

# Install development version
brew install praqma/tap/pac --HEAD
```

### Chocolatey (Windows)

```powershell
# Install latest stable
choco install pac

# Install specific version
choco install pac --version=4.0.0
```

### Winget (Windows)

```powershell
# Install latest stable
winget install Praqma.pac

# Install specific version
winget install Praqma.pac --version 4.0.0
```

### Debian/Ubuntu

```bash
# Download the .deb file from GitHub releases
wget https://github.com/Praqma/Praqmatic-Automated-Changelog/releases/download/v4.0.0/pac_4.0.0_amd64.deb

# Install
sudo dpkg -i pac_4.0.0_amd64.deb
```

### RPM-based distros (Fedora, RHEL, CentOS)

```bash
# Download the .rpm file from GitHub releases
wget https://github.com/Praqma/Praqmatic-Automated-Changelog/releases/download/v4.0.0/pac-4.0.0.x86_64.rpm

# Install
sudo rpm -i pac-4.0.0.x86_64.rpm
```

### Docker

```bash
# Latest stable
docker pull ghcr.io/praqma/pac:latest

# Specific version
docker pull ghcr.io/praqma/pac:4.0.0

# Dev (development)
docker pull ghcr.io/praqma/pac:dev

# Run
docker run --rm -v $(pwd):/repo ghcr.io/praqma/pac:latest from v1.0.0 to v2.0.0
```

### Direct Binary Download

Download the appropriate archive from the [GitHub Releases](https://github.com/Praqma/Praqmatic-Automated-Changelog/releases) page.

## Local Development

### Building Locally

```bash
# Build for current platform
make build

# Build snapshot for all platforms
make snapshot

# Build snapshot for current platform only (faster)
make snapshot-single
```

### Testing the Release Configuration

```bash
# Validate GoReleaser configuration
make release-check

# Dry-run release (builds everything but doesn't publish)
make snapshot
```

## Configuration

The release configuration is in [.goreleaser.yaml](../.goreleaser.yaml).

### Required Secrets

For the GitHub Actions workflows to work, you need to configure these secrets:

| Secret | Description | Required For |
|--------|-------------|--------------|
| `GITHUB_TOKEN` | Automatic | GitHub releases, Docker images (ghcr.io) |
| `HOMEBREW_TAP_GITHUB_TOKEN` | PAT with repo access | Homebrew tap updates |
| `CHOCOLATEY_API_KEY` | Chocolatey API key | Chocolatey publishing |
| `WINGET_GITHUB_TOKEN` | PAT with repo access | Winget manifest updates |

### External Repositories

You need to create these repositories for package managers:

1. **Homebrew Tap**: `Praqma/homebrew-tap`
3. **Winget Manifests**: `Praqma/winget-pkgs` (or submit to official winget-pkgs)

## Changelog Generation

GoReleaser automatically generates changelogs from commit messages. Follow these conventions:

- `feat:` - New features
- `fix:` - Bug fixes
- `docs:` - Documentation changes (excluded from changelog)
- `test:` - Test changes (excluded from changelog)
- `chore:` - Maintenance (excluded from changelog)
- `ci:` - CI changes (excluded from changelog)
- `!` suffix - Breaking changes (e.g., `feat!: breaking change`)

## Troubleshooting

### Docker Build Failures

If Docker builds fail with QEMU errors:

```bash
# Set up QEMU for cross-platform builds
docker run --privileged --rm tonistiigi/binfmt --install all
```

### GoReleaser Validation

```bash
# Check configuration
goreleaser check

# Run with verbose logging
goreleaser release --snapshot --clean --verbose
```

### Version Issues

GoReleaser uses git tags for versioning. Ensure:

1. Tags follow semantic versioning: `v1.2.3`
2. Tags are annotated: `git tag -a v1.2.3 -m "Release v1.2.3"`
3. Full git history is available: `git fetch --unshallow` (in CI)
