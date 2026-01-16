# PAC (Praqmatic Automated Changelog) Makefile
# Simplified Makefile that uses GoReleaser for builds and releases

BINARY_NAME := pac
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR := ./bin
DIST_DIR := ./dist
GO := go
GOFLAGS := -trimpath
LDFLAGS := -ldflags "-s -w -X main.Version=$(VERSION)"

# Default target
.PHONY: all
all: build

# ============================================================================
# Development Targets
# ============================================================================

# Build for current platform (fast, for development)
.PHONY: build
build:
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/pac

# Run the application
.PHONY: run
run: build
	$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

# Install to GOPATH/bin
.PHONY: install
install:
	$(GO) install $(GOFLAGS) $(LDFLAGS) ./cmd/pac

# ============================================================================
# Testing Targets
# ============================================================================

# Run all tests
.PHONY: test
test:
	$(GO) test -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GO) test -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

# Run tests with race detector
.PHONY: test-race
test-race:
	$(GO) test -race ./...

# Run integration tests
.PHONY: test-integration
test-integration:
	$(GO) test -v -tags=integration ./test/...

# ============================================================================
# GoReleaser Targets (Primary build/release method)
# ============================================================================

# Check GoReleaser configuration
.PHONY: release-check
release-check:
	@which goreleaser > /dev/null || (echo "Installing goreleaser..." && go install github.com/goreleaser/goreleaser@latest)
	goreleaser check

# Build snapshot (all platforms, no publish)
.PHONY: snapshot
snapshot:
	@which goreleaser > /dev/null || (echo "Installing goreleaser..." && go install github.com/goreleaser/goreleaser@latest)
	goreleaser release --snapshot --clean

# Build for a single platform (fast local testing)
.PHONY: snapshot-single
snapshot-single:
	@which goreleaser > /dev/null || (echo "Installing goreleaser..." && go install github.com/goreleaser/goreleaser@latest)
	goreleaser build --snapshot --clean --single-target

# Full release (requires tag and GitHub token)
.PHONY: release
release:
	@which goreleaser > /dev/null || (echo "Installing goreleaser..." && go install github.com/goreleaser/goreleaser@latest)
	goreleaser release --clean

# ============================================================================
# Docker Targets
# ============================================================================

# Build Docker image for local testing
.PHONY: docker
docker:
	docker build -t $(BINARY_NAME):$(VERSION) -t $(BINARY_NAME):latest .

# Build multi-arch Docker image (requires buildx)
.PHONY: docker-multiarch
docker-multiarch:
	docker buildx build --platform linux/amd64,linux/arm64 \
		-t $(BINARY_NAME):$(VERSION) \
		-t $(BINARY_NAME):latest \
		--push .

# Run Docker container
.PHONY: docker-run
docker-run:
	docker run --rm -v $(PWD):/repo $(BINARY_NAME):latest $(ARGS)

# ============================================================================
# Code Quality Targets
# ============================================================================

# Run linter
.PHONY: lint
lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

# Format code
.PHONY: fmt
fmt:
	$(GO) fmt ./...
	@which goimports > /dev/null && goimports -w . || true

# Vet code
.PHONY: vet
vet:
	$(GO) vet ./...

# ============================================================================
# Dependency Management
# ============================================================================

# Check for outdated dependencies
.PHONY: deps-check
deps-check:
	$(GO) list -u -m all

# Update dependencies
.PHONY: deps-update
deps-update:
	$(GO) get -u ./...
	$(GO) mod tidy

# Download dependencies
.PHONY: deps
deps:
	$(GO) mod download

# ============================================================================
# Cleanup
# ============================================================================

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	rm -rf $(DIST_DIR)
	rm -f coverage.out coverage.html

# ============================================================================
# Utilities
# ============================================================================

# Show version
.PHONY: version
version:
	@echo $(VERSION)

# Generate (if any code generation is needed)
.PHONY: generate
generate:
	$(GO) generate ./...

# ============================================================================
# Help
# ============================================================================

.PHONY: help
help:
	@echo "PAC (Praqmatic Automated Changelog) Build System"
	@echo ""
	@echo "Development Targets:"
	@echo "  build           Build for current platform"
	@echo "  run             Build and run (use ARGS='...' for arguments)"
	@echo "  install         Install to GOPATH/bin"
	@echo ""
	@echo "Testing Targets:"
	@echo "  test            Run all tests"
	@echo "  test-coverage   Run tests with coverage report"
	@echo "  test-race       Run tests with race detector"
	@echo "  test-integration Run integration tests"
	@echo ""
	@echo "GoReleaser Targets (Recommended for releases):"
	@echo "  release-check   Validate GoReleaser configuration"
	@echo "  snapshot        Build snapshot for all platforms (no publish)"
	@echo "  snapshot-single Build snapshot for current platform only"
	@echo "  release         Full release (requires tag and GITHUB_TOKEN)"
	@echo ""
	@echo "Docker Targets:"
	@echo "  docker          Build Docker image for local testing"
	@echo "  docker-multiarch Build multi-arch Docker image"
	@echo "  docker-run      Run Docker container"
	@echo ""
	@echo "Code Quality:"
	@echo "  lint            Run golangci-lint"
	@echo "  fmt             Format code"
	@echo "  vet             Run go vet"
	@echo ""
	@echo "Dependency Management:"
	@echo "  deps            Download dependencies"
	@echo "  deps-check      Check for outdated dependencies"
	@echo "  deps-update     Update dependencies"
	@echo ""
	@echo "Utilities:"
	@echo "  clean           Remove build artifacts"
	@echo "  version         Show version"
	@echo "  help            Show this help"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION         Override version (default: git describe)"
	@echo "  ARGS            Arguments for 'run' and 'docker-run' targets"
