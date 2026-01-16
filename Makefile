# PAC (Praqmatic Automated Changelog) Makefile
# Supports building for multiple platforms

BINARY_NAME := pac
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_DIR := ./bin
GO := go
GOFLAGS := -trimpath
LDFLAGS := -ldflags "-s -w -X main.Version=$(VERSION)"

# Platforms for cross-compilation
PLATFORMS := linux/amd64 linux/arm64 darwin/amd64 darwin/arm64 windows/amd64

# Default target
.PHONY: all
all: build

# Include modular Makefiles for packaging and Docker
-include Makefile.deb
-include Makefile.chocolatey
-include Makefile.winget
-include Makefile.docker

# Build for current platform
.PHONY: build
build:
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/pac

# Build for all platforms
.PHONY: build-all
build-all: $(PLATFORMS)

# Cross-compilation targets
.PHONY: $(PLATFORMS)
$(PLATFORMS):
	$(eval GOOS := $(word 1,$(subst /, ,$@)))
	$(eval GOARCH := $(word 2,$(subst /, ,$@)))
	$(eval EXT := $(if $(filter windows,$(GOOS)),.exe,))
	@echo "Building for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=0 $(GO) build $(GOFLAGS) $(LDFLAGS) \
		-o $(BUILD_DIR)/$(GOOS)-$(GOARCH)/$(BINARY_NAME)$(EXT) ./cmd/pac

# Create release archives
.PHONY: release
release: build-all
	@echo "Creating release archives..."
	@mkdir -p $(BUILD_DIR)/release
	@for platform in $(PLATFORMS); do \
		GOOS=$$(echo $$platform | cut -d'/' -f1); \
		GOARCH=$$(echo $$platform | cut -d'/' -f2); \
		if [ "$$GOOS" = "windows" ]; then \
			EXT=".exe"; \
			cd $(BUILD_DIR) && zip -q release/$(BINARY_NAME)-$(VERSION)-$$GOOS-$$GOARCH.zip $(BINARY_NAME)-$$GOOS-$$GOARCH$$EXT && cd ..; \
		else \
			cd $(BUILD_DIR) && tar -czf release/$(BINARY_NAME)-$(VERSION)-$$GOOS-$$GOARCH.tar.gz $(BINARY_NAME)-$$GOOS-$$GOARCH && cd ..; \
		fi; \
	done
	@echo "Release archives created in $(BUILD_DIR)/release/"

# Run tests
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
	$(GO) test -v ./internal/integration/... ./test/...

# Install to GOPATH/bin
.PHONY: install
install:
	$(GO) install $(GOFLAGS) $(LDFLAGS) ./cmd/pac

# Clean build artifacts
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

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

# Check for outdated dependencies
.PHONY: deps-check
deps-check:
	$(GO) list -u -m all

# Update dependencies
.PHONY: deps-update
deps-update:
	$(GO) get -u ./...
	$(GO) mod tidy

# Generate (if any code generation is needed)
.PHONY: generate
generate:
	$(GO) generate ./...

# Show version
.PHONY: version
version:
	@echo $(VERSION)

# Run the application
.PHONY: run
run: build
	$(BUILD_DIR)/$(BINARY_NAME) $(ARGS)

# Show help
.PHONY: help
help:
	@echo "PAC (Praqmatic Automated Changelog) Build System"
	@echo ""
	@echo "Build Targets:"
	@echo "  build          Build for current platform"
	@echo "  build-all      Build for all platforms (linux, darwin, windows)"
	@echo "  release        Create release archives for all platforms"
	@echo "  install        Install to GOPATH/bin"
	@echo ""
	@echo "Testing Targets:"
	@echo "  test           Run all tests"
	@echo "  test-coverage  Run tests with coverage report"
	@echo "  test-race      Run tests with race detector"
	@echo "  test-integration Run integration tests"
	@echo ""
	@echo "Package Targets (see Makefile.packages):"
	@echo "  deb            Build Debian package (.deb)"
	@echo "  winget         Create Windows Package Manager manifests"
	@echo "  chocolatey     Create Chocolatey package"
	@echo ""
	@echo "Docker Targets (see Makefile.docker):"
	@echo "  docker         Build Docker image"
	@echo "  docker-multiarch Build multi-arch Docker image"
	@echo ""
	@echo "Maintenance Targets:"
	@echo "  clean          Remove build artifacts"
	@echo "  lint           Run golangci-lint"
	@echo "  fmt            Format code"
	@echo "  deps-check     Check for outdated dependencies"
	@echo "  deps-update    Update dependencies"
	@echo "  version        Show version"
	@echo "  help           Show this help"
	@echo ""
	@echo "Variables:"
	@echo "  VERSION        Override version (default: git describe)"
	@echo "  ARGS           Arguments for 'run' target"
