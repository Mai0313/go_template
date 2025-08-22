package install

import (
	"fmt"
	"os"
	"path/filepath"

	"claude_analysis/cmd/installer/internal/logger"
	"claude_analysis/cmd/installer/internal/platform"
)

// InstallClaudeAnalysisBinary installs the claude_analysis binary to ~/.claude
func InstallClaudeAnalysisBinary() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("unable to get home dir: %w", err)
	}
	targetDir := filepath.Join(homeDir, ".claude")
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("failed to create %s: %w", targetDir, err)
	}

	// Determine source binary path: same directory as this installer
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("os.Executable failed: %w", err)
	}
	srcDir := filepath.Dir(exe)
	srcName := platform.ExeName("claude_analysis")
	srcPath := filepath.Join(srcDir, srcName)
	if _, err := os.Stat(srcPath); err != nil {
		return fmt.Errorf("expected %s next to installer: %w", srcName, err)
	}

	// Destination filename uses the same simple naming convention
	destName := platform.ExeName("claude_analysis")
	destPath := filepath.Join(targetDir, destName)

	// Copy the binary to destination and keep the original
	if err := platform.CopyFile(srcPath, destPath, 0o755); err != nil {
		return fmt.Errorf("failed to install claude_analysis to %s: %w", destPath, err)
	}
	logger.Success("✅ Claude Analysis binary installed successfully", fmt.Sprintf("Location: %s", destPath))
	return nil
}
