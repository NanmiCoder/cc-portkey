package claude

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/nanmi/cc-portkey/internal/config"
)

const (
	claudeDirName      = ".claude"
	settingsFileName   = "settings.json"
)

// SettingsPath returns the path to Claude's settings.json
func SettingsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, claudeDirName, settingsFileName), nil
}

// Settings represents Claude Code's settings.json structure
// We use map[string]interface{} to preserve unknown fields
type Settings map[string]interface{}

// removeTrailingCommas removes trailing commas from JSON that would cause parse errors
// This handles cases like: {"key": "value",} -> {"key": "value"}
func removeTrailingCommas(data []byte) []byte {
	// Remove trailing commas before } or ]
	re := regexp.MustCompile(`,\s*([\]}])`)
	return re.ReplaceAll(data, []byte("$1"))
}

// Load reads Claude's settings.json
func Load() (Settings, error) {
	path, err := SettingsPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			// Return empty settings if file doesn't exist
			return Settings{}, nil
		}
		return nil, fmt.Errorf("failed to read settings file: %w", err)
	}

	// Clean up JSON with trailing commas (common in hand-edited files)
	data = removeTrailingCommas(data)

	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("failed to parse settings file: %w", err)
	}

	return settings, nil
}

// Save writes Claude's settings.json
func Save(settings Settings) error {
	path, err := SettingsPath()
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create claude config directory: %w", err)
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	// Write to temp file first for atomic operation
	tmpPath := path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write settings file: %w", err)
	}

	if err := os.Rename(tmpPath, path); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("failed to save settings file: %w", err)
	}

	return nil
}

// ApplyProfile applies a profile's settings to Claude's settings.json
func ApplyProfile(profile *config.Profile) error {
	settings, err := Load()
	if err != nil {
		return err
	}

	// Get or create env map
	env, ok := settings["env"].(map[string]interface{})
	if !ok {
		env = make(map[string]interface{})
	}

	// Expand environment variables
	apiKey := config.ExpandEnv(profile.APIKey)
	baseURL := config.ExpandEnv(profile.BaseURL)

	// Apply base_url (empty after expansion means use official API)
	if baseURL != "" {
		env["ANTHROPIC_BASE_URL"] = baseURL
	} else {
		delete(env, "ANTHROPIC_BASE_URL")
	}

	// Apply API key
	env["ANTHROPIC_AUTH_TOKEN"] = apiKey

	// Apply timeout
	if profile.TimeoutMS > 0 {
		env["API_TIMEOUT_MS"] = strconv.Itoa(profile.TimeoutMS)
	}

	// Apply models
	if model, ok := profile.Models["default"]; ok && model != "" {
		env["ANTHROPIC_MODEL"] = model
	} else {
		delete(env, "ANTHROPIC_MODEL")
	}

	if model, ok := profile.Models["small_fast"]; ok && model != "" {
		env["ANTHROPIC_SMALL_FAST_MODEL"] = model
	} else {
		delete(env, "ANTHROPIC_SMALL_FAST_MODEL")
	}

	if model, ok := profile.Models["opus"]; ok && model != "" {
		env["ANTHROPIC_DEFAULT_OPUS_MODEL"] = model
	} else {
		delete(env, "ANTHROPIC_DEFAULT_OPUS_MODEL")
	}

	if model, ok := profile.Models["sonnet"]; ok && model != "" {
		env["ANTHROPIC_DEFAULT_SONNET_MODEL"] = model
	} else {
		delete(env, "ANTHROPIC_DEFAULT_SONNET_MODEL")
	}

	if model, ok := profile.Models["haiku"]; ok && model != "" {
		env["ANTHROPIC_DEFAULT_HAIKU_MODEL"] = model
	} else {
		delete(env, "ANTHROPIC_DEFAULT_HAIKU_MODEL")
	}

	// Disable nonessential traffic for third-party providers
	if baseURL != "" {
		env["CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC"] = "1"
	} else {
		delete(env, "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC")
	}

	settings["env"] = env

	return Save(settings)
}
