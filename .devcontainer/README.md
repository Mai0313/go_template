# Dev Container for Go Project

This directory hosts configuration for a reproducible Go development environment using [VS Code Dev Containers](https://code.visualstudio.com/docs/devcontainers/containers).

## What's Included?

- **Dockerfile**: Go 1.x base image with zsh, oh-my-zsh, powerlevel10k, fonts, and common shell plugins.
- **devcontainer.json**: VS Code settings and extension recommendations (`golang.go`, Docker, YAML, TOML, etc.).
- Mounts for your `.gitconfig`, `.ssh`, and `.p10k.zsh`.

## Usage

1. Open this folder in VS Code with the Dev Containers extension installed.
2. Run “Dev Containers: Reopen in Container”.
3. On start, the container verifies `go version` and you're ready to `make build` / `make test`.

## Customization

- Add system packages in the Dockerfile as needed.
- Add VS Code extensions in `devcontainer.json`.
- Mount more files by editing the `mounts` array.

## Useful Commands

- Rebuild container after Dockerfile changes: “Dev Containers: Rebuild Container”.
- Inside the container: `make build`, `make test`, `go mod tidy`.

## Troubleshooting

- If SSH or Git behave unexpectedly, ensure your local files are mounted as configured.
- See the [VS Code Dev Containers docs](https://code.visualstudio.com/docs/devcontainers/containers) for more details.
