package testutil

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestingT is an interface that both *testing.T and *testing.B implement
type TestingT interface {
	TempDir() string
	Fatalf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// CaptureOutput captures stdout and stderr during function execution
func CaptureOutput(t *testing.T, fn func()) (stdout, stderr string) {
	// Capture stdout
	oldStdout := os.Stdout
	stdoutR, stdoutW, _ := os.Pipe()
	os.Stdout = stdoutW

	// Capture stderr
	oldStderr := os.Stderr
	stderrR, stderrW, _ := os.Pipe()
	os.Stderr = stderrW

	// Channel to coordinate goroutines
	done := make(chan bool, 1)

	// Read stdout
	var stdoutBuf bytes.Buffer
	go func() {
		io.Copy(&stdoutBuf, stdoutR)
		done <- true
	}()

	// Read stderr
	var stderrBuf bytes.Buffer
	go func() {
		io.Copy(&stderrBuf, stderrR)
		done <- true
	}()

	// Execute the function
	fn()

	// Close write ends
	stdoutW.Close()
	stderrW.Close()

	// Wait for readers to finish
	<-done
	<-done

	// Restore original stdout and stderr
	os.Stdout = oldStdout
	os.Stderr = oldStderr

	return stdoutBuf.String(), stderrBuf.String()
}

// CreateTempConfig creates a temporary config file with the given content
func CreateTempConfig(t TestingT, content string) (string, func()) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	err := os.WriteFile(configPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp config: %v", err)
	}

	// Return path and cleanup function
	return configPath, func() {
		os.Remove(configPath)
	}
}

// SetEnv sets an environment variable for testing and returns a cleanup function
func SetEnv(t TestingT, key, value string) func() {
	original := os.Getenv(key)
	os.Setenv(key, value)

	return func() {
		if original == "" {
			os.Unsetenv(key)
		} else {
			os.Setenv(key, original)
		}
	}
}

// WithTempDir creates a temporary directory for testing
func WithTempDir(t *testing.T, fn func(dir string)) {
	tempDir := t.TempDir()
	fn(tempDir)
}

// CreateInvalidConfig creates an invalid config file for error testing
func CreateInvalidConfig(t TestingT, invalidJSON string) (string, func()) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	err := os.WriteFile(configPath, []byte(invalidJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid config: %v", err)
	}

	return configPath, func() {
		os.Remove(configPath)
	}
}

// AssertFileExists checks if a file exists
func AssertFileExists(t *testing.T, path string) {
	_, err := os.Stat(path)
	require.NoError(t, err, "file should exist: %s", path)
}

// AssertFileNotExists checks if a file does not exist
func AssertFileNotExists(t *testing.T, path string) {
	_, err := os.Stat(path)
	require.True(t, os.IsNotExist(err), "file should not exist: %s", path)
}
