package cmd

import (
	"fmt"

	"github.com/nanmi/cc-portkey/internal/config"
	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:    "current",
	Short:  "Show current active profile",
	Hidden: true, // 'list' shows current
	RunE:  runCurrent,
}

func init() {
	rootCmd.AddCommand(currentCmd)
}

func runCurrent(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if cfg.Current == "" {
		fmt.Println("No profile is currently active.")
		fmt.Printf("Run %s to switch to a profile.\n", cyan("cc-portkey use <profile>"))
		return nil
	}

	profile, ok := cfg.Profiles[cfg.Current]
	if !ok {
		fmt.Printf("%s Current profile '%s' not found in config.\n", yellow("Warning:"), cfg.Current)
		return nil
	}

	displayName := profile.DisplayName
	if displayName == "" {
		displayName = cfg.Current
	}

	fmt.Printf("%s (%s)\n", cyan(cfg.Current), displayName)

	return nil
}
