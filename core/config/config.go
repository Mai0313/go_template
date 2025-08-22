package config

import (
	"os"
	"os/user"
	"time"

	"github.com/denisbrodbeck/machineid"
)

// Config holds the application configuration
type Config struct {
	API             APIConfig `json:"api"`
	UserName        string    `json:"user_name"`
	ExtensionName   string    `json:"extension_name"`
	MachineID       string    `json:"machine_id"`
	InsightsVersion string    `json:"insights_version"`
	Mode            string    `json:"mode"`
}

// APIConfig holds API-related configuration
type APIConfig struct {
	Endpoint string        `json:"endpoint"`
	Timeout  time.Duration `json:"timeout"`
}

// Default returns the default configuration
func Default() *Config {
	machineID, err := machineid.ID()
	if err != nil {
		// 处理获取 machine ID 失败的情况，如果需要可以增加异常处理逻辑
		machineID = "unknown-machine-id"
	}
	userName := "unknown"
	if u, err := user.Current(); err == nil {
		userName = u.Username
	}

	// Load mode from environment or .env file
	mode := loadModeFromEnv()

	return &Config{
		API: APIConfig{
			Endpoint: "https://gaia.mediatek.inc/o11y/upload_locs",
			Timeout:  10 * time.Second,
		},
		UserName:        userName,
		ExtensionName:   "Claude-Code",
		MachineID:       machineID,
		InsightsVersion: "0.0.1",
		Mode:            mode,
	}
}

// loadModeFromEnv loads MODE from process env, allowing an optional .env file in CWD.
func loadModeFromEnv() string {
	// Best effort: read .env and export keys to process env
	if data, err := os.ReadFile(".env"); err == nil {
		lines := string(data)
		start := 0
		for i := 0; i <= len(lines); i++ {
			if i == len(lines) || lines[i] == '\n' || lines[i] == '\r' {
				line := lines[start:i]
				start = i + 1
				line = trimSpace(line)
				if line == "" || line[0] == '#' {
					continue
				}
				eq := -1
				for j := 0; j < len(line); j++ {
					if line[j] == '=' {
						eq = j
						break
					}
				}
				if eq <= 0 {
					continue
				}
				key := trimSpace(line[:eq])
				val := trimSpace(line[eq+1:])
				val = trimQuotes(val)
				if key != "" {
					_ = os.Setenv(key, val)
				}
			}
		}
	}

	mode := os.Getenv("MODE")
	if mode == "" {
		mode = os.Getenv("mode")
	}
	if mode == "" {
		mode = "STOP"
	}
	return mode
}

func trimSpace(s string) string {
	// minimal space trim to avoid importing strings
	start := 0
	for start < len(s) && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	end := len(s)
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

func trimQuotes(s string) string {
	if len(s) >= 2 {
		if (s[0] == '"' && s[len(s)-1] == '"') || (s[0] == '\'' && s[len(s)-1] == '\'') {
			return s[1 : len(s)-1]
		}
	}
	return s
}
