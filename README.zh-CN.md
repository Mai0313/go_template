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
- Docker（可选，用于容器构建）

本地开发：

```bash
make build            # 编译到 ./build/
make run              # 编译并运行主命令
make test             # 运行单元测试并生成覆盖率
make fmt              # go fmt ./...
make build-all        # 交叉编译常见 OS/ARCH
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
3. 运行重命名脚本或按照下方手动步骤操作

### 手动重命名步骤

**必需修改**（将 `{your_project}` 替换为你的实际项目名称）：

1. **Go 模块**：
   - 更新 `go.mod`：`module go_template` → `module {your_project}`
   - 重命名 `cmd/go_template/` → `cmd/{your_project}/`
   - 更新 `cmd/{your_project}/main.go` 中的导入
   - 更新 `Makefile` 的 LDFLAGS（第17-19行）和 `BIN_NAME`（第23行）

2. **CLI 包装器**（如果使用 npm/PyPI 分发）：
   - Node.js：更新 `cli/nodejs/package.json` 和 `cli/nodejs/bin/start.js`
   - Python：更新 `cli/python/pyproject.toml` 并重命名 `cli/python/src/go_template/`

3. **Docker**：
   - 更新 `docker/Dockerfile` 标签和二进制路径
   - 更新 `.devcontainer/Dockerfile` 标签

4. **文档**：
   - 更新 `README.md`、`README.zh-CN.md`、`README.zh-TW.md` 中的徽章 URL
   - 更新 `.github/CODEOWNERS`

**验证**：

```bash
make clean && make build
./build/{your_project} --version
grep -r "go_template" --exclude-dir=.git --exclude-dir=build .
```

详细说明请参见 [.github/copilot-instructions.md](.github/copilot-instructions.md)。

## 项目结构

```text
cmd/go_template/     # 主 CLI 入口
core/version/        # 版本工具与测试
build/               # 编译输出（已被 .gitignore 忽略）
docker/Dockerfile    # 多阶段镜像构建
```

## Docker

```bash
# 本地构建与运行
docker build -t your/image:dev -f docker/Dockerfile .
docker run --rm -it your/image:dev
```

## CI/CD（GitHub Actions）

- 测试：`.github/workflows/test.yml`
- 质量：`.github/workflows/code-quality-check.yml`
- 发布打包：`.github/workflows/build_release.yml`
- Docker 镜像：`.github/workflows/build_image.yml`
- 发布草稿：`.github/workflows/release_drafter.yml`
- 标签与语义化：`.github/workflows/auto_labeler.yml`, `semantic-pull-request.yml`
- 安全：`.github/workflows/code_scan.yml`（gitleaks、trivy、codeql）

## 贡献指南

- 提交前执行 `make fmt && make test`
- PR 聚焦单一变更并附带测试
- 使用 Conventional Commits 提交信息
