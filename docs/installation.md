# Installation Guide

This guide covers installing and running the Go version of PAC.

## Requirements

- Git (for repository operations)
- Network access (if using Jira or other task systems)

## Installation Methods

### Pre-built Binaries

Download the latest release for your platform from the [GitHub Releases](https://github.com/Praqma/Praqmatic-Automated-Changelog/releases) page.

Available platforms:
| OS | Architecture | Binary |
|----|--------------|--------|
| Linux | x86_64 (amd64) | `pac-linux-amd64` |
| Linux | ARM64 | `pac-linux-arm64` |
| macOS | x86_64 (Intel) | `pac-darwin-amd64` |
| macOS | ARM64 (Apple Silicon) | `pac-darwin-arm64` |
| Windows | x86_64 | `pac-windows-amd64.exe` |

```bash
# Linux/macOS example
curl -LO https://github.com/Praqma/Praqmatic-Automated-Changelog/releases/latest/download/pac-linux-amd64
chmod +x pac-linux-amd64
sudo mv pac-linux-amd64 /usr/local/bin/pac

# Verify installation
pac --version
```

### Using Go Install

If you have Go 1.21+ installed:

```bash
go install github.com/Praqma/Praqmatic-Automated-Changelog/cmd/pac@latest
```

This installs the `pac` binary to your `$GOPATH/bin` directory.

### From Source

```bash
# Clone the repository
git clone https://github.com/Praqma/Praqmatic-Automated-Changelog.git
cd Praqmatic-Automated-Changelog

# Build
make build

# The binary is in ./bin/pac
./bin/pac --version

# Optionally install system-wide
sudo cp ./bin/pac /usr/local/bin/
```

### Docker

The Docker image includes everything needed to run PAC:

```bash
# Pull the image
docker pull ghcr.io/praqma/pac:latest

# Run PAC (mount your repository and config)
docker run --rm \
  -v $(pwd):/repo \
  -v $(pwd)/pac_settings.yml:/settings.yml \
  ghcr.io/praqma/pac:latest \
  from v1.0.0 --settings /settings.yml
```

#### Docker with custom templates

```bash
docker run --rm \
  -v $(pwd):/repo \
  -v $(pwd)/templates:/templates \
  -v $(pwd)/settings.yml:/settings.yml \
  ghcr.io/praqma/pac:latest \
  from v1.0.0 --settings /settings.yml
```

## Verification

After installation, verify PAC is working:

```bash
# Check version
pac --version

# View help
pac --help

# View command help
pac from --help
```

## Shell Completion

PAC supports shell completion for Bash, Zsh, Fish, and PowerShell:

```bash
# Bash
pac completion bash > /etc/bash_completion.d/pac

# Zsh
pac completion zsh > "${fpath[1]}/_pac"

# Fish
pac completion fish > ~/.config/fish/completions/pac.fish

# PowerShell
pac completion powershell > pac.ps1
```

## Troubleshooting

### "command not found: pac"

Ensure the binary is in your PATH:

```bash
# Check if pac is in PATH
which pac

# If using go install, ensure GOPATH/bin is in PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

### Permission denied

Make the binary executable:

```bash
chmod +x /path/to/pac
```

### Git repository not found

PAC needs to run from within a Git repository, or you need to specify the repository location in your settings:

```yaml
vcs:
  type: git
  repo_location: '/path/to/your/repo'
```

### TLS/SSL errors

If you're behind a corporate proxy or have certificate issues:

```bash
# Set custom CA certificate
export SSL_CERT_FILE=/path/to/ca-certificates.crt
```
