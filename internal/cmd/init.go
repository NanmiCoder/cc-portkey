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
	if config.Exists() {
		path, _ := config.ConfigPath()
		fmt.Printf("%s Configuration already exists at %s\n", yellow("Warning:"), path)
		fmt.Println("Use 'cc-portkey edit' to modify it or delete the file to reinitialize.")
		return nil
	}

	cfg := config.DefaultConfig()
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	path, _ := config.ConfigPath()
	fmt.Printf("%s Configuration created at %s\n", green("Success!"), path)
	fmt.Println()

	// Auto-create symlinks
	fmt.Println("Creating shortcut commands...")
	if err := CreateSymlinks("", true); err != nil {
		fmt.Printf("%s Failed to create symlinks: %v\n", yellow("Warning:"), err)
		fmt.Println("You can try again later with: cc-portkey link")
	}

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

	return nil
}
