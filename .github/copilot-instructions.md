<!-- Workspace-specific instructions for GitHub Copilot. Keep this in sync with the repo. -->

⚠️ IMPORTANT: Update this document whenever tooling, commands, or workflows change.

### Project Background

This repository is a Golang project template to bootstrap services and CLIs quickly. It ships with a clean Go layout (`cmd/`, `core/`), Makefile tasks, Docker multi-stage builds, and a comprehensive GitHub Actions suite.

### Core Infrastructure

- Go 1.24+
- `cmd/` for binaries, `core/` for shared packages
- Makefile for build/test/cross-compile/format
- Dockerfile (multi-stage) under `docker/`

### Local Development

- Build: `make build` (outputs to `./build/`)
- Run: `make run` or execute `./build/go-template`
- Test: `make test` (produces `coverage.out`)
- Format: `make fmt` (runs `go fmt ./...`)
- Dead code checks: `make lint-deadcode`

Make targets reference (from `Makefile`):

- `make clean` — remove build/test caches and artifacts
- `make build-all` — cross-compile common OS/ARCH targets

### CLI Entrypoint

Main binary lives at `cmd/go-template`. It supports `--version` which prints version, build time, git commit, and Go version injected via `-ldflags`.

### Coding Style

- Use standard Go formatting: `go fmt ./...`
- Optional linting: `golangci-lint run` (action provided in CI)
- Tests colocated with code as `*_test.go`

### Dependencies

- Managed via Go modules (`go.mod`, `go.sum`)
- Tidy modules with `go mod tidy`

### Docker

- Multi-stage builder in `docker/Dockerfile`
- Local build: `docker build -t your/image:dev -f docker/Dockerfile .`

### CI/CD Workflows (GitHub Actions)

All workflows live in `.github/workflows/`:

- `test.yml`: Run `make test` with coverage and upload HTML artifact
- `code-quality-check.yml`: Run `golangci-lint`
- `build_release.yml`: Cross-compile on tags and upload release assets
- `build_image.yml`: Build/push Docker image with buildx cache
- `release_drafter.yml`: Maintain a draft release using Conventional Commits
- `auto_labeler.yml`: Auto-apply labels based on `.github/labeler.yml`
- `code_scan.yml`: Security scans (gitleaks, trivy) and CodeQL
- `semantic-pull-request.yml`: Enforce Conventional Commit PR titles

### Conventions

- Use Conventional Commit PR titles (enforced by workflow)
- Keep PRs small and focused with tests
- Update this file plus `README.md`, `README.zh-TW.md`, and `README.zh-CN.md` when commands or workflows change
