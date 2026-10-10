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

本機開發：

```bash
make build            # 編譯到 ./build/
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
3. 依照 [AGENTS.md](AGENTS.md) 重新命名 `go_template` 及其他模板識別字，該檔案列出它們出現的位置以及驗證方式

## 開發

開發環境設定、專案結構、本機與 Docker 建置、測試、CI 與程式碼慣例請見 [CONTRIBUTING.md](./.github/CONTRIBUTING.md)。
