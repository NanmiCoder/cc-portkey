package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nanmi/cc-portkey/internal/config"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:    "add <profile>",
	Short:  "Add a new profile interactively",
	Hidden: true, // Use 'edit' instead
	Long: `Add a new provider profile with interactive prompts.

You will be asked to provide:
  - Display name
  - Base URL
  - API key (can use ${ENV_VAR} syntax)
  - Timeout (optional)`,
	Args: cobra.ExactArgs(1),
	RunE: runAdd,
}

func init() {
	rootCmd.AddCommand(addCmd)
}

func runAdd(cmd *cobra.Command, args []string) error {
	profileName := args[0]

	cfg, err := config.Load()
	if err != nil {
		// If config doesn't exist, create a new one
		cfg = &config.Config{
			Profiles: make(map[string]config.Profile),
		}
	}

	if _, exists := cfg.Profiles[profileName]; exists {
		return fmt.Errorf("profile '%s' already exists. Use 'cc-portkey edit' to modify it", profileName)
	}

	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Adding new profile: %s\n\n", cyan(profileName))

	// Display name
	fmt.Printf("Display name [%s]: ", profileName)
	displayName, _ := reader.ReadString('\n')
	displayName = strings.TrimSpace(displayName)
	if displayName == "" {
		displayName = profileName
	}

	// Base URL
	fmt.Print("Base URL (e.g., https://api.example.com/anthropic): ")
	baseURL, _ := reader.ReadString('\n')
	baseURL = strings.TrimSpace(baseURL)

	// API Key
	fmt.Print("API Key (or ${ENV_VAR} reference): ")
	apiKey, _ := reader.ReadString('\n')
	apiKey = strings.TrimSpace(apiKey)

	// Timeout
	fmt.Print("Timeout in ms [120000]: ")
	timeoutStr, _ := reader.ReadString('\n')
	timeoutStr = strings.TrimSpace(timeoutStr)
	timeout := 120000
	if timeoutStr != "" {
		fmt.Sscanf(timeoutStr, "%d", &timeout)
	}

	profile := config.Profile{
		DisplayName: displayName,
		BaseURL:     baseURL,
		APIKey:      apiKey,
		TimeoutMS:   timeout,
		Models:      make(map[string]string),
	}

	cfg.Profiles[profileName] = profile

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Printf("\n%s Profile '%s' added successfully.\n", green("OK"), profileName)
	fmt.Printf("Run %s to start using it.\n", cyan(fmt.Sprintf("cc-portkey use %s", profileName)))

	return nil
}
