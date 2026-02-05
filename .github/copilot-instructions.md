# Go Template Project - AI Coding Agent Instructions

## Using This Template (CRITICAL - READ FIRST)

**This is a project template, not a library.** When creating a new project from this template, you MUST rename all occurrences of `go_template` to your new project name. Follow this checklist:

### Step 1: Choose Your Project Name

- **Project name**: New module name (e.g., `myapp`, `awesome-cli`)
- **GitHub username**: Your GitHub username (currently: `Mai0313`)
- **Author info**: Your name and email (currently: `Wei Lee <mai@mai0313.com>`)
- **npm scope** (optional): Your npm org scope (currently: `@mai0313`)

### Step 2: Required File Modifications

Execute these changes systematically (use find/replace with whole-word matching):

#### Go Module Files (Critical - affects all imports)

1. **[go.mod](../go.mod)**: Replace `module go_template` → `module {your_project}`
2. **[cmd/go_template/](../cmd/go_template/)**: Rename directory → `cmd/{your_project}/`
3. **[cmd/{your_project}/main.go](../cmd/go_template/main.go)**: Update import `"go_template/core/version"` → `"{your_project}/core/version"`
4. **[Makefile](../Makefile)**:
    - Line 17-19: Update LDFLAGS `-X go_template/core/version.*` → `-X {your_project}/core/version.*`
    - Line 23: Update `BIN_NAME := go_template` → `BIN_NAME := {your_project}`

#### Node.js CLI Wrapper

1. **[cli/nodejs/package.json](../cli/nodejs/package.json)**:
    - `name`: `"go_template"` → `"{your_project}"` (or `"@yourscope/{your_project}"`)
    - `author`: Update to your name and email
    - `homepage`, `repository.url`, `bugs.url`: Replace `Mai0313/go_template` → `{username}/{your_project}`
    - `bin`: Update command names (both keys reference `go_template`)
2. **[cli/nodejs/bin/start.js](../cli/nodejs/bin/start.js)**:
    - Lines 14, 15, 23, 24, 29, 30: Update `binary: 'go_template'` → `binary: '{your_project}'`

#### Python CLI Wrapper

1. **[cli/python/pyproject.toml](../cli/python/pyproject.toml)**:
    - `name`: `"go_template"` → `"{your_project}"`
    - `authors`: Update to your info
    - `project.urls`: Replace `Mai0313/go_template` → `{username}/{your_project}`
    - `project.scripts`: Update command names
2. **[cli/python/src/go_template/](../cli/python/src/go_template/)**: Rename directory → `cli/python/src/{your_project}/`
3. **[cli/python/src/{your_project}/__init__.py](../cli/python/src/go_template/__init__.py)**:
    - Lines 27, 28, 33, 34, 39, 40: Update `binary: "go_template"` → `binary: "{your_project}"`

#### Docker & DevContainer

1. **[docker/Dockerfile](../docker/Dockerfile)**:
    - Line 2, 12: Update `maintainer` label to your info
    - Line 13: Update `org.label-schema.name="go_template"` → `org.label-schema.name="{your_project}"`
    - Line 14: Update `org.label-schema.vendor` to your name
    - Line 18: Update binary path `/usr/local/bin/go_template` → `/usr/local/bin/{your_project}`
    - Line 20: Update CMD to your binary name
2. **[.devcontainer/Dockerfile](../.devcontainer/Dockerfile)**:
    - Line 3: Update `maintainer` label to your info
    - Line 4: Update `org.label-schema.name="go_template"` → `org.label-schema.name="{your_project}"`
    - Line 5: Update `org.label-schema.vendor` to your name

#### Documentation & CI/CD

1. **[README.md](../README.md)**, **[README.zh-CN.md](../README.zh-CN.md)**, **[README.zh-TW.md](../README.zh-TW.md)**:
    - All badge URLs: Replace `Mai0313/go_template` → `{username}/{your_project}`
    - npm badges: Replace `@mai0313/go_template` → your npm package name
    - Update project description if needed
2. **[.github/CODEOWNERS](../.github/CODEOWNERS)**: Replace `@Mai0313` → `@{your_username}`

### Step 3: Verification

After renaming, verify your changes work:

```bash
# 1. Check module name in all Go files
grep -r "go_template" --include="*.go" .

# 2. Verify imports compile
go mod tidy && go build ./...

# 3. Test build
make clean && make build

# 4. Verify binary works
./build/{your_project} --version

# 5. Check no old references remain (excluding .git/)
grep -r "go_template" --exclude-dir=.git --exclude-dir=build .
grep -r "Mai0313" --exclude-dir=.git .
grep -r "mai@mai0313.com" --exclude-dir=.git .
```

### Common Mistakes to Avoid

- **CRITICAL**: Rename `cmd/go_template` directory AND update imports everywhere
- **CRITICAL**: Update Makefile LDFLAGS — version embedding will break otherwise
- **CRITICAL**: CLI wrapper binary names must match the Go binary name
- **CRITICAL**: Update all GitHub URLs — badges and links must point to your new repo

### Key Points for AI Agents

When helping users adopt this template:

1. **Always suggest a complete search-and-replace strategy** rather than file-by-file edits
2. **Verify imports after renaming** — `go mod tidy` and `go build ./...` must succeed
3. **Check binary name consistency** across Go (Makefile), Node.js (package.json + start.js), and Python (pyproject.toml + __init__.py)
4. **Update all three README files** (English, zh-CN, zh-TW) to maintain documentation consistency
5. **Test the build** with `make clean && make build` before committing changes

---

## Architecture Overview

This is a **multi-language wrapper project**: a Go core binary wrapped by Node.js (`cli/nodejs/`) and Python (`cli/python/`) CLI packages. The wrappers download platform-specific pre-built binaries and execute them transparently.

