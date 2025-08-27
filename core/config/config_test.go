package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad_DefaultConfig(t *testing.T) {
	// Test loading default config when no file exists
	originalConfigPath := os.Getenv("CONFIG_PATH")
	defer os.Setenv("CONFIG_PATH", originalConfigPath)

	// Set an invalid config path
	os.Setenv("CONFIG_PATH", "/non/existent/path/config.json")

	cfg, err := Load()
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	// Check default values
	assert.Equal(t, "1.0.0", cfg.Version)
	assert.Equal(t, "development", cfg.Environment)
	assert.Equal(t, "info", cfg.LogLevel)
	assert.False(t, cfg.Debug)
}

func TestLoad_FromFile(t *testing.T) {
	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	testConfig := Config{
		Version:     "2.0.0",
		Environment: "production",
		LogLevel:    "error",
		Debug:       true,
	}

	// Write test config to file
	data, err := json.Marshal(testConfig)
	require.NoError(t, err)
	err = os.WriteFile(configPath, data, 0644)
	require.NoError(t, err)

	// Set CONFIG_PATH environment variable
	originalConfigPath := os.Getenv("CONFIG_PATH")
	defer os.Setenv("CONFIG_PATH", originalConfigPath)
	os.Setenv("CONFIG_PATH", configPath)

	// Load config
	cfg, err := Load()
	require.NoError(t, err)
	assert.NotNil(t, cfg)

	// Verify loaded values
	assert.Equal(t, "2.0.0", cfg.Version)
	assert.Equal(t, "production", cfg.Environment)
	assert.Equal(t, "error", cfg.LogLevel)
	assert.True(t, cfg.Debug)
}

func TestLoadFromFile_FileNotExists(t *testing.T) {
	cfg, err := loadFromFile()
	assert.Nil(t, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config file does not exist")
}

func TestLoadFromFile_InvalidJSON(t *testing.T) {
	// Create a temporary file with invalid JSON
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	err := os.WriteFile(configPath, []byte("invalid json"), 0644)
	require.NoError(t, err)

	// Set CONFIG_PATH environment variable
	originalConfigPath := os.Getenv("CONFIG_PATH")
	defer os.Setenv("CONFIG_PATH", originalConfigPath)
	os.Setenv("CONFIG_PATH", configPath)

	cfg, err := loadFromFile()
	assert.Nil(t, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse config file")
}

func TestLoadFromFile_ReadError(t *testing.T) {
	// Create a directory instead of a file to simulate read error
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config_dir")
	err := os.Mkdir(configPath, 0755)
	require.NoError(t, err)

	// Set CONFIG_PATH environment variable
	originalConfigPath := os.Getenv("CONFIG_PATH")
	defer os.Setenv("CONFIG_PATH", originalConfigPath)
	os.Setenv("CONFIG_PATH", configPath)

	cfg, err := loadFromFile()
	assert.Nil(t, cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to read config file")
}

func TestGetConfigPath_WithEnvVar(t *testing.T) {
	originalConfigPath := os.Getenv("CONFIG_PATH")
	defer os.Setenv("CONFIG_PATH", originalConfigPath)

	expectedPath := "/custom/config/path.json"
	os.Setenv("CONFIG_PATH", expectedPath)

	path := getConfigPath()
	assert.Equal(t, expectedPath, path)
}

func TestGetConfigPath_WithoutEnvVar(t *testing.T) {
	originalConfigPath := os.Getenv("CONFIG_PATH")
	defer os.Setenv("CONFIG_PATH", originalConfigPath)

	os.Setenv("CONFIG_PATH", "")

	path := getConfigPath()

	// Should contain home directory path or fallback to ./config.json
	homeDir, err := os.UserHomeDir()
	if err != nil {
		assert.Equal(t, "./config.json", path)
	} else {
		expectedPath := filepath.Join(homeDir, ".go-template", "config.json")
		assert.Equal(t, expectedPath, path)
	}
}

func TestConfig_JSONMarshalUnmarshal(t *testing.T) {
	originalConfig := Config{
		Version:     "1.5.0",
		Environment: "staging",
		LogLevel:    "debug",
		Debug:       true,
	}

	// Marshal to JSON
	data, err := json.Marshal(originalConfig)
	require.NoError(t, err)

	// Unmarshal back
	var newConfig Config
	err = json.Unmarshal(data, &newConfig)
	require.NoError(t, err)

	// Compare
	assert.Equal(t, originalConfig, newConfig)
}

func TestLoad_Integration(t *testing.T) {
	// Integration test that simulates a complete workflow
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.json")

	// Test 1: No config file, should return default
	originalConfigPath := os.Getenv("CONFIG_PATH")
	defer os.Setenv("CONFIG_PATH", originalConfigPath)
	os.Setenv("CONFIG_PATH", configPath)

	cfg, err := Load()
	require.NoError(t, err)
	assert.Equal(t, "1.0.0", cfg.Version)

	// Test 2: Create config file and load it
	testConfig := Config{
		Version:     "3.0.0",
		Environment: "test",
		LogLevel:    "warn",
		Debug:       false,
	}

	data, err := json.Marshal(testConfig)
	require.NoError(t, err)
	err = os.WriteFile(configPath, data, 0644)
	require.NoError(t, err)

	cfg, err = Load()
	require.NoError(t, err)
	assert.Equal(t, testConfig, *cfg)
}
