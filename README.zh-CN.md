# Go 模板项目

[English](README.md) | [繁體中文](README.zh-TW.md) | 简体中文

## 概述

这是一个具有良好结构化目录布局和完整 GitHub Actions 工作流程的 Go 项目模板。它为构建具有多个命令、跨平台构建和自动化 CI/CD 管道的 Go 应用程序提供了坚实的基础。

## 项目结构

```
.
├── cmd/                    # 命令行应用程序
│   ├── app/               # 主应用程序
│   └── cli/               # CLI 工具
├── core/                  # 核心业务逻辑
│   └── config/           # 配置管理
├── docker/               # Docker 配置
├── .github/              # GitHub Actions 工作流程
│   ├── workflows/        # CI/CD 管道
│   ├── ISSUE_TEMPLATE/   # 问题模板
│   └── ...              # 其他 GitHub 配置
├── go.mod               # Go 模块定义
├── Makefile            # 构建自动化
└── README.md           # 此文件
```

## 特性

- ✅ **多命令结构**：分离的 `app` 和 `cli` 命令
- ✅ **跨平台构建**：支持 Linux、macOS、Windows（AMD64 和 ARM64）
- ✅ **GitHub Actions**：完整的 CI/CD 管道，具有自动构建和发布功能
- ✅ **Docker 支持**：即用型 Docker 配置
- ✅ **配置管理**：灵活的配置系统，支持环境变量
- ✅ **Makefile 自动化**：简单的构建、测试和打包命令
- ✅ **多语言 README**：英文、繁体中文、简体中文

## 快速开始

### 先决条件

- Go 1.23.0 或更高版本
- Make（用于使用 Makefile 命令）
- Docker（可选，用于容器化）

### 安装

1. **克隆或使用此模板**：
   ```bash
   git clone <your-repo-url>
   cd go-template
   ```

2. **更新模块名称**：
   ```bash
   # 在 go.mod 中将 'go-template' 替换为您的实际模块名称
   go mod edit -module your-module-name
   ```

3. **安装依赖项**：
   ```bash
   go mod tidy
   ```

### 构建

#### 本地构建所有命令：
```bash
make all
```

#### 为特定平台构建：
```bash
make build_linux_amd64
make build_windows_amd64
make build_darwin_arm64
```

#### 为所有平台构建：
```bash
make build-all
```

#### 创建分发包：
```bash
make package-all
```

### 运行

#### 运行主应用程序：
```bash
./build/app
# 或
make run
```

#### 运行 CLI 工具：
```bash
./build/cli --help
./build/cli --version
```

## 开发

### 项目自定义

1. **更新应用程序名称**：
   - 修改 `Makefile` 中的 `BIN_NAME` 和 `CLI_NAME`
   - 更新 GitHub Actions 工作流程中的二进制文件名称

2. **添加您的业务逻辑**：
   - 在 `cmd/app/main.go` 中实现您的应用程序逻辑
   - 在 `cmd/cli/main.go` 中添加 CLI 命令和功能
   - 在 `core/` 下创建额外的包用于共享逻辑

3. **配置**：
   - 修改 `core/config/config.go` 以添加您的配置字段
   - 根据需要更新默认配置值

4. **Docker**：
   - 为您的应用程序需求自定义 `docker/Dockerfile`
   - 在 GitHub Actions 中更新 Docker 镜像名称

### 可用的 Make 命令

```bash
make help              # 显示可用命令
make all               # 本地构建所有命令
make build-all         # 为所有平台构建
make package-all       # 创建分发包
make clean             # 移除构建产物
make fmt               # 格式化 Go 代码
make test              # 运行测试
make run               # 构建并运行应用程序
```

## GitHub Actions 工作流程

此模板包含几个预配置的工作流程：

- **`build_package.yml`**：为所有平台构建和发布包
- **`build_image.yml`**：构建和推送 Docker 镜像
- **`auto_labeler.yml`**：自动标记拉取请求
- **`jira.yml`**：JIRA 集成用于问题跟踪
- **`updater.yml`**：自动依赖项更新

### 自定义工作流程

1. **更新仓库引用**：
   - 将 `gitea.mediatek.inc/IT-GAIA/go-template` 替换为您的仓库 URL
   - 如需要，更新 Docker 注册表 URL

2. **配置密钥**：
   - `GITHUB_TOKEN`：用于仓库访问
   - `GT_TOKEN`：用于 Docker 注册表访问
   - `JIRA_TOKEN`：用于 JIRA 集成（可选）
   - `SSH_KEY`：用于 SSH 访问（如需要）

3. **修改构建目标**：
   - 如果您不需要所有平台，请在工作流程中更新平台目标
   - 根据需要调整构建命令

## 配置

应用程序支持通过以下方式进行配置：

1. **配置文件**：`~/.go-template/config.json`
2. **环境变量**：`CONFIG_PATH` 指定自定义配置位置
3. **默认值**：用于开发的内建默认值

配置示例：
```json
{
  "version": "1.0.0",
  "environment": "production",
  "log_level": "info",
  "debug": false
}
```

## Docker 支持

构建 Docker 镜像：
```bash
docker build -f docker/Dockerfile -t your-app .
```

使用 Docker 运行：
```bash
docker run --rm your-app
```

## 贡献

1. Fork 仓库
2. 创建功能分支
3. 进行更改
4. 如适用，添加测试
5. 提交拉取请求

## 许可

此模板按原样提供供您使用。根据需要添加您自己的许可。

## 支持

- 为错误或功能请求创建问题
- 检查现有文档和示例
- 查看 GitHub Actions 日志以解决构建问题

---

## 下一步

设置此模板后：

1. **自定义应用程序逻辑**在 `cmd/` 目录中
2. **添加您的业务逻辑**在 `core/` 包中
3. **更新配置**以符合您的需求
4. **修改 GitHub Actions**以满足您的 CI/CD 要求
5. **在适当的目录中添加测试**
6. **更新文档**以反映您的应用程序

祝您编程愉快！🚀