package install

import (
	"errors"
	"fmt"

	"claude_analysis/cmd/installer/internal/env"
	"claude_analysis/cmd/installer/internal/logger"
	"claude_analysis/cmd/installer/internal/platform"
)

// InstallOrUpdateClaude installs/updates Claude CLI
func InstallOrUpdateClaude() error {
	logger.Progress("🤖 Installing/Updating Claude Code CLI...")

	if err := installClaudeCLI(); err != nil {
		return fmt.Errorf("failed to install/update Claude CLI: %w", err)
	}

	logger.Success("✅ Claude Code CLI installation/update completed!")
	return InstallClaudeAnalysisBinary()
}

// installClaudeCLI installs the @anthropic-ai/claude-code package using npm.
// It first tries the default npm registry, and if that fails, it looks for a fallback registry from the available environments.
// It verifies the installation by checking if the `claude --version` command works.
func installClaudeCLI() error {
	baseArgs := []string{"install", "-g", "@anthropic-ai/claude-code@latest", "--no-color"}

	// --- 步驟 1: 嘗試使用預設 registry 安裝 ---
	logger.Info("📦 Attempting to install @anthropic-ai/claude-code with default registry...")
	err := platform.RunLoggedCmd(platform.GetNpmPath(), baseArgs...)

	// 如果第一次嘗試就成功，直接進行驗證並返回
	if err == nil {
		logger.Success("✅ Installation with default registry succeeded.")
		// 驗證安裝
		if verifyErr := verifyClaudeInstalled(); verifyErr != nil {
			return fmt.Errorf("installation verification failed: %w", verifyErr)
		}
		logger.Success("✅ Claude CLI installed successfully!")
		return nil
	}

	// --- 步驟 2: 如果第一次失敗，則尋找備用 registry 重試 ---
	logger.Warning("⚠️ Default registry failed, looking for a fallback...", fmt.Sprintf("Error: %v", err))

	chosen := env.SelectAvailableURL()
	if chosen.RegistryURL == "" {
		// 如果沒有找到備用 registry，返回第一次嘗試的錯誤
		return fmt.Errorf("npm install failed with default registry and no fallback registry is available: %w", err)
	}

	// 建立帶有 registry 的新參數
	retryArgs := append(baseArgs, "--registry="+chosen.RegistryURL)
	logger.Info("📦 Retrying installation with fallback registry", fmt.Sprintf("Registry: %s", chosen.RegistryURL))

	// 執行重試
	if retryErr := platform.RunLoggedCmd(platform.GetNpmPath(), retryArgs...); retryErr != nil {
		// 如果重試也失敗，返回重試時的錯誤
		return fmt.Errorf("npm install also failed on retry with registry %s: %w", chosen.RegistryURL, retryErr)
	}

	// --- 成功後的驗證 ---
	// 如果重試成功，進行驗證
	if verifyErr := verifyClaudeInstalled(); verifyErr != nil {
		return fmt.Errorf("installation verification failed after retry: %w", verifyErr)
	}

	logger.Success("✅ Claude CLI installed successfully!")
	return nil
}

// verifyClaudeInstalled checks if the claude CLI is installed by running `claude --version`.
func verifyClaudeInstalled() error {
	if path, ok := platform.FindClaudeBinary(); ok {
		return platform.RunLoggedCmd(path, "--version")
	}
	return errors.New("claude CLI not found after installation")
}

// RunFullInstall performs the complete installation process
func RunFullInstall() error {
	logger.Progress("🚀 Starting full Claude Code installation...")
	logger.SendProgress(0, 3, "Initializing installation...")

	// 1) Node.js check/install guidance
	logger.SendProgress(1, 3, "Checking and installing Node.js...")
	if err := InstallNodeJS(); err != nil {
		return err
	}

	// 2) Install @anthropic-ai/claude-code with registry fallbacks
	// and move claude_analysis to ~/.claude with platform-specific name
	logger.SendProgress(2, 3, "Installing Claude CLI and components...")
	if err := InstallOrUpdateClaude(); err != nil {
		return err
	}

	logger.SendProgress(3, 3, "Installation completed!")
	logger.Success("🎉 Installation completed successfully!")
	logger.Info("🔧 Automatically switching to GAISF API Key configuration...")
	return nil
}
