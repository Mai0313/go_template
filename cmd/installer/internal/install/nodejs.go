package install

import (
	"fmt"
	"runtime"

	"claude_analysis/cmd/installer/internal/logger"
	"claude_analysis/cmd/installer/internal/platform"
)

// InstallNodeJS checks for Node.js and installs if necessary
func InstallNodeJS() error {
	logger.Progress("📦 Step 1: Checking Node.js...")
	if !platform.CheckNodeVersion() {
		if platform.IsCommandAvailable("node") {
			logger.Warning("⚡ Node.js found but version is less than 22. Upgrading...")
		} else {
			logger.Info("📦 Node.js not found. Installing...")
		}

		switch runtime.GOOS {
		case "windows":
			if err := platform.InstallNodeWindows(); err != nil {
				return fmt.Errorf("failed to install Node.js on Windows: %w", err)
			}
		case "darwin":
			if err := platform.InstallNodeDarwin(); err != nil {
				return fmt.Errorf("failed to install Node.js on macOS: %w", err)
			}
		case "linux":
			if err := platform.InstallNodeLinux(); err != nil {
				return fmt.Errorf("failed to install Node.js on Linux: %w", err)
			}
		default:
			return fmt.Errorf("unsupported OS: %s", runtime.GOOS)
		}
	}

	logger.Success("✅ Node.js version >= 22 found. Skipping Node.js installation.")
	return nil
}
