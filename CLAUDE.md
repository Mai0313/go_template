# Initializing a New Project from This Go Template

This repository is a project template, not a shared library to extend in place. Before implementing product features, initialize the project identity so the Go binary, npm and PyPI packages, Docker image, GitHub Actions, and documentation all describe the same project. Do not rename only the visible product name, and do not hide unresolved references to make a search appear clean.

## Establish the Project Identity First

Confirm the following with the user or derive it from the new repository. If required information is unavailable, ask rather than guessing a publishing name, owner, or account.

| Information                                                     | Used by                                                                                   |
| --------------------------------------------------------------- | ----------------------------------------------------------------------------------------- |
| GitHub owner and repository name                                | Module path, README badges, URLs, CODEOWNERS, Actions publishing, and container image     |
| Go module path                                                  | `go.mod` and every internal Go import. This is usually `github.com/<owner>/<repository>`. |
| Primary binary name                                             | `cmd/<binary>/`, `Makefile` `BIN_NAME`, Docker, and wrapper executable names              |
| Display name, description, license, author, and contact details | README, npm metadata, PyPI metadata, and Docker labels                                    |
| Whether npm, PyPI, and Docker image distribution are required   | The basis for retaining or adapting wrappers and publishing workflows                     |
| npm scope and PyPI package name                                 | Publishing workflow, `package.json`, `pyproject.toml`, and package entry points           |
| Supported operating systems and CPU architectures               | `Makefile` cross-compilation matrix and both wrapper platform maps                        |

Choose a lowercase, CLI-safe primary binary name. The Go package name, npm package name, PyPI package name, and GitHub repository name do not have to be identical, so record their mapping explicitly. The Python import package must be a valid Python identifier. If the public module path includes `/vN`, ensure the Makefile `-X <module>/core/version.*` targets exactly match the Go import path.

## Inventory Before Renaming

The template placeholders are not limited to `go_template`. Before initialization, build an inventory and distinguish intentional historical examples from values that must become new-project metadata.

```bash
rg -n --hidden -g '!/.git/**' -g '!build/**' \
    -e 'go_template|gotemp|Mai0313|mai@mai0313\\.com|@mai0313|post_hook' .
```

Review at least the following locations. Do not depend on stale line numbers.

| Area                                   | Required review or update                                                                                                                                                                                                                                                                                                         |
| -------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Go core                                | Update the module path in `go.mod`. Rename `cmd/go_template/` to the primary binary directory. Update the internal import in `main.go`. Keep `core/version` as the single build-time version source unless the new project has an explicit replacement design.                                                                    |
| Makefile                               | Update all three `LDFLAGS` `-X` targets to the new module path. Keep `BIN_NAME` aligned with the primary `cmd/` directory and Docker binary. `CMDS` automatically compiles every `cmd/*` directory. Correct the stale `post_hook` comment placeholder as well.                                                                    |
| Node.js wrapper                        | Update `cli/nodejs/package.json` name, author, description, URLs, keywords, `bin`, and package contents. Every binary name in `bin/start.js` must match its Go release artifact. Remove the `gotemp` example alias unless the new project deliberately supports it.                                                               |
| Python wrapper                         | Update project metadata, URLs, scripts, and classifiers in `cli/python/pyproject.toml`. Rename `cli/python/src/go_template/` to the new Python import package, then update its entry point and every binary name in `__init__.py`.                                                                                                |
| Containers and development environment | Update maintainer labels, image name labels, binary `COPY`, `chmod`, and `CMD` in `docker/Dockerfile`, plus labels in `.devcontainer/Dockerfile`. Treat the devcontainer as configurable development tooling: review its disabled SSL verification and personal shell customization instead of retaining them without a decision. |
| GitHub metadata                        | Review `.github/CODEOWNERS`, issue templates, label configuration, Dependabot, and release settings for the new team. Check every owner, repository URL, package scope, and secret reference.                                                                                                                                     |
| Documentation                          | Update project names, purpose, installation and execution examples, project structure, badge URLs, and related links in `README.md`, `README.zh-TW.md`, and `README.zh-CN.md`. All three READMEs must describe the real product.                                                                                                  |

Use `git mv` for directory renames so Git can preserve history, then update content. Do not restrict the search to `*.go`: module paths also occur in the Makefile, Dockerfiles, wrappers, READMEs, Actions, and release metadata.

## Preserve README Badges and Turn READMEs into Product Documentation

