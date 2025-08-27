package integration

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"go-template/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestAppIntegration tests the app binary end-to-end
func TestAppIntegration(t *testing.T) {
	// Build the app binary first
	buildDir := t.TempDir()
	appBinary := filepath.Join(buildDir, "app")

	cmd := exec.Command("go", "build", "-o", appBinary, "../../cmd/app")
	cmd.Dir = "."
	err := cmd.Run()
	require.NoError(t, err, "Failed to build app binary")

	testutil.AssertFileExists(t, appBinary)

	// Test running the app with default config
	cmd = exec.Command(appBinary)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "App should run successfully")

	outputStr := string(output)
	assert.Contains(t, outputStr, "Go Template App v1.0.0")
	assert.Contains(t, outputStr, "This is a template application.")
}

func TestAppIntegrationWithCustomConfig(t *testing.T) {
	// Build the app binary first
	buildDir := t.TempDir()
	appBinary := filepath.Join(buildDir, "app")

	cmd := exec.Command("go", "build", "-o", appBinary, "../../cmd/app")
	cmd.Dir = "."
	err := cmd.Run()
	require.NoError(t, err, "Failed to build app binary")

	// Create custom config
	configContent := `{
		"version": "2.5.0",
		"environment": "integration-test",
		"log_level": "debug",
		"debug": true
	}`

	configPath, cleanup := testutil.CreateTempConfig(t, configContent)
	defer cleanup()

	// Run app with custom config
	cmd = exec.Command(appBinary)
	cmd.Env = append(os.Environ(), "CONFIG_PATH="+configPath)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "App should run successfully with custom config")

	outputStr := string(output)
	assert.Contains(t, outputStr, "Go Template App v2.5.0")
}

// TestCLIIntegration tests the CLI binary end-to-end
func TestCLIIntegration(t *testing.T) {
	// Build the CLI binary first
	buildDir := t.TempDir()
	cliBinary := filepath.Join(buildDir, "cli")

	cmd := exec.Command("go", "build", "-o", cliBinary, "../../cmd/cli")
	cmd.Dir = "."
	err := cmd.Run()
	require.NoError(t, err, "Failed to build CLI binary")

	testutil.AssertFileExists(t, cliBinary)

	tests := []struct {
		name           string
		args           []string
		expectedOutput []string
		expectError    bool
	}{
		{
			name: "Default behavior",
			args: []string{},
			expectedOutput: []string{
				"Go Template CLI v1.0.0",
				"This is a template CLI application.",
			},
			expectError: false,
		},
		{
			name: "Version flag",
			args: []string{"-version"},
			expectedOutput: []string{
				"Go Template CLI v1.0.0",
			},
			expectError: false,
		},
		{
			name: "Help flag",
			args: []string{"-help"},
			expectedOutput: []string{
				"Usage:",
				"Options:",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command(cliBinary, tt.args...)
			output, err := cmd.CombinedOutput()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			outputStr := string(output)
			for _, expected := range tt.expectedOutput {
				assert.Contains(t, outputStr, expected)
			}
		})
	}
}

func TestCLIIntegrationWithCustomConfig(t *testing.T) {
	// Build the CLI binary first
	buildDir := t.TempDir()
	cliBinary := filepath.Join(buildDir, "cli")

	cmd := exec.Command("go", "build", "-o", cliBinary, "../../cmd/cli")
	cmd.Dir = "."
	err := cmd.Run()
	require.NoError(t, err, "Failed to build CLI binary")

	// Create custom config
	configContent := `{
		"version": "3.1.4",
		"environment": "integration-test",
		"log_level": "warn",
		"debug": false
	}`

	configPath, cleanup := testutil.CreateTempConfig(t, configContent)
	defer cleanup()

	// Test with version flag and custom config
	cmd = exec.Command(cliBinary, "-version")
	cmd.Env = append(os.Environ(), "CONFIG_PATH="+configPath)
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "CLI should run successfully with custom config")

	outputStr := strings.TrimSpace(string(output))
	assert.Equal(t, "Go Template CLI v3.1.4", outputStr)
}

