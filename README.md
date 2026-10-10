<div align="center" markdown="1">

# Go Project Template

[![Go](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go&logoColor=white)](https://go.dev/dl/)
[![npm version](https://img.shields.io/npm/v/@mai0313/go_template?logo=npm&style=flat-square&color=CB3837)](https://www.npmjs.com/package/@mai0313/go_template)
[![npm downloads](https://img.shields.io/npm/dt/@mai0313/go_template?logo=npm&style=flat-square)](https://www.npmjs.com/package/@mai0313/go_template)
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

Local setup:

```bash
make build            # build binaries into ./build/
```

Run the example CLI:

```bash
./build/go_template --version
```

## Using as a Template

**IMPORTANT**: This is a template, not a library. You must rename `go_template` to your project name.

### Quick Setup

1. Click **Use this template** to create your repository
2. Clone your new repository
3. Rename `go_template` and the other template identifiers by following [AGENTS.md](AGENTS.md), which lists where they appear and how to verify the rename

## Development

Contributor setup, project layout, local and Docker builds, tests, CI, and code conventions live in [CONTRIBUTING.md](./.github/CONTRIBUTING.md).
