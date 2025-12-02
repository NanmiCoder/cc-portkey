package cmd

import (
	"fmt"

	"github.com/nanmi/cc-portkey/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration file with default profiles",
	Long: `Initialize the cc-portkey configuration file with default profiles
for common providers (Claude, DeepSeek, GLM, MiniMax).

The configuration will be created at ~/.cc-portkey/config.json`,
	RunE: runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) error {
	path, _ := config.ConfigPath()
	configExists := config.Exists()

	if configExists {
		fmt.Printf("%s Found existing configuration at %s\n", cyan("Info:"), path)
		fmt.Println("Updating shortcut commands...")
	} else {
		// Create default configuration
		cfg := config.DefaultConfig()
		if err := config.Save(cfg); err != nil {
			return fmt.Errorf("failed to create config: %w", err)
		}
		fmt.Printf("%s Configuration created at %s\n", green("Success!"), path)
		fmt.Println()
		fmt.Println("Creating shortcut commands...")
	}

	// Create/update symlinks for all aliases
	if err := CreateSymlinks("", true); err != nil {
		fmt.Printf("%s Failed to create symlinks: %v\n", yellow("Warning:"), err)
		fmt.Println("You can try again later with: cc-portkey link")
	}

	// Only show next steps if we just created the config
	if !configExists {
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Printf("  1. Edit config to add your API keys:\n")
		fmt.Printf("     %s\n", cyan("cc-portkey edit"))
		fmt.Println()
		fmt.Printf("  2. Or set environment variables:\n")
		fmt.Printf("     %s\n", cyan("export DEEPSEEK_API_KEY=sk-xxx"))
		fmt.Println()
		fmt.Printf("  3. Launch Claude Code with different providers:\n")
		fmt.Printf("     %s  # DeepSeek\n", cyan("ds"))
		fmt.Printf("     %s  # GLM\n", cyan("glm"))
		fmt.Printf("     %s  # MiniMax\n", cyan("mm"))
		fmt.Printf("     %s # Claude (Official)\n", cyan("ccc"))
	}

	return nil
}
