package main

import (
	"bytes"
	"flag"
	"io"
	"os"
	"strings"
	"testing"

	"go-template/core/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain_DefaultBehavior(t *testing.T) {
	// Reset flags for testing
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Setup command line arguments
	oldArgs := os.Args
	os.Args = []string{"cli"}
	defer func() { os.Args = oldArgs }()

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	main()

	w.Close()
	os.Stdout = old
	output := <-outC

	// Check default output
	assert.Contains(t, output, "Go Template CLI v1.0.0")
	assert.Contains(t, output, "This is a template CLI application.")
	assert.Contains(t, output, "Add your CLI commands and functionality here.")
}

func TestMain_VersionFlag(t *testing.T) {
	// Reset flags for testing
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Setup command line arguments
	oldArgs := os.Args
	os.Args = []string{"cli", "-version"}
	defer func() { os.Args = oldArgs }()

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	main()

	w.Close()
	os.Stdout = old
	output := strings.TrimSpace(<-outC)

	// Should only show version when -version flag is used
	assert.Equal(t, "Go Template CLI v1.0.0", output)
	assert.NotContains(t, output, "This is a template CLI application.")
}

func TestMain_HelpFlag(t *testing.T) {
	// Reset flags for testing
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Setup command line arguments
	oldArgs := os.Args
	os.Args = []string{"cli", "-help"}
	defer func() { os.Args = oldArgs }()

	// Capture stderr (where help is printed)
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	main()

	w.Close()
	os.Stderr = old
	output := <-outC

	// Check help output
	assert.Contains(t, output, "Usage:")
	assert.Contains(t, output, "Options:")
	assert.Contains(t, output, "-help")
	assert.Contains(t, output, "-version")
}

func TestShowHelp(t *testing.T) {
	// Reset flags for testing
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Add the flags that main() adds
	flag.Bool("version", false, "Show version information")
	flag.Bool("help", false, "Show help information")

	// Setup command line arguments
	oldArgs := os.Args
	os.Args = []string{"test-cli"}
	defer func() { os.Args = oldArgs }()

	// Capture stderr
	old := os.Stderr
	r, w, _ := os.Pipe()
	os.Stderr = w

	showHelp()

	w.Close()
	output, _ := io.ReadAll(r)
	os.Stderr = old

	outputStr := string(output)

	// Verify help output format
	assert.Contains(t, outputStr, "Usage: test-cli [options]")
	assert.Contains(t, outputStr, "Options:")
	assert.Contains(t, outputStr, "-help")
	assert.Contains(t, outputStr, "-version")
	assert.Contains(t, outputStr, "Show version information")
	assert.Contains(t, outputStr, "Show help information")
}

func TestMain_WithCustomConfig(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := tempDir + "/config.json"
	configJSON := `{
		"version": "3.0.0",
		"environment": "production",
		"log_level": "error", 
		"debug": false
	}`

	err := os.WriteFile(configPath, []byte(configJSON), 0644)
	require.NoError(t, err)

	// Set environment variable
	originalConfigPath := os.Getenv("CONFIG_PATH")
	defer os.Setenv("CONFIG_PATH", originalConfigPath)
	os.Setenv("CONFIG_PATH", configPath)

	// Reset flags for testing
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Setup command line arguments for version flag
	oldArgs := os.Args
	os.Args = []string{"cli", "-version"}
	defer func() { os.Args = oldArgs }()

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	main()

	w.Close()
	os.Stdout = old
	output := strings.TrimSpace(<-outC)

	// Should show custom version from config
	assert.Equal(t, "Go Template CLI v3.0.0", output)
}

func TestMain_ConfigLoadError(t *testing.T) {
	// Reset flags for testing
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	// Create an invalid config file (directory instead of file)
	tempDir := t.TempDir()
	invalidConfigPath := tempDir + "/invalid_config"
	err := os.Mkdir(invalidConfigPath, 0755)
	require.NoError(t, err)

	originalConfigPath := os.Getenv("CONFIG_PATH")
	defer os.Setenv("CONFIG_PATH", originalConfigPath)
	os.Setenv("CONFIG_PATH", invalidConfigPath)

	// Setup command line arguments
	oldArgs := os.Args
	os.Args = []string{"cli"}
	defer func() { os.Args = oldArgs }()

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	main()

	w.Close()
	os.Stdout = old
	output := <-outC

	// Should still work with default config
	assert.Contains(t, output, "Go Template CLI v1.0.0")
}

// Benchmark test for CLI main function
func BenchmarkMain_Default(b *testing.B) {
	// Redirect stdout to discard
	old := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)
	defer func() { os.Stdout = old }()

	for i := 0; i < b.N; i++ {
		b.StopTimer()

		// Reset flags for each benchmark iteration
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
		oldArgs := os.Args
		os.Args = []string{"cli"}

		b.StartTimer()

		// Benchmark the config loading part
		_, err := config.Load()
		if err != nil {
			b.Fatalf("Config load failed: %v", err)
		}

		b.StopTimer()
		os.Args = oldArgs
	}
}

// Test multiple flag combinations
func TestMain_FlagCombinations(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedStdout string
		expectedStderr string
		checkStdout    bool
		checkStderr    bool
	}{
		{
			name:           "No flags",
			args:           []string{"cli"},
			expectedStdout: "Go Template CLI v1.0.0",
			checkStdout:    true,
		},
		{
			name:           "Version flag",
			args:           []string{"cli", "-version"},
			expectedStdout: "Go Template CLI v1.0.0",
			checkStdout:    true,
		},
		{
			name:           "Help flag",
			args:           []string{"cli", "-help"},
			expectedStderr: "Usage:",
			checkStderr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags for testing
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)

			// Setup command line arguments
			oldArgs := os.Args
			os.Args = tt.args
			defer func() { os.Args = oldArgs }()

			if tt.checkStdout {
				old := os.Stdout
				r, w, _ := os.Pipe()
				os.Stdout = w

				outC := make(chan string)
				go func() {
					var buf bytes.Buffer
					io.Copy(&buf, r)
					outC <- buf.String()
				}()

				main()

				w.Close()
				os.Stdout = old
				output := <-outC

				assert.Contains(t, output, tt.expectedStdout)
			}

			if tt.checkStderr {
				// Reset flags again
				flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
				os.Args = tt.args

				old := os.Stderr
				r, w, _ := os.Pipe()
				os.Stderr = w

				outC := make(chan string)
				go func() {
					var buf bytes.Buffer
					io.Copy(&buf, r)
					outC <- buf.String()
				}()

				main()

				w.Close()
				os.Stderr = old
				output := <-outC

				assert.Contains(t, output, tt.expectedStderr)
			}
		})
	}
}
