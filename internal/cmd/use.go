package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/nanmi/cc-portkey/internal/claude"
	"github.com/nanmi/cc-portkey/internal/config"
	"github.com/spf13/cobra"
)

var useCmd = &cobra.Command{
	Use:   "use <profile>",
	Short: "Switch to specified profile",
	Long: `Switch Claude Code to use the specified profile's configuration.

This updates ~/.claude/settings.json with the profile's base URL, API key,
and model settings.`,
	Args: cobra.ExactArgs(1),
	RunE: runUse,
}

func init() {
	rootCmd.AddCommand(useCmd)
}

func runUse(cmd *cobra.Command, args []string) error {
	return switchToProfile(args[0], false, nil)
}

// switchToProfile switches to the specified profile
// If launchClaude is true, starts Claude Code CLI after switching with given args
func switchToProfile(profileName string, launchClaude bool, claudeArgs []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	profile, ok := cfg.Profiles[profileName]
	if !ok {
		return fmt.Errorf("profile '%s' not found. Run 'cc-portkey list' to see available profiles", profileName)
	}

	// Apply profile to Claude settings
	if err := claude.ApplyProfile(&profile); err != nil {
		return fmt.Errorf("failed to apply profile: %w", err)
	}

	// Update current profile in config
	cfg.Current = profileName
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to update config: %w", err)
	}

	displayName := profile.DisplayName
	if displayName == "" {
		displayName = profileName
	}

	// Print switch confirmation with details
	fmt.Printf("%s Switched to %s (%s)\n", green("OK"), cyan(profileName), displayName)
	fmt.Println()

	// Expand env vars for display
	expandedURL := config.ExpandEnv(profile.BaseURL)
	expandedKey := config.ExpandEnv(profile.APIKey)

	// Show Base URL
	if expandedURL != "" {
		fmt.Printf("  Base URL:  %s\n", expandedURL)
	} else {
		fmt.Printf("  Base URL:  %s\n", cyan("https://api.anthropic.com (Official)"))
	}

	// Show masked API Key
	maskedKey := config.MaskAPIKey(expandedKey)
	fmt.Printf("  API Key:   %s\n", maskedKey)

	// Show model if configured
	if model, ok := profile.Models["default"]; ok && model != "" {
		fmt.Printf("  Model:     %s\n", model)
	}

	// Launch Claude Code CLI if requested
	if launchClaude {
		fmt.Println()
		fmt.Printf("Starting Claude Code...\n\n")
		return launchClaudeCLI(claudeArgs)
	}

	return nil
}

// launchClaudeCLI starts the Claude Code CLI, replacing the current process
func launchClaudeCLI(claudeArgs []string) error {
	claudePath, err := exec.LookPath("claude")
	if err != nil {
		return fmt.Errorf("claude command not found. Is Claude Code CLI installed?")
	}

	// Build arguments: "claude" followed by any additional args
	args := []string{"claude"}
	if claudeArgs != nil {
		args = append(args, claudeArgs...)
	}

	// Replace current process with claude (exec)
	return syscall.Exec(claudePath, args, os.Environ())
}
