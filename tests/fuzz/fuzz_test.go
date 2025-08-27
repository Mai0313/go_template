package fuzz

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"unicode/utf8"

	"go-template/core/config"
)

// FuzzConfigParsing fuzzes the JSON parsing functionality
func FuzzConfigParsing(f *testing.F) {
	// Seed with valid JSON examples
	seeds := []string{
		`{"version":"1.0.0","environment":"dev","log_level":"info","debug":false}`,
		`{"version":"2.0.0","environment":"prod","log_level":"error","debug":true}`,
		`{}`,
		`{"version":"","environment":"","log_level":"","debug":false}`,
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, jsonData string) {
		var cfg config.Config
		err := json.Unmarshal([]byte(jsonData), &cfg)

		// We don't expect any panics, but errors are acceptable
		// for invalid JSON
		if err == nil {
			// If parsing succeeded, verify the config is reasonable
			// (no specific requirements, just shouldn't panic)
			_ = cfg.Version
			_ = cfg.Environment
			_ = cfg.LogLevel
			_ = cfg.Debug
		}
	})
}

// FuzzConfigLoad fuzzes the config loading functionality
func FuzzConfigLoad(f *testing.F) {
	// Seed with various config file contents
	seeds := []string{
		`{"version":"1.0.0","environment":"test","log_level":"debug","debug":true}`,
		`invalid json`,
		`{"version":123}`, // Wrong type
		`null`,
		`[]`,
		`{"extra_field":"value"}`,
	}

	for _, seed := range seeds {
		f.Add([]byte(seed))
	}

	f.Fuzz(func(t *testing.T, configData []byte) {
		// Create a temporary config file
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "fuzz_config.json")

		err := os.WriteFile(configPath, configData, 0644)
		if err != nil {
			t.Skip("Failed to write config file")
		}

		// Set environment variable
		originalConfigPath := os.Getenv("CONFIG_PATH")
		defer os.Setenv("CONFIG_PATH", originalConfigPath)
		os.Setenv("CONFIG_PATH", configPath)

		// This should not panic, even with invalid input
		cfg, err := config.Load()

		if err == nil {
			// If loading succeeded, the config should be valid
			if cfg.Version == "" {
				cfg.Version = "1.0.0" // Default should be set
			}
		}
		// Errors are acceptable for invalid input
	})
}

// FuzzConfigPathHandling fuzzes the config path handling
func FuzzConfigPath(f *testing.F) {
	// Seed with various path examples
	paths := []string{
		"/valid/path/config.json",
		"",
		"/",
		"relative/path",
		"../config.json",
		"/tmp/config.json",
		"config.json",
		"/non/existent/very/deep/path/config.json",
		"~/.config/app/config.json",
	}

	for _, path := range paths {
		f.Add(path)
	}

	f.Fuzz(func(t *testing.T, configPath string) {
		originalConfigPath := os.Getenv("CONFIG_PATH")
		defer os.Setenv("CONFIG_PATH", originalConfigPath)

		os.Setenv("CONFIG_PATH", configPath)

		// This should not panic, regardless of the path
		_, err := config.Load()

		// We expect this to return a default config in most cases
		// since most fuzz paths won't exist
		_ = err // Errors are acceptable
	})
}

// FuzzJSONFieldValues fuzzes individual JSON field values
func FuzzJSONFieldValues(f *testing.F) {
	// Seed with different field combinations
	seeds := []struct {
		version     string
		environment string
		logLevel    string
		debug       bool
	}{
		{"1.0.0", "production", "error", false},
		{"", "", "", true},
		{"v2.0.0-beta", "staging", "warn", false},
		{"invalid.version", "dev123", "invalid_level", true},
	}

	for _, seed := range seeds {
		f.Add(seed.version, seed.environment, seed.logLevel, seed.debug)
	}

	f.Fuzz(func(t *testing.T, version, environment, logLevel string, debug bool) {
		// Skip non-UTF8 strings to avoid comparison issues
		if !utf8.ValidString(version) || !utf8.ValidString(environment) || !utf8.ValidString(logLevel) {
			t.Skip("Skipping non-UTF8 strings")
		}

		// Skip very long strings that might cause issues
		if len(version) > 1000 || len(environment) > 1000 || len(logLevel) > 1000 {
			t.Skip("Skipping very long strings")
		}

		// Create a JSON config with fuzzed values
		configMap := map[string]interface{}{
			"version":     version,
			"environment": environment,
			"log_level":   logLevel,
			"debug":       debug,
		}

		jsonData, err := json.Marshal(configMap)
		if err != nil {
			t.Skip("Failed to marshal config")
		}

		// Try to parse it
		var cfg config.Config
		err = json.Unmarshal(jsonData, &cfg)

		if err == nil {
			// Verify fields were set correctly (only if parsing succeeded)
			if cfg.Version != version {
				t.Errorf("Version mismatch: got %q, want %q", cfg.Version, version)
			}
			if cfg.Environment != environment {
				t.Errorf("Environment mismatch: got %q, want %q", cfg.Environment, environment)
			}
			if cfg.LogLevel != logLevel {
				t.Errorf("LogLevel mismatch: got %q, want %q", cfg.LogLevel, logLevel)
			}
			if cfg.Debug != debug {
				t.Errorf("Debug mismatch: got %v, want %v", cfg.Debug, debug)
			}
		}
		// If parsing failed, that's also valid - just test that it doesn't crash
	})
}
