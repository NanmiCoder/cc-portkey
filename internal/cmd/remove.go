package cmd

import (
	"fmt"

	"github.com/nanmi/cc-portkey/internal/config"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove <profile>",
	Aliases: []string{"rm", "delete"},
	Short:   "Remove a profile",
	Hidden:  true, // Use 'edit' instead
	Args:    cobra.ExactArgs(1),
	RunE:    runRemove,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}

func runRemove(cmd *cobra.Command, args []string) error {
	profileName := args[0]

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if _, exists := cfg.Profiles[profileName]; !exists {
		return fmt.Errorf("profile '%s' not found", profileName)
	}

	delete(cfg.Profiles, profileName)

	// If we deleted the current profile, clear current
	if cfg.Current == profileName {
		cfg.Current = ""
		fmt.Printf("%s Profile '%s' was the current profile. No profile is now active.\n", yellow("Note:"), profileName)
	}

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("%s Profile '%s' removed.\n", green("OK"), profileName)

	return nil
}
