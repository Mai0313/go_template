# Contributing Guide

Thank you for your interest in contributing to this Go project. This document describes how to set up the development environment, the conventions used by the project, and the workflow expected for issues and pull requests.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Ways to Contribute](#ways-to-contribute)
- [Reporting Issues](#reporting-issues)
- [Development Setup](#development-setup)
- [Local Workflow](#local-workflow)
- [Testing](#testing)
- [Branching Model](#branching-model)
- [Commit Convention](#commit-convention)
- [Pull Request Process](#pull-request-process)
- [Code Review](#code-review)
- [Coding Standards](#coding-standards)
- [Security Reports](#security-reports)
- [Licensing](#licensing)

## Code of Conduct

All contributors are expected to behave professionally and respectfully. Personal attacks, harassment, and discriminatory language are not tolerated. By participating, you agree to uphold a welcoming environment for everyone.

## Ways to Contribute

- Reporting bugs and reproducible issues
- Proposing or implementing new features
- Improving documentation, examples, and tutorials
- Reviewing pull requests and providing constructive feedback
- Suggesting tooling, performance, or security improvements

## Reporting Issues

Before opening a new issue:

1. Search existing issues to avoid duplicates.
2. Confirm the problem reproduces on the latest release or `main`.
3. Use the appropriate issue template.

Please include:

- A clear, descriptive title
- Go version (`go version`), OS, and architecture
- Project version, commit hash, or release tag
- Minimal reproduction steps and a code snippet when applicable
- Expected vs. actual behavior
- Full stack traces, logs, or screenshots

## Development Setup

```bash
# Verify your Go toolchain matches the version declared in go.mod
go version

# Clone your fork
git clone https://github.com/<your-username>/<repo>.git
cd <repo>

# Download dependencies
go mod download

# Build all commands declared under cmd/
make all
```

The required Go toolchain version is declared in `go.mod`. Install matching versions via [`gvm`](https://github.com/moovweb/gvm), [`asdf`](https://asdf-vm.com/), or your platform's package manager.

## Local Workflow

Common tasks are exposed via the `Makefile`. Run `make help` to list all targets. Frequently used ones:

```bash
make all       # Build all commands
make test      # Run the full test suite
make clean     # Remove build artifacts and caches
```

Recommended commands directly via `go`:

```bash
go fmt ./...                   # Format code
go vet ./...                   # Static analysis
go test ./...                  # Run all tests
go test -race -cover ./...     # Run with race detector and coverage
golangci-lint run              # Run the configured linters (if installed)
```

Always run `go fmt`, `go vet`, and the test suite before opening a pull request.

## Testing

- Tests live alongside source files using the `_test.go` suffix and Go's standard `testing` package.
- New behavior must be covered by tests. Bug fixes should include a regression test.
- Run with `-race` when changes touch concurrency-sensitive code.
- Use table-driven tests where they improve clarity.

Useful commands:

```bash
go test ./...                                  # Run all tests
go test ./path/to/pkg -run TestName            # Run a single test
go test -count=1 ./...                         # Disable test caching
go test -bench=. ./...                         # Run benchmarks
```

## Branching Model

- `main` is the default branch and must always be releasable.
- Feature branches: `feat/<short-description>`
- Bug fix branches: `fix/<short-description>`
- Documentation branches: `docs/<short-description>`

## Commit Convention

Commit messages follow [Conventional Commits](https://www.conventionalcommits.org/) and **must be written in English**.

Format:

```
<type>(<optional scope>): <short summary>

<optional body>

<optional footer>
```

Allowed types:

| Type | Purpose |
| --- | --- |
| `feat` | A new feature |
| `fix` | A bug fix |
| `refactor` | Code change that neither fixes a bug nor adds a feature |
| `doc` | Documentation-only changes |
| `perf` | Performance improvement |
| `style` | Formatting or stylistic changes |
| `test` | Adding or correcting tests |
| `chore` | Build, tooling, or auxiliary changes |
| `ci` | Continuous integration changes |
| `revert` | Reverting a previous commit |

Append `!` after the type or include `BREAKING CHANGE:` in the footer to indicate a breaking change. Reference issues with `Closes #123` or `Refs #123`.

## Pull Request Process

1. Ensure your branch is up to date with the target branch.
2. Run `go fmt`, `go vet`, and `make test` locally; all must pass.
3. Ensure CI checks pass on the pull request.
4. Use a descriptive title following the commit convention; it is validated by **semantic-pull-request**.
5. Fill out the pull request template, including motivation, summary, and testing notes.
6. Link related issues and design documents.
7. Mark the PR as **draft** while still in progress.
8. Request review only after self-review and a green CI.

Pull requests are typically merged via **squash merge** to keep history linear.

## Code Review

- Address all review comments or explain why a change is not needed.
- Keep discussions technical, focused, and respectful.
- Resolve conversations only after the concern has been addressed.

## Coding Standards

- **Formatting**: enforced by `gofmt` / `goimports`. Code must be formatted before committing.
- **Static analysis**: `go vet` must pass; `golangci-lint` is recommended.
- **Naming**: follow [Effective Go](https://go.dev/doc/effective_go) and standard idioms.
- **Errors**: wrap errors with context using `fmt.Errorf("...: %w", err)` and prefer `errors.Is` / `errors.As` for checks.
- **Concurrency**: prefer channels and `context.Context` for cancellation; avoid leaking goroutines.
- **Public APIs**: document all exported identifiers with complete sentences.

Prefer clarity over cleverness, and avoid unrelated refactors in feature or fix pull requests.

## Security Reports

Please **do not** report security vulnerabilities through public issues. Refer to [`SECURITY.md`](./SECURITY.md) for the responsible disclosure process.

## Licensing

By contributing, you agree that your contributions will be licensed under the project's license (see [`LICENSE`](../LICENSE)). Ensure that you have the right to submit any code, content, or assets you contribute.
