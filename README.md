# cc-portkey

![cc-portkey logo](assets/logo.png)

[English](./README_EN.md) | 中文

一个用于快速切换 Claude Code 模型服务商的命令行工具（支持 Claude、DeepSeek、GLM、MiniMax 等）



> **为什么叫 Portkey?**
>
> Portkey（门钥匙）是《哈利·波特》中的魔法物品——碰一下就能瞬间传送到另一个地方。
>
> 这个工具也一样：敲一个命令，就能把 Claude Code 瞬移到另一个模型服务商。

## 功能特性

- 一条命令切换不同的 API 服务商
- 支持环境变量引用 (`${VAR_NAME}`)
- 快捷别名命令 (`ccc`, `ds`, `glm`, `mm`) 秒切
- 跨平台支持 (Windows, macOS, Linux)
- 切换时保留现有的 Claude Code 配置

## 安装

### 下载预编译二进制

从 [Releases](https://github.com/NanmiCoder/cc-portkey/releases) 下载对应平台的二进制文件，放到 PATH 目录中。

### 从源码编译（所有平台）

需要 [Go 1.21+](https://go.dev/dl/)

```bash
# 克隆仓库
git clone https://github.com/NanmiCoder/cc-portkey.git
cd cc-portkey

# 编译
go build -o cc-portkey ./cmd/cc-portkey/

# 移动到 PATH 目录
# Linux/macOS:
sudo mv cc-portkey /usr/local/bin/

# Windows (管理员 PowerShell):
Move-Item cc-portkey.exe C:\Windows\System32\
```

## 快速开始

```bash
# 1. 初始化（自动创建配置 + 快捷命令）
cc-portkey init

# 2. 编辑配置，添加你的 API Keys
cc-portkey edit

# 3. 一键启动（切换 + 启动 Claude Code）
ds   # DeepSeek
glm  # GLM (智谱)
mm   # MiniMax
ccc  # Claude (官方)
```

## 配置说明

### 配置文件位置

- **cc-portkey 配置**: `~/.cc-portkey/config.json`
- **Claude 配置**: `~/.claude/settings.json`（由 cc-portkey 修改）

### 配置文件结构

```json
{
  "current": "claude",
  "profiles": {
    "claude": {
      "display_name": "Claude (官方)",
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
      "display_name": "GLM (智谱)",
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

### Profile 字段说明

| 字段 | 说明 |
|------|------|
| `display_name` | 显示名称 |
| `base_url` | API 地址（官方 Claude 留空） |
| `api_key` | API Key 或 `${环境变量}` 引用 |
| `timeout_ms` | 请求超时（毫秒） |
| `models.default` | 默认模型 |
| `models.small_fast` | 快速任务模型 |
| `models.opus` | 映射到 Claude Opus |
| `models.sonnet` | 映射到 Claude Sonnet |
| `models.haiku` | 映射到 Claude Haiku |

### 环境变量配置

配置文件支持环境变量引用：

```json
{
  "api_key": "${DEEPSEEK_API_KEY}"
}
```

**Linux/macOS** (添加到 `~/.bashrc` 或 `~/.zshrc`):
```bash
# Claude (支持中转)
export ANTHROPIC_API_KEY="sk-ant-xxx"
export ANTHROPIC_BASE_URL="https://your-proxy.com/v1"  # 可选，不设置则用官方

# 其他服务商
export DEEPSEEK_API_KEY="sk-xxx"
export GLM_API_KEY="xxx"
export MINIMAX_API_KEY="xxx"
```

**Windows** (PowerShell 或系统环境变量):
```powershell
# 临时设置
$env:DEEPSEEK_API_KEY = "sk-xxx"

# 永久设置（需要管理员权限）
[Environment]::SetEnvironmentVariable("DEEPSEEK_API_KEY", "sk-xxx", "User")
```

或通过 系统属性 > 环境变量 设置。

## 命令列表

### `cc-portkey init`

初始化配置文件，生成默认 profiles。

```bash
cc-portkey init
```

### `cc-portkey list`

列出所有配置的 profiles。

```bash
cc-portkey list

# 输出:
# Profiles:
#
# * claude        Claude (Official)  [current]
#   deepseek      DeepSeek
#   glm           GLM (Zhipu)
#   minimax       MiniMax
```

### `cc-portkey use <profile>`

切换到指定的 profile。

```bash
cc-portkey use deepseek

# 输出:
# OK Switched to deepseek (DeepSeek)
#
#   Base URL:  https://api.deepseek.com/anthropic
#   API Key:   sk-a...9xyz
#   Model:     deepseek-chat
```

### `cc-portkey edit`

用编辑器打开配置文件（使用 `$EDITOR`）。

```bash
cc-portkey edit
```

## 快捷别名

`init` 命令会自动在 `~/.local/bin/` 创建以下快捷命令：

| 别名 | 作用 |
|------|------|
| `ds` | DeepSeek |
| `glm` | GLM (智谱) |
| `mm` | MiniMax |
| `ccc` | Claude (官方) |

### 各平台设置方法

#### Linux/macOS

```bash
# 创建符号链接
cc-portkey link

# 确保 ~/.local/bin 在 PATH 中（添加到 ~/.bashrc 或 ~/.zshrc）
export PATH="$HOME/.local/bin:$PATH"

# 重新加载配置
source ~/.bashrc  # 或 source ~/.zshrc
```

#### Windows

**方式一：符号链接（需要管理员权限）**

```powershell
# 以管理员身份运行 PowerShell
cc-portkey link C:\Windows\System32
```

**方式二：批处理文件**

在 PATH 目录中创建批处理文件：

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

## 工作原理

cc-portkey 通过修改 `~/.claude/settings.json` 的 `env` 字段来切换配置：

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

你在 `settings.json` 中的其他配置会被保留。

## 添加自定义服务商

编辑配置文件添加任何兼容 Anthropic API 的服务商：

```bash
cc-portkey edit
```

添加新的 profile：

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

## 常见问题

### "config file not found"

运行 `cc-portkey init` 创建默认配置。

### API Key 显示为 `${VAR...}`

环境变量未设置。两种解决方式：
1. 设置环境变量：`export DEEPSEEK_API_KEY=sk-xxx`
2. 直接在配置文件中写入 API Key（安全性较低）

### Windows 上符号链接不工作

Windows 创建符号链接需要管理员权限。解决方式：
1. 以管理员身份运行 PowerShell
2. 使用批处理文件代替（见上方 Windows 设置）
3. 在 Windows 设置中开启开发者模式

### 切换后不生效

Claude Code 在启动时读取 `settings.json`。如果 Claude Code 正在运行，需要重启它才能应用新配置。

## 许可证

Apache License 2.0
