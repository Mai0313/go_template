# Go 模板專案

[English](README.md) | 繁體中文 | [简体中文](README.zh-CN.md)

## 概述

這是一個具有良好結構化目錄佈局和完整 GitHub Actions 工作流程的 Go 專案模板。它為構建具有多個命令、跨平台構建和自動化 CI/CD 管道的 Go 應用程式提供了堅實的基礎。

## 專案結構

```
.
├── cmd/                    # 命令行應用程式
│   ├── app/               # 主應用程式
│   └── cli/               # CLI 工具
├── core/                  # 核心業務邏輯
│   └── config/           # 配置管理
├── docker/               # Docker 配置
├── .github/              # GitHub Actions 工作流程
│   ├── workflows/        # CI/CD 管道
│   ├── ISSUE_TEMPLATE/   # 問題模板
│   ├── CODEOWNERS        # 程式碼擁有者定義
│   ├── dependabot.yml    # 依賴項自動更新配置
│   ├── labeler.yml       # 自動標籤配置
│   └── ...              # 其他 GitHub 配置
├── .dockerignore         # Docker 忽略檔案
├── .gitignore           # Git 忽略檔案
├── go.mod               # Go 模組定義
├── Makefile            # 構建自動化
└── README.md           # 此檔案
```

## 特性

- ✅ **多命令結構**：分離的 `app` 和 `cli` 命令
- ✅ **跨平台構建**：支援 Linux、macOS、Windows（AMD64 和 ARM64）
- ✅ **GitHub Actions**：完整的 CI/CD 管道，具有自動構建和發布功能
- ✅ **Docker 支援**：即用型 Docker 配置
- ✅ **配置管理**：靈活的配置系統，支援環境變數
- ✅ **Makefile 自動化**：簡單的構建、測試和打包命令
- ✅ **多語言 README**：英文、繁體中文、簡體中文

## 快速開始

### 先決條件

- Go 1.23.0 或更高版本（當前支援到 Go 1.24.3）
- Make（用於使用 Makefile 命令）
- Docker（可選，用於容器化）

### 安裝

1. **克隆或使用此模板**：
   ```bash
   git clone <your-repo-url>
   cd go-template
   ```

2. **更新模組名稱**：
   ```bash
   # 在 go.mod 中將 'go-template' 替換為您的實際模組名稱
   go mod edit -module your-module-name
   ```

3. **安裝依賴項**：
   ```bash
   go mod tidy
   ```

### 構建

#### 本地構建所有命令：
```bash
make all
```

#### 為特定平台構建：
```bash
make build_linux_amd64
make build_windows_amd64
make build_darwin_arm64
```

#### 為所有平台構建：
```bash
make build-all
```

### 運行

#### 運行主應用程式：
```bash
./build/app
# 或
make run
```

#### 運行 CLI 工具：
```bash
./build/cli --help
./build/cli --version
```

## 開發

### 專案自定義

1. **更新應用程式名稱**：
   - 修改 `Makefile` 中的 `BIN_NAME` 和 `CLI_NAME`
   - 更新 GitHub Actions 工作流程中的二進制檔案名稱

2. **添加您的業務邏輯**：
   - 在 `cmd/app/main.go` 中實現您的應用程式邏輯
   - 在 `cmd/cli/main.go` 中添加 CLI 命令和功能
   - 在 `core/` 下創建額外的包用於共享邏輯

3. **配置**：
   - 修改 `core/config/config.go` 以添加您的配置欄位
   - 根據需要更新預設配置值

4. **Docker**：
   - 為您的應用程式需求自定義 `docker/Dockerfile`
   - 在 GitHub Actions 中更新 Docker 映像名稱

### 可用的 Make 命令

```bash
make help              # 顯示可用命令
make all               # 本地構建所有命令
make build-all         # 為所有平台構建
make clean             # 移除構建產物
make fmt               # 格式化 Go 程式碼
make test              # 運行測試
make test-verbose      # 運行詳細測試
make run               # 構建並運行應用程式
make install           # 安裝二進制檔案到系統（/usr/local/bin）
```

## GitHub Actions 工作流程

此模板包含幾個預配置的工作流程：

- **`build_package.yml`**：為所有平台構建和發布包
- **`build_image.yml`**：構建和推送 Docker 映像
- **`auto_labeler.yml`**：自動標記拉取請求
- **`jira.yml`**：JIRA 整合用於問題追蹤
- **`updater.yml`**：自動依賴項更新
- **`secret_scan.yml`**：安全掃描工作流程

### 自定義工作流程

1. **更新儲存庫引用**：
   - 將範例儲存庫 URL 替換為您的實際儲存庫 URL
   - 如需要，更新 Docker 註冊表 URL

2. **配置密鑰**：
   - `GITHUB_TOKEN`：用於儲存庫存取
   - `GT_TOKEN`：用於 Docker 註冊表存取
   - `JIRA_TOKEN`：用於 JIRA 整合（可選）
   - `SSH_KEY`：用於 SSH 存取（如需要）

3. **修改構建目標**：
   - 如果您不需要所有平台，請在工作流程中更新平台目標
   - 根據需要調整構建命令

## 配置

應用程式支援通過以下方式進行配置：

1. **配置檔案**：`~/.go-template/config.json`
2. **環境變數**：`CONFIG_PATH` 指定自定義配置位置
3. **預設值**：用於開發的內建預設值

配置範例：
```json
{
  "version": "1.0.0",
  "environment": "production",
  "log_level": "info",
  "debug": false
}
```

## Docker 支援

構建 Docker 映像：
```bash
docker build -f docker/Dockerfile -t your-app .
```

使用 Docker 運行：
```bash
docker run --rm your-app
```

## 貢獻

1. Fork 儲存庫
2. 創建功能分支
3. 進行更改
4. 如適用，添加測試
5. 提交拉取請求

## 授權

此模板按原樣提供供您使用。根據需要添加您自己的授權。

## 支援

- 為錯誤或功能請求創建問題
- 檢查現有文檔和範例
- 查看 GitHub Actions 日誌以解決構建問題

---

## 下一步

設置此模板後：

1. **自定義應用程式邏輯**在 `cmd/` 目錄中
2. **添加您的業務邏輯**在 `core/` 包中
3. **更新配置**以符合您的需求
4. **修改 GitHub Actions**以滿足您的 CI/CD 要求
5. **在適當的目錄中添加測試**
6. **更新文檔**以反映您的應用程式

祝您編程愉快！🚀