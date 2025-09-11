<!-- Workspace-specific instructions for GitHub Copilot. Keep this in sync with the repo. -->

⚠️ IMPORTANT: Update this document whenever tooling, commands, or workflows change.

### Project Background

This repository is a Python project template intended to help developers bootstrap projects quickly. It ships with modern packaging, CLI entry points, docs generation, Docker/Compose, and a comprehensive CI/CD setup via GitHub Actions.

### Core Infrastructure

- Python 3.10/3.11/3.12 supported
- Dependency management: `uv`
- Project layout: `src/` packaging
- Dockerfile with multi-stage builds; `docker-compose.yaml` for local services
- MkDocs Material documentation with mkdocstrings

### Local Development

- Install deps: `make uv-install && uv sync`
- Run quality hooks: `make format`
- Run tests: `make test`
- Generate docs: `make gen-docs`
- Serve docs: `uv run mkdocs serve`

Dependency groups:

- `uv sync --group test` for test-only deps
- `uv sync --group docs` for docs-only deps

Make targets reference (from `Makefile`):

- `make clean`: remove caches, artifacts and generated docs
- `make format`: run pre-commit hooks (ruff etc.)
- `make test`: run pytest
- `make gen-docs`: generate docs for `src/` and `scripts/`
- `make submodule-init|submodule-update`: submodule helpers (optional)
- `make uv-install`: install `uv` on the system

### CLI Entrypoints

Defined in `pyproject.toml` under `[project.scripts]`:

- `repo_template` and `cli` → `repo_template.cli:main`

Example:

```bash
uv run repo_template
uv run cli
```

### Coding Style

- Use ruff with repo-configured rules (run via pre-commit or `make format`)
- PEP 8 naming: snake_case (functions/variables), PascalCase (classes), UPPER_CASE (constants)
- Prefer full type hints on public functions and datamodels
- Use Pydantic models with `Field(..., description=...)` where appropriate
- Tests in `tests/`, discoverable by `pytest`

### Type Hints and Docs

- Type-annotate function parameters and returns
- Prefer Google-style docstrings (configured in mkdocstrings)
- Keep code/doc comments concise; English for docs

### Dependencies (uv)

- Prod: `uv add <package>`, `uv remove <package>`
- Dev: `uv add <package> --dev`, `uv remove <package> --dev`

Build and publish:

```bash
uv build                 # create wheel/sdist in dist/
UV_PUBLISH_TOKEN=... uv publish   # publish to PyPI
```

### Documentation Generation

- Script: `scripts/gen_docs.py` supports `.py` and `.ipynb`
- Examples:

```bash
uv run python ./scripts/gen_docs.py --source ./src --output ./docs/Reference gen_docs
uv run python ./scripts/gen_docs.py --source ./src --output ./docs/Reference --mode file gen_docs
```

### Docker/Compose

- `docker-compose.yaml` provides optional `redis`, `postgresql`, `mongodb`, `mysql` services and an example `app` service
- Configure via `.env` (see `README.md` for keys)

### CI/CD Workflows (GitHub Actions)

All workflows live in `.github/workflows/`:

- `test.yml`: Run pytest on PRs to `master`/`release/*` (3.10/3.11/3.12)
- `code-quality-check.yml`: Run pre-commit hooks on PRs
- `deploy.yml`: Build and publish MkDocs site on pushes to `master` and tags `v*`
- `build_package.yml`: Build wheel/sdist on tags `v*`, upload artifacts, generate changelog; optional PyPI publish with `UV_PUBLISH_TOKEN`
- `build_image.yml`: Build and push Docker image to GHCR on `master` and tags `v*`
- `build_executable.yml`: Example Windows packaging flow on tags `v*` (stub)
- `release_drafter.yml`: Maintain a draft release using Conventional Commits
- `auto_labeler.yml`: Auto-apply labels based on `.github/labeler.yml`
- `secret_scan.yml`: Run gitleaks on push/PR
- `semantic-pull-request.yml`: Enforce Conventional Commit PR titles

### Conventions

- Conventional Commit PR titles (enforced by workflow)
- Prefer small, focused PRs with tests and docs
- Update this file plus `README.md`, `README.zh-TW.md`, and `README.zh-CN.md` when workflows, commands, or usage changes

### Running from PyPI via uvx

After publishing your package, the CLI can be executed without prior install using `uvx`:

```bash
uvx repo_template
uvx --from your-package-name==0.1.0 your-entrypoint
```