**Core pattern**: `cmd/go_template/main.go` → Go binary → wrapped by JS/Python → distributed via npm/PyPI

**Why CLI wrappers?**

- Enables distribution via npm/PyPI in addition to direct binary downloads
- Platform-specific binaries are downloaded on first install (not bundled)
- Users can run `npx {your-package}` or `pip install {your-package}` without Go toolchain
- Single release process generates artifacts for all platforms

## Module & Import Conventions

- **Module name**: `go_template` (declared in [go.mod](../go.mod))
- **Import pattern**: All internal imports use `go_template/` prefix (e.g., `go_template/core/version`)
- **Command discovery**: Makefile auto-discovers commands in `cmd/*` — each subdirectory becomes a separate binary
- **Binary naming**: Output binary name matches the `cmd/` subdirectory name

## Version Management

Version info is **embedded at build time** via Makefile ldflags:

```go
// Set in core/version/version.go via -ldflags at build time
Version = "v1.2.3"           // from git describe --tags
BuildTime = "2026-02-05T..." // ISO 8601 timestamp
GitCommit = "abc1234"        // short commit hash
```

The [core/version/version.go](../core/version/version.go) package also implements semantic version parsing that supports pre-release and build metadata (e.g., `1.2.3-alpha.1+build.123`).

**To display version**: `./build/go_template --version`

## Critical Build Commands

```bash
make build         # Build all commands in cmd/* (outputs to build/)
make run           # Build and run the main binary
make test          # Run tests with coverage (generates coverage.out)
make fmt           # Format all Go code
make package-all   # Cross-compile for 6 platforms + create versioned archives
make lint-deadcode # Detect unused code (staticcheck U1000 + deadcode)
make clean         # Remove build artifacts, caches, and run git gc
```

**Platform targets** (6 total): `linux/amd64`, `linux/arm64`, `windows/amd64`, `windows/arm64`, `darwin/amd64`, `darwin/arm64`

## Cross-Platform Build & Packaging

The Makefile generates platform-specific archives with this naming convention:

```
{command}-v{version}-{platform}.{ext}
```

Where:

- `{platform}` = `macos-x64|macos-arm64|linux-x64-gnu|linux-arm64-gnu|windows-x64|windows-arm64`
- `{ext}` = `tar.gz` (unix) or `zip` (windows)

Example: `go_template-v0.1.0-macos-arm64.tar.gz`

CLI wrappers (Node.js/Python) expect binaries in `binaries/{platform-dir}/go_template[.exe]`.

## Testing Conventions

- Tests live alongside source files (`*_test.go`)
- Coverage report: `coverage.out` (HTML: `coverage.html`)
- CI runs tests on every push/PR (ignoring `**/*.md` changes)
- CI skips tests for branches starting with `chore/`, `ci/`, or `docs/`

## Pre-commit Hooks

Pre-commit runs on **4 events**: `pre-commit`, `post-checkout`, `post-merge`, `post-rewrite`

Key hooks enforced:

- **shellcheck** for shell scripts
- **mdformat** with plugins (gofmt, footnotes, GFM, frontmatter)
- **codespell** spell checking
- **gitleaks** secret scanning
- Standard checks: JSON/YAML/TOML validation, trailing whitespace, EOF fixer

## Docker Workflow

Multi-stage [Dockerfile](../docker/Dockerfile):

1. **Builder stage** (`golang:1.24-alpine`): Compiles binary with `make`
2. **Production stage** (`alpine:3.21`): Minimal image with just the binary + CA certs

```bash
docker build -t your/image:dev -f docker/Dockerfile .
docker run --rm -it your/image:dev
```

## CI/CD Workflows

- [test.yml](workflows/test.yml): Runs `make test`, uploads coverage artifact
- [code-quality-check.yml](workflows/code-quality-check.yml): Linting with golangci-lint
- [build_release.yml](workflows/build_release.yml): Triggered on `v*` tags, runs `make package-all`, creates GitHub release
- [build_image.yml](workflows/build_image.yml): Builds and pushes Docker image
- [code_scan.yml](workflows/code_scan.yml): Security scanning (gitleaks, trivy, CodeQL)

**Release process**: Push a git tag starting with `v` (e.g., `v1.2.3`) to trigger automated builds and release creation.

## Project-Specific Patterns

1. **Auto-command discovery**: Makefile uses `$(notdir $(wildcard cmd/*))` to find all commands automatically
2. **Platform info lookup**: Both JS and Python wrappers use a `platform_map` dictionary to resolve `{os}-{arch}` → binary path
3. **Graceful version fallback**: If ldflags version is missing, [version.go](../core/version/version.go) attempts to read from `debug.ReadBuildInfo()` (for `go install` scenarios)
4. **Clean target is aggressive**: `make clean` removes build dir, coverage files, clears all Go caches, and runs `git gc --aggressive`

## Adding New Commands

To add a new command:

1. Create `cmd/{newcommand}/main.go`
2. Run `make build` — Makefile auto-discovers it
3. Binary appears at `build/{newcommand}`

Makefile changes are required only for advanced build configurations.

## Key Files Reference

- [Makefile](../Makefile): Build orchestration, cross-compilation, packaging
- [cmd/go_template/main.go](../cmd/go_template/main.go): CLI entrypoint with `--version` flag
- [core/version/version.go](../core/version/version.go): Version info struct, semantic version parser
- [cli/nodejs/bin/start.js](../cli/nodejs/bin/start.js): Node.js wrapper that spawns Go binary
- [cli/python/src/go_template/__init__.py](../cli/python/src/go_template/__init__.py): Python wrapper equivalent
- [.pre-commit-config.yaml](../.pre-commit-config.yaml): Pre-commit hook configuration
