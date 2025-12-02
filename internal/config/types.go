package config

// Profile represents a single provider configuration
type Profile struct {
	DisplayName string            `json:"display_name"`
	BaseURL     string            `json:"base_url"`
	APIKey      string            `json:"api_key"`
	TimeoutMS   int               `json:"timeout_ms,omitempty"`
	Models      map[string]string `json:"models,omitempty"`
}

// Config represents the main configuration file structure
type Config struct {
	Current  string             `json:"current"`
	Profiles map[string]Profile `json:"profiles"`
}

// AliasMapping maps short aliases to profile names
var AliasMapping = map[string]string{
	"ccc": "claude", // ccc = Claude Code CLI (避免与 C 编译器 cc 冲突)
	"ds":  "deepseek",
	"glm": "glm",
	"mm":  "minimax",
}

// DefaultConfig returns a default configuration with common providers
func DefaultConfig() *Config {
	return &Config{
		Current: "claude",
		Profiles: map[string]Profile{
			"claude": {
				DisplayName: "Claude",
				BaseURL:     "${ANTHROPIC_BASE_URL}", // 支持中转，留空则使用官方地址
				APIKey:      "${ANTHROPIC_API_KEY}",
				TimeoutMS:   120000,
				Models:      map[string]string{},
			},
			"deepseek": {
				DisplayName: "DeepSeek",
				BaseURL:     "https://api.deepseek.com/anthropic",
				APIKey:      "${DEEPSEEK_API_KEY}",
				TimeoutMS:   600000,
				Models: map[string]string{
					"default":    "deepseek-chat",
					"small_fast": "deepseek-chat",
				},
			},
			"glm": {
				DisplayName: "GLM (Zhipu)",
				BaseURL:     "https://open.bigmodel.cn/api/anthropic",
				APIKey:      "${GLM_API_KEY}",
				TimeoutMS:   3000000,
				Models: map[string]string{
					"opus":   "glm-4.6",
					"sonnet": "glm-4.6",
					"haiku":  "glm-4.5-air",
				},
			},
			"minimax": {
				DisplayName: "MiniMax",
				BaseURL:     "https://api.minimaxi.com/anthropic",
				APIKey:      "${MINIMAX_API_KEY}",
				TimeoutMS:   3000000,
				Models: map[string]string{
					"default":    "MiniMax-M2",
					"small_fast": "MiniMax-M2",
					"opus":       "MiniMax-M2",
					"sonnet":     "MiniMax-M2",
					"haiku":      "MiniMax-M2",
				},
			},
		},
	}
}