Every existing badge in every README is permanent. Never delete an existing README badge. During initialization, only add badges or modify an existing badge's URL, package name, workflow link, or display data to match the new project. For example, if npm publishing is not yet available, update its badge to a correct and verifiable target or document the release status. This rule does not have a later-project exception.

After initialization, READMEs must be product documentation rather than copies of the template setup guide. Describe what the project does, prerequisites, real commands, common workflows, container usage, and release method. `CLAUDE.md` is the Agent instruction document for initialization and maintenance.

## GitHub Actions and Release Contracts

Treat the existing `.github/workflows/` files as a starting point. During initialization, modify workflows rather than deleting them. Each currently represents a potentially needed protection or release capability: tests, code quality, CodeQL and secret scanning, image publishing, release drafts, release artifacts, Dependabot, PR labels, and semantic PR titles.

GitHub Actions may be reduced only later in the project lifecycle, after there is evidence that a capability is no longer needed and the user explicitly confirms its removal. Remove a workflow only with its associated configuration, documentation, and release assumptions, rather than deleting it merely to make Actions pass.

`build_release.yml` is particularly prone to partial migrations. On a `v*` tag, it runs `make package-all` for six targets, uploads archives to a GitHub Release, and extracts them into the `binaries/<platform>/` locations expected by the Node.js and Python wrappers. When changing the binary name, module, package names, npm scope, Python import package, supported platforms, or release options, verify all of the following:

1. The archive names generated by the Makefile are still parsed correctly by both shell extraction steps.
2. The repository name used by Actions equals the artifact binary prefix. If it does not, introduce an explicit binary variable instead of relying on coincidence.
3. The npm matrix scope, package name, and trusted publisher belong to the new project. Remove the template `gotemp` publishing entry or replace it with an intentional supported alias.
4. The PyPI `src/<package>/binaries` path matches the renamed Python import package.
5. The Docker image destination, GitHub package permissions, PyPI token or trusted publishing, and related secrets are configured for the new repository.

If the new project does not yet publish npm, PyPI, or Docker artifacts, make the related workflow explicitly disabled or validation-only and document the current status in the README. Do not delete it during initialization.

When creating or editing workflows, do not add `container` fields or `Setup MTK Certification` steps. Order job attributes as `name`, `needs`, `runs-on`, `if`, then other attributes. Order step attributes as `name`, `id`, `continue-on-error`, `if`, `uses`, `with`, `env`, `shell`, `run`. Write a GitHub expression directly in a command when it is used once rather than adding a redundant environment variable.

## Development and Verification Order

After renaming and metadata changes, make the shortest Go path work before testing the retained distribution surfaces.

```bash
go mod tidy
make fmt
make test
make build
./build/<binary-name> --version
```

`make build` compiles every `cmd/*` directory, so each additional command must have a valid entry point. Before a release, also run `make clean && make build`, but inspect the Makefile first: the current `clean` target removes build artifacts, clears Go caches, runs `git fetch --prune`, and performs aggressive `git gc`. Do not run it casually where repository state must not change or network access is unavailable.

Add verification appropriate to every retained distribution surface.

| Surface           | Expected verification                                                                                                                    |
| ----------------- | ---------------------------------------------------------------------------------------------------------------------------------------- |
| Go release        | `make package-all` succeeds, archive names and contents are correct for all six platforms, and Windows binaries use `.exe`.              |
| Node.js           | Package metadata is correct and `node cli/nodejs/bin/start.js --version` forwards arguments when the current platform binary is present. |
| Python            | Build from `cli/python` with `uv`, confirm the package entry point forwards `--version`, and verify package data includes binaries.      |
| Docker            | `docker build -f docker/Dockerfile .` succeeds and the container's default command is the new binary.                                    |
| CI quality checks | `pre-commit run -a` and `golangci-lint run ./...` pass, or record exactly which unavailable tool prevented execution.                    |

Finally, search for template identifiers again, excluding `.git`, `build`, and intentional template descriptions in this instruction file. Do not delete files or add search exclusions to hide unfinished references. Review every remaining match and decide whether it is justified.

## Ongoing Development Rules

Add a Go command under `cmd/<name>/main.go` because the Makefile discovers it automatically. Only commands intended for external distribution should also be added to the Docker entry point, wrapper binary mappings, and release documentation. Preserve the `core/version` ldflags contract during version changes; otherwise `--version` and release artifact traceability break.

After code changes, run the applicable formatter and tests. Commit messages and GitHub PR titles must be in English and follow Conventional Commits, for example `feat(cli): add status command`. Do not claim that template initialization is complete until every retained dependent surface has been verified.
