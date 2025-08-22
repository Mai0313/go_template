package tests

import (
	"os"
	"path/filepath"
	"testing"

	"claude_analysis/core/config"
)

func TestMode_DefaultSTOP_WithoutEnvAndDotenv(t *testing.T) {
	// Save and clear env
	oldMODE := os.Getenv("MODE")
	oldmode := os.Getenv("mode")
	_ = os.Unsetenv("MODE")
	_ = os.Unsetenv("mode")
	defer func() {
		_ = os.Setenv("MODE", oldMODE)
		_ = os.Setenv("mode", oldmode)
	}()

	// Use empty temp dir as CWD (no .env)
	tmpDir := t.TempDir()
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	cfg := config.Default()
	if cfg.Mode != "STOP" {
		t.Errorf("expected Mode=STOP by default, got %q", cfg.Mode)
	}
}

func TestMode_FromEnvVar_MODE(t *testing.T) {
	oldMODE := os.Getenv("MODE")
	oldmode := os.Getenv("mode")
	defer func() {
		_ = os.Setenv("MODE", oldMODE)
		_ = os.Setenv("mode", oldmode)
	}()
	_ = os.Setenv("MODE", "POST_TOOL")
	_ = os.Unsetenv("mode")

	// Any dir
	cfg := config.Default()
	if cfg.Mode != "POST_TOOL" {
		t.Errorf("expected Mode=POST_TOOL from MODE env, got %q", cfg.Mode)
	}
}

func TestMode_FromDotEnv_File(t *testing.T) {
	// Clear env so .env takes effect
	oldMODE := os.Getenv("MODE")
	oldmode := os.Getenv("mode")
	_ = os.Unsetenv("MODE")
	_ = os.Unsetenv("mode")
	defer func() {
		_ = os.Setenv("MODE", oldMODE)
		_ = os.Setenv("mode", oldmode)
	}()

	tmpDir := t.TempDir()
	cwd, _ := os.Getwd()
	defer os.Chdir(cwd)
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	// Write .env with whitespace and quotes
	content := "# sample env\n  MODE = \"POST_TOOL\"  \n OTHER=abc\n"
	if err := os.WriteFile(filepath.Join(tmpDir, ".env"), []byte(content), 0o644); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	cfg := config.Default()
	if cfg.Mode != "POST_TOOL" {
		t.Errorf("expected Mode=POST_TOOL from .env, got %q", cfg.Mode)
	}
}
