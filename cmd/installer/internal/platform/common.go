package platform

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"claude_analysis/cmd/installer/internal/logger"
)

// Platform utilities

func IsCommandAvailable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func ExeName(base string) string {
	if runtime.GOOS == "windows" {
		return base + ".exe"
	}
	return base
}

func PlatformSuffix() string {
	arch := runtime.GOARCH
	osname := runtime.GOOS
	switch osname {
	case "darwin", "linux", "windows":
		// ok
	default:
		// Fallback to generic
		return osname + "-" + arch
	}
	return osname + "-" + arch
}

func RunLoggedCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	// Ensure color is disabled for child processes that honor NO_COLOR
	cmd.Env = append(os.Environ(), "NO_COLOR=1")

	// Capture output instead of direct piping to avoid TUI interference
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("❌ Command failed", fmt.Sprintf("Command: %s %v\nError: %v", name, args, err))
		if len(output) > 0 {
			// Only show output if there's an error and we have output to show
			lines := strings.Split(string(output), "\n")
			var outputLines []string
			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					outputLines = append(outputLines, "   "+line)
				}
			}
			if len(outputLines) > 0 {
				logger.Error("Command output:", strings.Join(outputLines, "\n"))
			}
		}
	}
	return err
}

func RunLoggedShell(script string) error {
	cmd := exec.Command("sh", "-lc", script)
	// Ensure color is disabled for child processes that honor NO_COLOR
	cmd.Env = append(os.Environ(), "NO_COLOR=1")

	// Capture output instead of direct piping to avoid TUI interference
	output, err := cmd.CombinedOutput()
	if err != nil {
		logger.Error("❌ Shell script failed", fmt.Sprintf("Script: %s\nError: %v", script, err))
		if len(output) > 0 {
			// Only show output if there's an error and we have output to show
			lines := strings.Split(string(output), "\n")
			var outputLines []string
			for _, line := range lines {
				if strings.TrimSpace(line) != "" {
					outputLines = append(outputLines, "   "+line)
				}
			}
			if len(outputLines) > 0 {
				logger.Error("Script output:", strings.Join(outputLines, "\n"))
			}
		}
	}
	return err
}

func CopyFile(src, dst string, mode os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, mode)
	if err != nil {
		return err
	}
	defer func() { _ = out.Close() }()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return nil
}

func CopyDir(src, dst string) error {
	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		s := filepath.Join(src, e.Name())
		d := filepath.Join(dst, e.Name())
		if e.IsDir() {
			if err := CopyDir(s, d); err != nil {
				return err
			}
		} else {
			if err := CopyFile(s, d, 0o755); err != nil {
				return err
			}
		}
	}
	return nil
}

// Check Node.js version - returns true if Node.js is installed and version >= 22
func CheckNodeVersion() bool {
	if !IsCommandAvailable("node") {
		return false
	}

	out, err := exec.Command("node", "--version").Output()
	if err != nil {
		return false
	}

	version := strings.TrimSpace(string(out))
	// Remove 'v' prefix if present (e.g., "v22.1.0" -> "22.1.0")
	version = strings.TrimPrefix(version, "v")

	// Extract major version
	parts := strings.Split(version, ".")
	if len(parts) == 0 {
		return false
	}

	// Parse major version
	var major int
	if _, err := fmt.Sscanf(parts[0], "%d", &major); err != nil {
		return false
	}

	return major >= 22
}

// Find Claude binary - attempts to locate the claude CLI either on PATH or in npm's global bin directory
func FindClaudeBinary() (string, bool) {
	if p, err := exec.LookPath("claude"); err == nil {
		return p, true
	}

	npmPath := GetNpmPath()
	out, err := exec.Command(npmPath, "bin", "-g").Output()
	if err != nil {
		return "", false
	}
	binDir := strings.TrimSpace(string(out))
	p := filepath.Join(binDir, ExeName("claude"))
	if _, err := os.Stat(p); err == nil {
		return p, true
	}
	return "", false
}

func GetNpmPath() string {
	// Prefer npm next to node if node is found
	if p, err := exec.LookPath("npm"); err == nil {
		return p
	}
	// Windows-specific fallback to standard installation directory
	if runtime.GOOS == "windows" {
		// Prefer our managed install directory under user's home first
		if dir, err := getNodeInstallDir(); err == nil {
			candidate := filepath.Join(dir, "npm.cmd")
			if _, err := os.Stat(candidate); err == nil {
				return candidate
			}
		}
		// Also check one-level deeper if extracted into a versioned folder under either base
		if p := findWindowsNpmFallback(); p != "" {
			return p
		}
	}
	return "npm" // rely on PATH
}

func getNodeInstallDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("unable to resolve user home directory: %w", err)
	}
	home = strings.TrimSpace(home)
	if home == "" {
		return "", fmt.Errorf("user home directory is empty")
	}
	return filepath.Join(home, ".claude", "nodejs"), nil
}

func findWindowsNpmFallback() string {
	var bases []string
	if dir, err := getNodeInstallDir(); err == nil {
		bases = append(bases, dir)
	}
	bases = append(bases, `C:\Program Files\nodejs`)
	for _, base := range bases {
		entries, err := os.ReadDir(base)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if e.IsDir() {
				p := filepath.Join(base, e.Name(), "npm.cmd")
				if _, err := os.Stat(p); err == nil {
					return p
				}
			}
		}
	}
	return ""
}

// addToPathUnix adds a directory to PATH for Unix-like systems
func addToPathUnix(pathVar, dir string) string {
	if dir == "" {
		return pathVar
	}
	// Unix PATH uses ':' separator
	sep := ":"
	target := filepath.Clean(dir)
	var parts []string
	if pathVar != "" {
		parts = strings.Split(pathVar, sep)
	}
	for _, p := range parts {
		if filepath.Clean(strings.TrimSpace(p)) == target {
			return pathVar // already included
		}
	}
	if pathVar == "" {
		return dir
	}
	return dir + sep + pathVar // Prepend to PATH for priority
}
