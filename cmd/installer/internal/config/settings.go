package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"claude_analysis/cmd/installer/internal/env"
	"claude_analysis/cmd/installer/internal/logger"
	"claude_analysis/cmd/installer/internal/platform"
)

type Settings struct {
	Env                        map[string]string `json:"env"`
	IncludeCoAuthoredBy        bool              `json:"includeCoAuthoredBy"`
	EnableAllProjectMcpServers bool              `json:"enableAllProjectMcpServers"`
	Hooks                      map[string][]Hook `json:"hooks"`
}

type Hook struct {
	Matcher string       `json:"matcher,omitempty"`
	Hooks   []HookAction `json:"hooks,omitempty"`
	// For leaf action
	Type    string `json:"type,omitempty"`
	Command string `json:"command,omitempty"`
}

type HookAction = Hook

// UpdateClaudeCodeSettings - TUI-based settings configuration with optional token parameter
func UpdateClaudeCodeSettings(token ...string) error {
	logger.Info("🔑 Updating GAISF API Key...")
	// Resolve settings path and load existing settings (if any) for merge
	homeDir, _ := os.UserHomeDir()
	targetDir := filepath.Join(homeDir, ".claude")
	hookPath := filepath.Join(homeDir, ".claude", platform.ExeName("claude_analysis"))
	target := filepath.Join(targetDir, "settings.json")

	var existingSettings *Settings

	if _, err := os.Stat(target); err == nil {
		if existingData, rerr := os.ReadFile(target); rerr == nil {
			var es Settings
			if jerr := json.Unmarshal(existingData, &es); jerr == nil {
				existingSettings = &es
			} else {
				logger.Warning("⚠️ Warning: Existing settings.json is not valid JSON; proceeding with defaults", fmt.Sprintf("Error: %v", jerr))
			}
		}
	}

	// Always use connectivity-based selection for MLOP URL via environment selection
	chosen := env.SelectAvailableURL()

	var gaisfToken string

	// Check if token is provided as parameter
	if len(token) > 0 && token[0] != "" {
		gaisfToken = token[0]
		logger.Success("🔑 Using provided token...")
	} else {
		// For now, just skip GAISF configuration
		// This will be implemented when we integrate with UI
		logger.Info("⏭️ Skipping GAISF configuration...")
	}

	// Build settings from existing (if any) and ensure unified defaults
	var settings Settings
	if existingSettings != nil {
		logger.Info("📋 Found existing settings, merging configurations...")
		settings = *existingSettings
	}
	EnsureDefaultSettings(&settings, hookPath, chosen.MLOPBaseURL, "")

	// Add custom headers if GAISF token was obtained
	if gaisfToken != "" {
		settings.Env["ANTHROPIC_CUSTOM_HEADERS"] = "api-key: " + gaisfToken
	} else {
		// Remove the header if no token provided (skip case)
		delete(settings.Env, "ANTHROPIC_CUSTOM_HEADERS")
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	// Create target directory and write settings.json
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return fmt.Errorf("failed to create %s: %w", targetDir, err)
	}
	if err := os.WriteFile(target, data, 0o644); err != nil {
		return fmt.Errorf("failed to write %s: %w", target, err)
	}
	logger.Success("✅ Settings saved successfully", fmt.Sprintf("Location: %s", target))
	return nil
}

// EnsureDefaultSettings applies unified defaults and required structure to settings.
// It is idempotent and can be called whether settings was empty or loaded from disk.
func EnsureDefaultSettings(settings *Settings, hookPath, baseURL, customHeader string) {
	if settings.Env == nil {
		settings.Env = make(map[string]string)
	}
	ApplyDefaultEnv(settings.Env, baseURL, customHeader)
	// Hard-enable flags required by the app
	settings.IncludeCoAuthoredBy = true
	settings.EnableAllProjectMcpServers = true
	// Ensure required Stop hook exists and points to provided hookPath
	if settings.Hooks == nil {
		settings.Hooks = make(map[string][]Hook)
	}
	settings.Hooks["Stop"] = []Hook{{Matcher: "*", Hooks: []Hook{{Type: "command", Command: hookPath}}}}
}

// ApplyDefaultEnv sets/overwrites the expected env defaults used by settings.json
func ApplyDefaultEnv(env map[string]string, baseURL string, customHeader string) {
	env["DISABLE_TELEMETRY"] = "1"
	env["CLAUDE_CODE_USE_BEDROCK"] = "1"
	env["ANTHROPIC_BEDROCK_BASE_URL"] = baseURL
	env["CLAUDE_CODE_SKIP_BEDROCK_AUTH"] = "1"
	env["CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC"] = "1"
	env["NODE_TLS_REJECT_UNAUTHORIZED"] = "0" // Allow self-signed certs for MLOP
	env["BASH_DEFAULT_TIMEOUT_MS"] = "36000000"
	env["BASH_MAX_TIMEOUT_MS"] = "36000000"
	env["MCP_TIMEOUT"] = "300000"      // 5 minutes
	env["MCP_TOOL_TIMEOUT"] = "300000" // 5 minutes for tool requests
	env["API_TIMEOUT_MS"] = "600000"   // 10 minutes for API requests

	// If custom header is provided, set it in the environment
	if customHeader != "" {
		env["ANTHROPIC_CUSTOM_HEADERS"] = customHeader
	}
}
