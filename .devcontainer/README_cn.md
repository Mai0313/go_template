# Go 專案 Dev Container 中文說明

本目錄提供 VS Code Dev Container 的設定，讓你在一致且可重現的 Go 開發環境中工作。

## 內容說明

- **Dockerfile**：以 Go 1.x 為基礎，安裝 zsh、oh-my-zsh、powerlevel10k 與常用外掛/字型。
- **devcontainer.json**：VS Code 容器設定（建議擴充：`golang.go`、Docker、YAML、TOML 等）。
- 掛載本機 `.gitconfig`、`.ssh`、`.p10k.zsh`。

## 使用方式

1. 以 VS Code 開啟本資料夾並安裝 Dev Containers 擴充套件。
2. 執行「Dev Containers: Reopen in Container」。
3. 進入容器後可直接執行 `make build`、`make test` 等指令。

## 自訂化

- 若需額外系統套件，請編輯 Dockerfile。
- 若需更多 VS Code 擴充套件，請調整 `devcontainer.json`。
- 若需掛載更多檔案，請修改 `mounts` 陣列。

## 常用指令

- 變更 Dockerfile 後，使用「Dev Containers: Rebuild Container」。
- 容器中可執行：`make build`、`make test`、`go mod tidy`。

## 疑難排解

- 若 SSH 或 Git 有異常，請確認本機檔案已依設定掛載。
- 參考 [VS Code Dev Containers 文件](https://code.visualstudio.com/docs/devcontainers/containers)。
