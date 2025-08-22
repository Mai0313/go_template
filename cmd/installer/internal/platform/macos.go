package platform

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"claude_analysis/cmd/installer/internal/logger"
)

// macOS-specific functionality

func InstallNodeDarwin() error {
	// Determine the correct Node.js archive based on architecture
	var nodeArchiveName string
	switch runtime.GOARCH {
	case "arm64":
		nodeArchiveName = "node-v22.18.0-darwin-arm64.tar.gz"
	case "amd64":
		nodeArchiveName = "node-v22.18.0-darwin-x64.tar.gz"
	default:
		return fmt.Errorf("unsupported architecture for macOS: %s", runtime.GOARCH)
	}
	targetDir, derr := getNodeInstallDir()
	if derr != nil {
		return derr
	}

	// Locate archive next to the installer executable
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("os.Executable failed: %w", err)
	}
	exeDir := filepath.Dir(exe)
	archivePath := filepath.Join(exeDir, nodeArchiveName)
	if _, err := os.Stat(archivePath); err != nil {
		return fmt.Errorf("required %s not found next to installer at %s: %w", nodeArchiveName, exeDir, err)
	}

	// Ensure target directory exists
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("create target dir %s: %w", targetDir, err)
	}

	logger.Info("📦 Extracting Node.js from tar.gz archive", fmt.Sprintf("From: %s\nTo: %s", archivePath, targetDir))
	if err := extractTarGZ(archivePath, targetDir); err != nil {
		return fmt.Errorf("extract node archive: %w", err)
	}

	// Some Node.js archives wrap files in a single version folder. Flatten it.
	if err := flattenIfSingleSubdir(targetDir); err != nil {
		logger.Warning("⚠️ Failed to flatten node directory", fmt.Sprintf("Error: %v", err))
	}

	// Update shell profile to persist environment variables
	if err := updateMacOSShellProfile(targetDir); err != nil {
		logger.Warning("⚠️ Failed to update shell profile", fmt.Sprintf("Error: %v", err))
	}

	// Also update current process environment so subsequent steps in this run can use node/npm immediately
	_ = os.Setenv("NODE_HOME", targetDir)
	_ = os.Setenv("PATH", addToPathUnix(os.Getenv("PATH"), targetDir))

	logger.Success("✅ Node.js installed on macOS.")
	return nil
}

// extractTarGZ extracts a tar.gz archive to the destination directory
func extractTarGZ(srcArchive, destDir string) error {
	cmd := exec.Command("tar", "-xzf", srcArchive, "-C", destDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract tar.gz archive: %w", err)
	}
	return nil
}

// updateMacOSShellProfile adds NODE_HOME and PATH to shell profile files
func updateMacOSShellProfile(nodeDir string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// List of shell profile files to update (macOS specific order)
	profiles := []string{
		filepath.Join(home, ".zshrc"),        // Default shell in macOS Catalina+
		filepath.Join(home, ".bash_profile"), // Bash on macOS
		filepath.Join(home, ".bashrc"),
		filepath.Join(home, ".profile"),
	}

	exportLines := []string{
		fmt.Sprintf("export NODE_HOME=\"%s\"", nodeDir),
		fmt.Sprintf("export PATH=\"%s:$PATH\"", nodeDir),
	}

	for _, profile := range profiles {
		if _, err := os.Stat(profile); err != nil {
			continue // Skip if file doesn't exist
		}

		// Check if our exports are already present
		content, err := os.ReadFile(profile)
		if err != nil {
			continue
		}

		needsUpdate := false
		for _, line := range exportLines {
			if !strings.Contains(string(content), line) {
				needsUpdate = true
				break
			}
		}

		if !needsUpdate {
			continue
		}

		// Append our exports
		f, err := os.OpenFile(profile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			continue
		}

		f.WriteString("\n# Added by Claude Code Installer\n")
		for _, line := range exportLines {
			if !strings.Contains(string(content), line) {
				f.WriteString(line + "\n")
			}
		}
		f.Close()
	}

	return nil
}
