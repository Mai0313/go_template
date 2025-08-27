package main

import (
	"bytes"
	"io"
	"os"
	"testing"

	"go-template/core/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain_Output(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	// Copy the output in a separate goroutine so reading doesn't block
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	// Run main function
	main()

	// Close the writer side of the pipe and get the output
	w.Close()
	os.Stdout = old
	output := <-outC

	// Check that version information is printed
	assert.Contains(t, output, "Go Template App v")
	assert.Contains(t, output, "This is a template application.")
	assert.Contains(t, output, "Replace this with your actual application logic.")
}

func TestMain_WithConfigFile(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := tempDir + "/config.json"
	configJSON := `{
		"version": "2.5.0",
		"environment": "test",
		"log_level": "debug",
		"debug": true
	}`

	err := os.WriteFile(configPath, []byte(configJSON), 0644)
	require.NoError(t, err)

	// Set environment variable
	originalConfigPath := os.Getenv("CONFIG_PATH")
	defer os.Setenv("CONFIG_PATH", originalConfigPath)
	os.Setenv("CONFIG_PATH", configPath)

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

	// Check that the custom version from config is used
	assert.Contains(t, output, "Go Template App v2.5.0")
}

func TestMain_ConfigLoadError(t *testing.T) {
	// Set invalid config path to trigger error
	originalConfigPath := os.Getenv("CONFIG_PATH")
	defer os.Setenv("CONFIG_PATH", originalConfigPath)

	// Create an invalid config file (directory instead of file)
	tempDir := t.TempDir()
	invalidConfigPath := tempDir + "/invalid_config"
	err := os.Mkdir(invalidConfigPath, 0755)
	require.NoError(t, err)

	os.Setenv("CONFIG_PATH", invalidConfigPath)

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
	assert.Contains(t, output, "Go Template App v1.0.0")
}

// Benchmark test for main function
func BenchmarkMain(b *testing.B) {
	// Redirect stdout to discard
	old := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)
	defer func() { os.Stdout = old }()

	for i := 0; i < b.N; i++ {
		// We can't benchmark main() directly due to os.Exit behavior
		// Instead we benchmark the core logic
		b.StopTimer()

		// This is a simplified benchmark that focuses on the config loading part
		// which is the main computational part of the application

		b.StartTimer()
		// Simulate the main operations without calling main()
		_, err := config.Load()
		if err != nil {
			b.Fatalf("Config load failed: %v", err)
		}
	}
}

// Test the application behavior with different environment variables
func TestMain_EnvironmentVariables(t *testing.T) {
	tests := []struct {
		name           string
		configPath     string
		expectedOutput []string
	}{
		{
			name:           "Default config path",
			configPath:     "",
			expectedOutput: []string{"Go Template App v1.0.0"},
		},
		{
			name:           "Custom config path (non-existent)",
			configPath:     "/tmp/non-existent-config.json",
			expectedOutput: []string{"Go Template App v1.0.0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalConfigPath := os.Getenv("CONFIG_PATH")
			defer os.Setenv("CONFIG_PATH", originalConfigPath)
			os.Setenv("CONFIG_PATH", tt.configPath)

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

			for _, expected := range tt.expectedOutput {
				assert.Contains(t, output, expected)
			}
		})
	}
}
