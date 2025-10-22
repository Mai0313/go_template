<div align="center" markdown="1">

# Go Project Template

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&logoColor=white)](https://go.dev/dl/)
[![tests](https://github.com/Mai0313/go_template/actions/workflows/test.yml/badge.svg)](.github/workflows/test.yml)
[![code-quality](https://github.com/Mai0313/go_template/actions/workflows/code-quality-check.yml/badge.svg)](https://github.com/Mai0313/go_template/actions/workflows/code-quality-check.yml)
[![pre-commit](https://img.shields.io/badge/pre--commit-enabled-brightgreen?logo=pre-commit)](https://github.com/pre-commit/pre-commit)
[![license](https://img.shields.io/badge/License-MIT-green.svg?labelColor=gray)](https://github.com/Mai0313/go_template/tree/master?tab=License-1-ov-file)
[![PRs](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/Mai0313/go_template/pulls)
[![contributors](https://img.shields.io/github/contributors/Mai0313/go_template.svg)](https://github.com/Mai0313/go_template/graphs/contributors)

</div>

🚀 A production‑ready Golang project template to bootstrap new Go services and CLIs quickly. It ships with a pragmatic layout, Makefile, Docker builds, and a complete CI/CD suite.

Click [Use this template](../../generate) to start a new repository from this scaffold.

Other Languages: [English](README.md) | [繁體中文](README.zh-TW.md) | [简体中文](README.zh-CN.md)

## ✨ Highlights

- Makefile tasks: build, test, cross‑compile, format, dead‑code scan
- Version embedding via `-ldflags` (version, build time, git commit)
- Example CLI under `cmd/go_template` with `--version`
- Unit tests with coverage artifact in CI
- Docker: multi‑stage image build with cache and minimal runtime
- GitHub Actions: test, lint (golangci‑lint), image build+push, release drafter, labels, secret/code scanning

## 🚀 Quick Start

Prerequisites:

- Go 1.24+
- Docker (optional, for container builds)

Local setup:

```bash
make build            # build binaries into ./build/
make run              # build and run the main command
make test             # run unit tests with coverage
make fmt              # format code (go fmt ./...)
make build-all        # cross‑compile common OS/ARCH targets
```

Run the example CLI:

```bash
./build/go_template --version
```

Use as a template:

1. Click Use this template to create your repository
2. Replace module name in `go.mod` as needed
3. Rename the command under `cmd/` if you want a different binary name

## Project Structure

```text
cmd/go_template/     # Main CLI entrypoint
core/version/        # Version utilities and tests
build/               # Build outputs (git‑ignored)
docker/Dockerfile    # Multi‑stage image build
```

## Docker

```bash
# Build & run image locally
docker build -t your/image:dev -f docker/Dockerfile .
docker run --rm -it your/image:dev
```

## CI/CD (GitHub Actions)

- tests: `.github/workflows/test.yml`
- quality: `.github/workflows/code-quality-check.yml`
- release package: `.github/workflows/build_release.yml`
- docker image: `.github/workflows/build_image.yml`
- release drafter: `.github/workflows/release_drafter.yml`
- labels & semantics: `.github/workflows/auto_labeler.yml`, `semantic-pull-request.yml`
- security: `.github/workflows/code_scan.yml` (gitleaks, trivy, codeql)

## Contribution

- Run `make fmt && make test` before pushing
- Keep PRs focused and small; include tests
- Use Conventional Commit messages
