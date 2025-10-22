# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A production-ready Go project template for bootstrapping services and CLIs. Uses a clean `cmd/`-based layout with version embedding via `-ldflags`.

## Build and Development Commands

### Core Commands

```bash
make build          # Build all binaries to ./build/
make run            # Build and run the main binary
make test           # Run tests with coverage (outputs coverage.out)
make test-verbose   # Run tests with verbose output
make fmt            # Format all Go code
make clean          # Remove build artifacts, caches, and run git gc
```

### Cross-Platform Builds

```bash
make build-all      # Cross-compile for: linux/amd64, linux/arm64, windows/amd64, windows/arm64, darwin/amd64, darwin/arm64
```

Outputs: `./build/{cmd}-{os}-{arch}[.exe]`

### Code Quality

```bash
make lint-deadcode  # Detect unused code using staticcheck (U1000) and deadcode
golangci-lint run   # Linting (configured in CI, not in Makefile)
```

### Docker

```bash
docker build -t your/image:dev -f docker/Dockerfile .
```

Multi-stage build optimized for caching and minimal runtime.

## Architecture

### Project Layout

```
cmd/               # Binary entrypoints (one subdirectory per binary)
  └── go_template/ # Main CLI (rename this when creating new projects)
core/              # Shared packages and utilities
  └── version/     # Version info with semantic version parsing
build/             # Build outputs (gitignored)
docker/            # Docker build files
```

### Version System

Version information is injected at build time via Makefile ldflags:

- `go_template/core/version.Version` - from git tags (strips 'v' prefix)
- `go_template/core/version.BuildTime` - ISO 8601 timestamp
- `go_template/core/version.GitCommit` - short commit hash

The `core/version` package provides:

- `Get()` - returns full `Info` struct with all version fields
- `ParseVersion()` - parses semantic versions (handles `v` prefix, pre-release, build metadata)
- `IsNewerVersion()` - compares two versions semantically

### Multi-Command Structure

The Makefile automatically discovers all subdirectories in `cmd/` and builds them as separate binaries. Add new commands by creating `cmd/newcommand/main.go`.

## CI/CD Workflows

All workflows in `.github/workflows/`:

- `test.yml` - Runs `make test`, uploads coverage HTML artifact
- `code-quality-check.yml` - Runs `golangci-lint`
- `build_release.yml` - Cross-compiles on git tags, uploads release assets, generates changelog with git-cliff
- `build_image.yml` - Builds/pushes Docker image with buildx cache
- `release_drafter.yml` - Maintains draft release using Conventional Commits
- `auto_labeler.yml` - Auto-labels PRs based on `.github/labeler.yml`
- `code_scan.yml` - Security: gitleaks, trivy, CodeQL
- `semantic-pull-request.yml` - Enforces Conventional Commit format in PR titles

## Conventions

- Use Conventional Commit format for PR titles (enforced by CI)
- Tests colocated as `*_test.go` files
- Dependencies managed via Go modules (`go mod tidy`)
- Pre-commit hooks configured in `.pre-commit-config.yaml` (shellcheck, mdformat, codespell, gitleaks, etc.)

## Testing

Run individual tests:

```bash
go test -v -run TestFunctionName ./core/version
```

The main test suite covers semantic version parsing and comparison logic in `core/version/version_test.go`.
