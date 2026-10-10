<div align="center" markdown="1">

# Go 项目模板

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

🚀 面向 Golang 的生产级项目模板，帮助你快速创建新的 Go 服务或 CLI。内置合理的目录结构、Makefile、Docker 多阶段构建，以及完整的 CI/CD 工作流。

点击 [使用此模板](../../generate) 开始。

其他语言: [English](README.md) | [繁體中文](README.zh-TW.md) | [简体中文](README.zh-CN.md)

## ✨ 特性

- Makefile 任务：build、test、交叉编译、fmt、dead‑code 扫描
- 版本信息嵌入：通过 `-ldflags` 注入 version、build time、git commit
- 示例 CLI：`cmd/go_template` 支持 `--version`
- 单元测试，CI 上传覆盖率 HTML 产物
- Docker：多阶段构建，最小化运行时镜像
- GitHub Actions：测试、静态检查（golangci‑lint）、镜像构建/推送、Release Drafter、标签、机密/代码扫描

## 🚀 快速开始

前置条件：

- Go 1.24+

本地开发：

```bash
make build            # 编译到 ./build/
```

运行示例 CLI：

```bash
./build/go_template --version
```

## 作为模板使用

**重要提示**：这是一个模板，不是库。你必须将 `go_template` 重命名为你的项目名称。

### 快速设置

1. 点击 **使用此模板** 创建你的仓库
2. 克隆你的新仓库
3. 按照 [AGENTS.md](AGENTS.md) 重命名 `go_template` 及其他模板标识符，该文件列出了它们出现的位置以及验证方法

## 开发

开发环境设置、项目结构、本地与 Docker 构建、测试、CI 与代码规范请见 [CONTRIBUTING.md](./.github/CONTRIBUTING.md)。
