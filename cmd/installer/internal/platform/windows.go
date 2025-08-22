package platform

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"claude_analysis/cmd/installer/internal/logger"
)

// Windows-specific functionality

func InstallNodeWindows() error {
	const nodeZipName = "node-v22.18.0-win-x64.zip"
	// Install under user's home to avoid requiring Administrator
	targetDir, derr := getNodeInstallDir()
	if derr != nil {
		return derr
	}

	// Locate zip next to the installer executable
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("os.Executable failed: %w", err)
	}
	exeDir := filepath.Dir(exe)
	zipPath := filepath.Join(exeDir, nodeZipName)
	if _, err := os.Stat(zipPath); err != nil {
		return fmt.Errorf("required %s not found next to installer at %s: %w", nodeZipName, exeDir, err)
	}

	// Ensure target directory exists
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("create target dir %s: %w (try running as Administrator)", targetDir, err)
	}

	logger.Info("📦 Extracting Node.js from zip archive", fmt.Sprintf("From: %s\nTo: %s", zipPath, targetDir))
	if err := unzip(zipPath, targetDir); err != nil {
		return fmt.Errorf("extract node zip: %w", err)
	}

	// Some Node.js zips wrap files in a single version folder. Flatten it.
	if err := flattenIfSingleSubdir(targetDir); err != nil {
		logger.Warning("⚠️ Failed to flatten node directory", fmt.Sprintf("Error: %v", err))
	}

	// Persist user environment variables (User scope)
	if err := setWindowsUserEnv("NODE_HOME", targetDir); err != nil {
		logger.Warning("⚠️ Failed to set NODE_HOME (user)", fmt.Sprintf("Error: %v", err))
	}
	// Ensure PATH includes targetDir
	if err := ensureWindowsUserPathIncludes(targetDir); err != nil {
		logger.Warning("⚠️ Failed to update PATH (user)", fmt.Sprintf("Error: %v", err))
	}

	// Also update current process environment so subsequent steps in this run can use node/npm immediately
	_ = os.Setenv("NODE_HOME", targetDir)
	_ = os.Setenv("PATH", addToPath(os.Getenv("PATH"), targetDir))

	// Broadcast environment change so future processes can pick up updated user env without reboot
	if err := broadcastWindowsEnvChange(); err != nil {
		logger.Warning("⚠️ Failed to broadcast environment change", fmt.Sprintf("Error: %v", err))
	}

	logger.Success("✅ Node.js installed on Windows.")
	return nil
}

// unzip extracts a zip archive to the destination directory. Overwrites existing files.
func unzip(srcZip, destDir string) error {
	r, err := zip.OpenReader(srcZip)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		// Resolve path and prevent ZipSlip
		fpath := filepath.Join(destDir, f.Name)
		if !strings.HasPrefix(filepath.Clean(fpath)+string(os.PathSeparator), filepath.Clean(destDir)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path in zip: %s", f.Name)
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(fpath, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(fpath), 0o755); err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		out, err := os.OpenFile(fpath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
		if err != nil {
			rc.Close()
			return err
		}
		if _, err := io.Copy(out, rc); err != nil {
			out.Close()
			rc.Close()
			return err
		}
		out.Close()
		rc.Close()
	}
	return nil
}

// flattenIfSingleSubdir moves contents up if destDir contains exactly one subdirectory and no files.
func flattenIfSingleSubdir(destDir string) error {
	entries, err := os.ReadDir(destDir)
	if err != nil {
		return err
	}
	var dirs []os.DirEntry
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e)
		} else {
			// file at root -> do nothing
			return nil
		}
	}
	if len(dirs) != 1 {
		return nil
	}
	sub := filepath.Join(destDir, dirs[0].Name())
	// Move all items from sub up to destDir
	items, err := os.ReadDir(sub)
	if err != nil {
		return err
	}
	for _, it := range items {
		from := filepath.Join(sub, it.Name())
		to := filepath.Join(destDir, it.Name())
		if err := os.Rename(from, to); err != nil {
			// Fallback to copy if rename fails across volumes (unlikely)
			if it.IsDir() {
				if err2 := CopyDir(from, to); err2 != nil {
					return err
				}
				_ = os.RemoveAll(from)
			} else {
				if err2 := CopyFile(from, to, 0o755); err2 != nil {
					return err
				}
				_ = os.Remove(from)
			}
		}
	}
	// Remove now-empty subdir
	_ = os.RemoveAll(sub)
	return nil
}

func FindWindowsNpmFallback() string {
	return findWindowsNpmFallback()
}

func setWindowsUserEnv(name, value string) error {
	// Use PowerShell to persist user-level environment variable
	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command",
		fmt.Sprintf("[Environment]::SetEnvironmentVariable('%s','%s','User')", name, value))

	// Capture output to avoid TUI interference
	output, err := cmd.CombinedOutput()
	if err != nil && len(output) > 0 {
		lines := strings.Split(string(output), "\n")
		var outputLines []string
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				outputLines = append(outputLines, "   "+line)
			}
		}
		if len(outputLines) > 0 {
			logger.Error("PowerShell output:", strings.Join(outputLines, "\n"))
		}
	}
	return err
}

func getWindowsUserEnv(name string) (string, error) {
	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command",
		fmt.Sprintf("[Environment]::GetEnvironmentVariable('%s','User')", name))
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func ensureWindowsUserPathIncludes(dir string) error {
	existing, err := getWindowsUserEnv("Path")
	if err != nil {
		existing = os.Getenv("PATH") // fallback to process PATH
	}
	updated := addToPath(existing, dir)
	if updated == existing {
		return nil // already present
	}
	return setWindowsUserEnv("Path", updated)
}

func addToPath(pathVar, dir string) string {
	if dir == "" {
		return pathVar
	}
	// Windows PATH uses ';' separator.
	sep := ";"
	// Normalize for comparison
	target := strings.ToLower(filepath.Clean(dir))
	var parts []string
	if pathVar != "" {
		parts = strings.Split(pathVar, sep)
	}
	for _, p := range parts {
		if strings.ToLower(filepath.Clean(strings.TrimSpace(p))) == target {
			return pathVar // already included
		}
	}
	if pathVar == "" {
		return dir
	}
	if strings.HasSuffix(pathVar, sep) {
		return pathVar + dir
	}
	return pathVar + sep + dir
}

// broadcastWindowsEnvChange notifies the system that environment variables changed.
// This helps new processes see updated user env without requiring a full logoff.
func broadcastWindowsEnvChange() error {
	ps := `Add-Type @"
using System;
using System.Runtime.InteropServices;
public static class NativeMethods {
	[DllImport("user32.dll", SetLastError=true, CharSet=CharSet.Auto)]
	public static extern IntPtr SendMessageTimeout(IntPtr hWnd, uint Msg, IntPtr wParam, string lParam, uint fuFlags, uint uTimeout, out IntPtr lpdwResult);
}
"@; [IntPtr]$r=[IntPtr]::Zero; [void][NativeMethods]::SendMessageTimeout([IntPtr]0xffff, 0x1A, [IntPtr]::Zero, 'Environment', 0x0002, 5000, [ref]$r)`
	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", ps)

	// Capture output to avoid TUI interference
	output, err := cmd.CombinedOutput()
	if err != nil && len(output) > 0 {
		lines := strings.Split(string(output), "\n")
		var outputLines []string
		for _, line := range lines {
			if strings.TrimSpace(line) != "" {
				outputLines = append(outputLines, "   "+line)
			}
		}
		if len(outputLines) > 0 {
			logger.Error("PowerShell broadcast output:", strings.Join(outputLines, "\n"))
		}
	}
	return err
}
