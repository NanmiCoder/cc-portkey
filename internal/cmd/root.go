package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/fatih/color"
	"github.com/nanmi/cc-portkey/internal/config"
	"github.com/spf13/cobra"
)

var (
	Version   = "0.1.0"
	cfgFile   string
	green     = color.New(color.FgGreen).SprintFunc()
	yellow    = color.New(color.FgYellow).SprintFunc()
	red       = color.New(color.FgRed).SprintFunc()
	cyan      = color.New(color.FgCyan).SprintFunc()
	bold      = color.New(color.Bold).SprintFunc()
)

var rootCmd = &cobra.Command{
	Use:   "cc-portkey",
	Short: "Claude Code Provider Switcher",
	Long: `Quickly switch Claude Code between providers.

Quick Start:
  cc-portkey init     # Initialize config + create shortcuts
  cc-portkey edit     # Add your API keys
  ds / glm / mm / ccc # Launch Claude Code with different providers`,
	Version: Version,
}

// Execute runs the root command
func Execute() {
	// Check if invoked via alias (cc, ds, glm, mm)
	if handleAlias() {
		return
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// handleAlias checks if the program was invoked via a shortcut alias
// and executes the corresponding 'use' command, then launches Claude Code
func handleAlias() bool {
	if len(os.Args) == 0 {
		return false
	}

	basename := filepath.Base(os.Args[0])

	// Check if basename matches any alias
	if profileName, ok := config.AliasMapping[basename]; ok {
		// Switch profile and launch Claude Code CLI with remaining arguments
		claudeArgs := os.Args[1:]
		if err := switchToProfile(profileName, true, claudeArgs); err != nil {
			fmt.Println(red("Error:"), err)
			os.Exit(1)
		}
		return true
	}

	return false
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path (default ~/.cc-portkey/config.json)")
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
}
