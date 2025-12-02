# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

cc-portkey is a Go CLI tool for switching Claude Code between different model providers (Claude, DeepSeek, GLM, MiniMax, etc.). It works by modifying `~/.claude/settings.json` to set environment variables that control which API endpoint and credentials Claude Code uses.

## Common Development Commands

### Building
```bash
# Build the binary (outputs to current directory)
go build -o cc-portkey ./cmd/cc-portkey/

# Build for specific platforms (see release.yml for full matrix)
GOOS=linux GOARCH=amd64 go build -o cc-portkey-linux-amd64 ./cmd/cc-portkey/
GOOS=darwin GOARCH=arm64 go build -o cc-portkey-darwin-arm64 ./cmd/cc-portkey/
GOOS=windows GOARCH=amd64 go build -o cc-portkey-windows-amd64.exe ./cmd/cc-portkey/
```

### Development
```bash
# Run the built binary
./cc-portkey --help

# Install locally (for testing)
go install ./cmd/cc-portkey/

# Update dependencies
go mod tidy
```

### Testing
No automated tests are currently implemented. Manual testing involves:
1. Building the binary
2. Running commands against a local `~/.cc-portkey/config.json`
3. Verifying changes to `~/.claude/settings.json`

### Release Process
Releases are automated via GitHub Actions (`.github/workflows/release.yml`). To create a release:
1. Tag the commit: `git tag v0.1.0`
2. Push the tag: `git push origin v0.1.0`
3. The workflow builds binaries for all platforms (darwin-amd64, darwin-arm64, linux-amd64, linux-arm64, windows-amd64) and creates a GitHub release.

## Code Architecture

### Project Structure
```
.
├── cmd/cc-portkey/main.go          # Entry point (calls cmd.Execute())
├── internal/cmd/                   # Command implementations
│   ├── root.go                     # Root command & alias handling
│   ├── use.go                      # `use` command (profile switching)
│   ├── init.go                     # `init` command (config creation)
│   ├── list.go                     # `list` command (profile listing)
│   ├── edit.go                     # `edit` command (config editing)
│   ├── show.go                     # `show` command (profile details)
│   ├── current.go                  # `current` command (show current profile)
│   ├── link.go, unlink.go          # Symlink management
│   └── add.go, remove.go           # Profile management (unused in default)
├── internal/config/
│   ├── config.go                   # Config file I/O & utilities
│   └── types.go                    # Config structs & defaults
└── internal/claude/
    └── settings.go                 # Claude settings.json manipulation
```

### Key Components

1. **Configuration System** (`internal/config/`)
   - `Config` struct: Contains `Current` profile and `Profiles` map
   - `Profile` struct: API endpoint, credentials, timeout, model mappings
   - Default config includes Claude, DeepSeek, GLM, MiniMax profiles
   - Supports environment variable references (`${VAR_NAME}`)

2. **Command Layer** (`internal/cmd/`)
   - Uses [Cobra](https://github.com/spf13/cobra) for CLI framework
   - `root.go`: Handles shortcut aliases (`ccc`, `ds`, `glm`, `mm`) via `handleAlias()`
   - Each command file follows pattern: define `cobra.Command`, register in `init()`, implement `RunE`

3. **Claude Integration** (`internal/claude/`)
   - `Settings` type: `map[string]interface{}` to preserve unknown fields
   - `ApplyProfile()`: Maps profile settings to Claude environment variables:
     - `ANTHROPIC_BASE_URL`, `ANTHROPIC_AUTH_TOKEN`
     - `API_TIMEOUT_MS`, `ANTHROPIC_MODEL`, etc.
   - Handles JSON with trailing commas (common in hand-edited files)

### Data Flow
1. User runs `cc-portkey use <profile>` or shortcut alias (`ds`, `glm`, etc.)
2. `config.Load()` reads `~/.cc-portkey/config.json`
3. `claude.ApplyProfile()` updates `~/.claude/settings.json` env section
4. For aliases: `switchToProfile(..., launchClaude: true)` also execs `claude` CLI

### Environment Variables Mapping
Profile settings map to Claude Code environment variables:
- `profile.BaseURL` → `ANTHROPIC_BASE_URL` (deleted if empty → official API)
- `profile.APIKey` → `ANTHROPIC_AUTH_TOKEN`
- `profile.TimeoutMS` → `API_TIMEOUT_MS`
- `profile.Models["default"]` → `ANTHROPIC_MODEL`
- `profile.Models["small_fast"]` → `ANTHROPIC_SMALL_FAST_MODEL`
- `profile.Models["opus"]` → `ANTHROPIC_DEFAULT_OPUS_MODEL`
- `profile.Models["sonnet"]` → `ANTHROPIC_DEFAULT_SONNET_MODEL`
- `profile.Models["haiku"]` → `ANTHROPIC_DEFAULT_HAIKU_MODEL`

For third-party providers (`baseURL != ""`), also sets `CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC=1`.

## Configuration File Structure

Default location: `~/.cc-portkey/config.json`
```json
{
  "current": "claude",
  "profiles": {
    "claude": {
      "display_name": "Claude (Official)",
      "base_url": "${ANTHROPIC_BASE_URL}",
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
    }
    // ... other profiles
  }
}
```

## Adding New Commands

1. Create file in `internal/cmd/` (e.g., `newcmd.go`)
2. Define `cobra.Command` with `Use`, `Short`, `Long`, `RunE`
3. Register in `init()`: `rootCmd.AddCommand(newCmd)`
4. Follow existing patterns for error handling and output formatting

## Cross-Platform Considerations

- Symlink creation (`link` command) works differently on Windows (requires Admin)
- Path separators: Use `filepath.Join()` not string concatenation
- Line endings: Use `fmt.Println()` not `\n` concatenation
- Executable naming: Windows adds `.exe` extension automatically

## Dependencies

- [spf13/cobra](https://github.com/spf13/cobra): CLI framework
- [fatih/color](https://github.com/fatih/color): Colored terminal output

See `go.mod` for exact versions.