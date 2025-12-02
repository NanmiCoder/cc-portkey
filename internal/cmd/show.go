package cmd

import (
	"fmt"

	"github.com/nanmi/cc-portkey/internal/config"
	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:    "show [profile]",
	Short:  "Show profile details",
	Hidden: true,
	Long: `Show detailed configuration of a profile.

If no profile is specified, shows the current profile.
API keys are masked for security.`,
	Args: cobra.MaximumNArgs(1),
	RunE: runShow,
}

func init() {
	rootCmd.AddCommand(showCmd)
}

func runShow(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	profileName := cfg.Current
	if len(args) > 0 {
		profileName = args[0]
	}

	if profileName == "" {
		return fmt.Errorf("no profile specified and no current profile set")
	}

	profile, ok := cfg.Profiles[profileName]
	if !ok {
		return fmt.Errorf("profile '%s' not found", profileName)
	}

	fmt.Printf("%s\n", bold(fmt.Sprintf("Profile: %s", profileName)))
	fmt.Println()

	fmt.Printf("  Display Name:  %s\n", profile.DisplayName)

	if profile.BaseURL != "" {
		fmt.Printf("  Base URL:      %s\n", profile.BaseURL)
	} else {
		fmt.Printf("  Base URL:      %s\n", cyan("(official Claude API)"))
	}

	// Mask API key
	maskedKey := config.MaskAPIKey(profile.APIKey)
	fmt.Printf("  API Key:       %s\n", maskedKey)

	fmt.Printf("  Timeout:       %dms\n", profile.TimeoutMS)

	if len(profile.Models) > 0 {
		fmt.Printf("  Models:\n")
		for key, value := range profile.Models {
			fmt.Printf("    %-12s %s\n", key+":", value)
		}
	}

	if cfg.Current == profileName {
		fmt.Println()
		fmt.Printf("  Status:        %s\n", green("[current]"))
	}

	return nil
}