// TestBuildTargets tests that all Makefile targets build successfully
func TestBuildTargets(t *testing.T) {
	// Change to project root
	projectRoot := "../../"

	// Test building individual targets
	targets := []string{"app", "cli"}

	for _, target := range targets {
		t.Run("Build_"+target, func(t *testing.T) {
			cmd := exec.Command("make", target)
			cmd.Dir = projectRoot
			output, err := cmd.CombinedOutput()

			assert.NoError(t, err, "Make %s should succeed. Output: %s", target, string(output))

			// Check that binary was created
			binaryPath := filepath.Join(projectRoot, "build", target)
			testutil.AssertFileExists(t, binaryPath)
		})
	}
}

// TestConfigFileHandling tests various config file scenarios
func TestConfigFileHandling(t *testing.T) {
	buildDir := t.TempDir()
	appBinary := filepath.Join(buildDir, "app")

	// Build the app binary
	cmd := exec.Command("go", "build", "-o", appBinary, "../../cmd/app")
	cmd.Dir = "."
	err := cmd.Run()
	require.NoError(t, err, "Failed to build app binary")

	tests := []struct {
		name           string
		configContent  string
		configPath     string
		expectError    bool
		expectedOutput string
	}{
		{
			name: "Valid config file",
			configContent: `{
				"version": "1.5.0",
				"environment": "production",
				"log_level": "error",
				"debug": false
			}`,
			expectError:    false,
			expectedOutput: "Go Template App v1.5.0",
		},
		{
			name:           "Non-existent config file",
			configPath:     "/non/existent/config.json",
			expectError:    false,
			expectedOutput: "Go Template App v1.0.0", // Should use default
		},
		{
			name:           "Invalid JSON config",
			configContent:  `{"version": "1.0.0"`, // Invalid JSON
			expectError:    false,                 // Should fallback to default
			expectedOutput: "Go Template App v1.0.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var configPath string
			var cleanup func()

			if tt.configContent != "" {
				configPath, cleanup = testutil.CreateTempConfig(t, tt.configContent)
				defer cleanup()
			} else {
				configPath = tt.configPath
			}

			cmd := exec.Command(appBinary)
			if configPath != "" {
				cmd.Env = append(os.Environ(), "CONFIG_PATH="+configPath)
			}

			output, err := cmd.CombinedOutput()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err, "Command should succeed. Output: %s", string(output))
			}

			if tt.expectedOutput != "" {
				assert.Contains(t, string(output), tt.expectedOutput)
			}
		})
	}
}

// TestConfigStructure tests that config files have the expected structure
func TestConfigStructure(t *testing.T) {
	configContent := `{
		"version": "1.2.3",
		"environment": "test",
		"log_level": "debug",
		"debug": true
	}`

	configPath, cleanup := testutil.CreateTempConfig(t, configContent)
	defer cleanup()

	// Read and parse the config file
	data, err := os.ReadFile(configPath)
	require.NoError(t, err)

	var config struct {
		Version     string `json:"version"`
		Environment string `json:"environment"`
		LogLevel    string `json:"log_level"`
		Debug       bool   `json:"debug"`
	}

	err = json.Unmarshal(data, &config)
	require.NoError(t, err)

	// Verify structure
	assert.Equal(t, "1.2.3", config.Version)
	assert.Equal(t, "test", config.Environment)
	assert.Equal(t, "debug", config.LogLevel)
	assert.True(t, config.Debug)
}

// BenchmarkAppStartup benchmarks application startup time
func BenchmarkAppStartup(b *testing.B) {
	// Build the app binary once
	buildDir := b.TempDir()
	appBinary := filepath.Join(buildDir, "app")

	cmd := exec.Command("go", "build", "-o", appBinary, "../../cmd/app")
	cmd.Dir = "."
	err := cmd.Run()
	require.NoError(b, err, "Failed to build app binary")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cmd := exec.Command(appBinary)
		err := cmd.Run()
		if err != nil {
			b.Fatalf("App execution failed: %v", err)
		}
	}
}

// BenchmarkCLIStartup benchmarks CLI startup time
func BenchmarkCLIStartup(b *testing.B) {
	// Build the CLI binary once
	buildDir := b.TempDir()
	cliBinary := filepath.Join(buildDir, "cli")

	cmd := exec.Command("go", "build", "-o", cliBinary, "../../cmd/cli")
	cmd.Dir = "."
	err := cmd.Run()
	require.NoError(b, err, "Failed to build CLI binary")

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cmd := exec.Command(cliBinary, "-version")
		err := cmd.Run()
		if err != nil {
			b.Fatalf("CLI execution failed: %v", err)
		}
	}
}
