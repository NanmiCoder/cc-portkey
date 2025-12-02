# cc-portkey

A CLI tool for quickly switching Claude Code between different model providers (Claude, DeepSeek, GLM, MiniMax, etc.)

## Features

- Switch between multiple API providers with a single command
- Support for environment variable references in API keys (`${VAR_NAME}`)
- Quick alias commands (`cc`, `ds`, `glm`, `mm`) for instant switching
- Cross-platform support (Windows, macOS, Linux)
- Preserves existing Claude Code settings when switching

## Installation

### From Source (All Platforms)

Requires [Go 1.21+](https://go.dev/dl/)

```bash
# Clone the repository
git clone https://github.com/nanmi/cc-portkey.git
cd cc-portkey

# Build
go build -o cc-portkey ./cmd/cc-portkey/

# Move to a directory in your PATH
# Linux/macOS:
sudo mv cc-portkey /usr/local/bin/

# Windows (PowerShell as Admin):
Move-Item cc-portkey.exe C:\Windows\System32\
```

### From Release Binary

Download the latest binary from [Releases](https://github.com/nanmi/cc-portkey/releases) and add it to your PATH.

## Quick Start

```bash
# 1. Initialize (creates config + shortcut commands)
cc-portkey init

# 2. Edit config to add your API keys
cc-portkey edit

# 3. Launch Claude Code with different providers
ds   # DeepSeek
glm  # GLM
mm   # MiniMax
ccc  # Claude (Official)
```

## Configuration

### Config File Location

- **Config file**: `~/.cc-portkey/config.json`
- **Claude settings**: `~/.claude/settings.json` (modified by cc-portkey)

### Config File Structure

```json
{
  "current": "claude",
  "profiles": {
    "claude": {
      "display_name": "Claude (Official)",
      "base_url": "",
      "api_key": "${ANTHROPIC_API_KEY}",
      "timeout_ms": 120000,
      "models": {}
    },
    "deepseek": {
      "display_name": "DeepSeek",
      "base_url": "https://api.deepseek.com/anthropic",
      "api_key": "${DEEPSEEK_API_KEY}",
      "timeout_ms": 600000,
      "models": {
        "default": "deepseek-chat",
        "small_fast": "deepseek-chat"
      }
    },
    "glm": {
      "display_name": "GLM (Zhipu)",
      "base_url": "https://open.bigmodel.cn/api/anthropic",
      "api_key": "${GLM_API_KEY}",
      "timeout_ms": 3000000,
      "models": {
        "opus": "glm-4.6",
        "sonnet": "glm-4.6",
        "haiku": "glm-4.5-air"
      }
    },
    "minimax": {
      "display_name": "MiniMax",
      "base_url": "https://api.minimaxi.com/anthropic",
      "api_key": "${MINIMAX_API_KEY}",
      "timeout_ms": 3000000,
      "models": {
        "default": "MiniMax-M2",
        "opus": "MiniMax-M2",
        "sonnet": "MiniMax-M2",
        "haiku": "MiniMax-M2"
      }
    }
  }
}
```

### Profile Fields

| Field | Description |
|-------|-------------|
| `display_name` | Human-readable name shown in output |
| `base_url` | API endpoint URL (empty for official Claude) |
| `api_key` | API key or `${ENV_VAR}` reference |
| `timeout_ms` | Request timeout in milliseconds |
| `models.default` | Default model name |
| `models.small_fast` | Model for quick tasks |
| `models.opus` | Model mapped to Claude Opus |
| `models.sonnet` | Model mapped to Claude Sonnet |
| `models.haiku` | Model mapped to Claude Haiku |

### Environment Variables

You can use environment variable references in your config:

```json
{
  "api_key": "${DEEPSEEK_API_KEY}"
}
```

Set the environment variable in your shell:

**Linux/macOS** (`~/.bashrc` or `~/.zshrc`):
```bash
# Claude (supports proxy)
export ANTHROPIC_API_KEY="sk-ant-xxx"
export ANTHROPIC_BASE_URL="https://your-proxy.com/v1"  # Optional, omit to use official API

# Other providers
export DEEPSEEK_API_KEY="sk-xxx"
export GLM_API_KEY="xxx"
export MINIMAX_API_KEY="xxx"
```

**Windows** (PowerShell profile or System Environment):
```powershell
$env:ANTHROPIC_API_KEY = "sk-ant-xxx"
$env:DEEPSEEK_API_KEY = "sk-xxx"
```

Or set permanently via System Properties > Environment Variables.

## Commands

### `cc-portkey init`

Initialize configuration with default profiles.

```bash
cc-portkey init
```

### `cc-portkey list`

List all configured profiles.

```bash
cc-portkey list

# Output:
# Profiles:
#
# * claude        Claude (Official)  [current]
#   deepseek      DeepSeek
#   glm           GLM (Zhipu)
#   minimax       MiniMax
```

### `cc-portkey use <profile>`

Switch to the specified profile.

```bash
cc-portkey use deepseek

# Output:
# OK Switched to deepseek (DeepSeek)
#
#   Base URL:  https://api.deepseek.com/anthropic
#   API Key:   sk-a...9xyz
#   Model:     deepseek-chat
```

### `cc-portkey current`

Show the currently active profile.

```bash
cc-portkey current
# Output: deepseek (DeepSeek)
```

### `cc-portkey show [profile]`

Show detailed configuration of a profile (API key is masked).

```bash
cc-portkey show deepseek
```

### `cc-portkey add <profile>`

Add a new profile interactively.

```bash
cc-portkey add openrouter
```

### `cc-portkey remove <profile>`

Remove a profile.

```bash
cc-portkey remove openrouter
```

### `cc-portkey edit`

Open config file in your default editor (`$EDITOR`).

```bash
cc-portkey edit
```

### `cc-portkey link [directory]`

Create shortcut symlinks for quick switching.

```bash
cc-portkey link              # Creates in ~/.local/bin/
cc-portkey link /usr/local/bin  # Custom directory
```

### `cc-portkey unlink [directory]`

Remove shortcut symlinks.

```bash
cc-portkey unlink
```

## Quick Aliases

The `init` command automatically creates these shortcuts in `~/.local/bin/`:

| Alias | Action |
|-------|--------|
| `ds` | DeepSeek |
| `glm` | GLM |
| `mm` | MiniMax |
| `ccc` | Claude (Official) |

### Setup by Platform

#### Linux/macOS

```bash
# Create symlinks
cc-portkey link

# Add ~/.local/bin to PATH if not already (add to ~/.bashrc or ~/.zshrc)
export PATH="$HOME/.local/bin:$PATH"

# Reload shell config
source ~/.bashrc  # or source ~/.zshrc
```

#### Windows

**Option 1: Symlinks (requires Admin)**

```powershell
# Run PowerShell as Administrator
cc-portkey link C:\Windows\System32
```

**Option 2: Batch files**

Create batch files in a directory in your PATH:

`ccc.bat`:
```batch
@echo off
cc-portkey use claude
```

`ds.bat`:
```batch
@echo off
cc-portkey use deepseek
```

## How It Works

cc-portkey modifies `~/.claude/settings.json` to set environment variables that Claude Code reads on startup:

```json
{
  "env": {
    "ANTHROPIC_BASE_URL": "https://api.deepseek.com/anthropic",
    "ANTHROPIC_AUTH_TOKEN": "sk-xxx",
    "API_TIMEOUT_MS": "600000",
    "ANTHROPIC_MODEL": "deepseek-chat",
    "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": "1"
  }
}
```

Your existing settings in `settings.json` are preserved.

## Adding Custom Providers

Edit your config file to add any provider with Anthropic-compatible API:

```bash
cc-portkey edit
```

Add a new profile:

```json
{
  "profiles": {
    "openrouter": {
      "display_name": "OpenRouter",
      "base_url": "https://openrouter.ai/api/v1",
      "api_key": "${OPENROUTER_API_KEY}",
      "timeout_ms": 300000,
      "models": {
        "default": "anthropic/claude-3.5-sonnet"
      }
    }
  }
}
```

## Troubleshooting

### "config file not found"

Run `cc-portkey init` to create the default configuration.

### API Key shows as `${VAR...}`

The environment variable is not set. Either:
1. Set the environment variable: `export DEEPSEEK_API_KEY=sk-xxx`
2. Or put the actual key in the config file (less secure)

### Symlinks not working on Windows

Windows requires Administrator privileges to create symlinks. Either:
1. Run PowerShell as Administrator
2. Use batch files instead (see Windows setup section)
3. Enable Developer Mode in Windows Settings

### Changes not taking effect

Claude Code reads `settings.json` on startup. If Claude Code is already running, restart it to apply changes.

## License

Apache License 2.0
