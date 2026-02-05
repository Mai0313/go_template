<div align="center" markdown="1">

# Go 專案模板

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

🚀 幫助 Golang 開發者「快速建立新專案」的模板。提供務實的專案結構、Makefile、Docker 多階段建置，以及完整的 GitHub Actions 工作流程。

點擊 [使用此模板](../../generate) 後即可開始。

其他語言: [English](README.md) | [繁體中文](README.zh-TW.md) | [简体中文](README.zh-CN.md)

## ✨ 重點特色

- Makefile 工作流：build、test、跨平台編譯、fmt、dead‑code 掃描
- 內建版本資訊：以 `-ldflags` 注入 version、build time、git commit
- 範例 CLI：`cmd/go_template`，支援 `--version`
- 單元測試與 CI 覆蓋率報告產物
- Docker：多階段建置，最小化執行環境
- GitHub Actions：測試、靜態檢查（golangci‑lint）、映像建置/推送、Release Drafter、標籤、自動秘密/程式碼掃描

## 🚀 快速開始

需求：

- Go 1.24+
- Docker（可選，用於容器化建置）

本機開發：

```bash
make build            # 編譯到 ./build/
make run              # 編譯並執行主程式
make test             # 執行測試並輸出覆蓋率
make fmt              # go fmt ./...
make build-all        # 跨平台編譯常見 OS/ARCH
```

執行範例 CLI：

```bash
./build/go_template --version
```

## 作為模板使用

**重要提示**：這是一個模板，不是函式庫。你必須將 `go_template` 重新命名為你的專案名稱。

### 快速設定

1. 點擊 **使用此模板** 建立你的倉庫
2. 複製你的新倉庫
3. 執行重新命名腳本或按照下方手動步驟操作

### 手動重新命名步驟

**必要修改**（將 `{your_project}` 替換為你的實際專案名稱）：

1. **Go 模組**：
   - 更新 `go.mod`：`module go_template` → `module {your_project}`
   - 重新命名 `cmd/go_template/` → `cmd/{your_project}/`
   - 更新 `cmd/{your_project}/main.go` 中的匯入
   - 更新 `Makefile` 的 LDFLAGS（第17-19行）和 `BIN_NAME`（第23行）

2. **CLI 包裝器**（如果使用 npm/PyPI 發佈）：
   - Node.js：更新 `cli/nodejs/package.json` 和 `cli/nodejs/bin/start.js`
   - Python：更新 `cli/python/pyproject.toml` 並重新命名 `cli/python/src/go_template/`

3. **Docker**：
   - 更新 `docker/Dockerfile` 標籤和二進位路徑
   - 更新 `.devcontainer/Dockerfile` 標籤

4. **文件**：
   - 更新 `README.md`、`README.zh-CN.md`、`README.zh-TW.md` 中的徽章 URL
   - 更新 `.github/CODEOWNERS`

**驗證**：

```bash
make clean && make build
./build/{your_project} --version
grep -r "go_template" --exclude-dir=.git --exclude-dir=build .
```

詳細說明請參見 [.github/copilot-instructions.md](.github/copilot-instructions.md)。

## 專案結構

```text
cmd/go_template/     # 主 CLI 入口
core/version/        # 版本工具與測試
build/               # 編譯輸出（已加入 .gitignore）
docker/Dockerfile    # 多階段映像建置
```

## Docker

```bash
# 本地建置與執行
docker build -t your/image:dev -f docker/Dockerfile .
docker run --rm -it your/image:dev
```

## CI/CD（GitHub Actions）

- 測試：`.github/workflows/test.yml`
- 品質：`.github/workflows/code-quality-check.yml`
- 釋出打包：`.github/workflows/build_release.yml`
- Docker 映像：`.github/workflows/build_image.yml`
- 釋出草稿：`.github/workflows/release_drafter.yml`
- 標籤與語義化：`.github/workflows/auto_labeler.yml`, `semantic-pull-request.yml`
- 安全性：`.github/workflows/code_scan.yml`（gitleaks、trivy、codeql）

## 貢獻指南

- 提交前請執行 `make fmt && make test`
- PR 請聚焦單一變更並附上測試
- 使用 Conventional Commits 作為提交訊息
