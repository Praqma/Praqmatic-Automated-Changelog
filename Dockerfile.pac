# Multi-stage Dockerfile for Go-based PAC
# Produces a minimal image with just the compiled binary

# ============================================================================
# Build stage: Compile the Go binary
# ============================================================================
FROM golang:1.24-alpine AS builder

# Install git for VCS information during build
RUN apk add --no-cache git ca-certificates

WORKDIR /build

# Copy go.mod and go.sum first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build arguments for version information
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown

# Build the binary with version information embedded
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags "-s -w -X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME}" \
    -o pac \
    ./cmd/pac

# ============================================================================
# Runtime stage: Minimal image with just the binary
# ============================================================================
FROM alpine:3.20 AS runtime

# Install git (required for repository operations) and ca-certificates
RUN apk add --no-cache git ca-certificates tzdata

# Create non-root user for security
RUN adduser -D -u 1000 pac

# Copy the compiled binary from builder
COPY --from=builder /build/pac /usr/local/bin/pac

# Copy default templates
COPY --from=builder /build/templates /usr/share/pac/templates

# Set up working directory
WORKDIR /repo

# Switch to non-root user
USER pac

# Default entrypoint
ENTRYPOINT ["pac"]
CMD ["--help"]

# ============================================================================
# Labels for container metadata
# ============================================================================
LABEL org.opencontainers.image.title="PAC - Praqmatic Automated Changelog"
LABEL org.opencontainers.image.description="Git changelog generator that integrates with task management systems"
LABEL org.opencontainers.image.url="https://github.com/Praqma/Praqmatic-Automated-Changelog"
LABEL org.opencontainers.image.source="https://github.com/Praqma/Praqmatic-Automated-Changelog"
LABEL org.opencontainers.image.licenses="MIT"
