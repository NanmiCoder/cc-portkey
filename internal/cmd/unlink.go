package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/nanmi/cc-portkey/internal/config"
	"github.com/spf13/cobra"
)

var unlinkCmd = &cobra.Command{
	Use:    "unlink [directory]",
	Short:  "Remove shortcut symlinks",
	Hidden: true,
	Long: `Remove the symbolic links created by 'cc-portkey link'.

Removes symlinks from the specified directory (default ~/.local/bin/):
  cc, ds, glm, mm`,
	Args: cobra.MaximumNArgs(1),
	RunE: runUnlink,
}

func init() {
	rootCmd.AddCommand(unlinkCmd)
}

func runUnlink(cmd *cobra.Command, args []string) error {
	// Determine target directory
	var targetDir string
	if len(args) > 0 {
		targetDir = args[0]
	} else {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		targetDir = filepath.Join(home, ".local", "bin")
	}

	// Get the path to the current executable for verification
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	execPath, _ = filepath.EvalSymlinks(execPath)

	removed := 0
	notFound := 0

	for alias := range config.AliasMapping {
		linkPath := filepath.Join(targetDir, alias)

		// On Windows, add .exe extension
		if runtime.GOOS == "windows" {
			linkPath += ".exe"
		}

		// Check if link exists
		fi, err := os.Lstat(linkPath)
		if os.IsNotExist(err) {
			notFound++
			continue
		}

		// Only remove if it's a symlink
		if fi.Mode()&os.ModeSymlink == 0 {
			fmt.Printf("  %s -> %s (not a symlink, skipped)\n", alias, yellow("WARNING"))
			continue
		}

		// Optionally verify it points to cc-portkey
		target, _ := os.Readlink(linkPath)
		if target != "" && target != execPath {
			fmt.Printf("  %s -> %s (points elsewhere, skipped)\n", alias, yellow("WARNING"))
			continue
		}

		if err := os.Remove(linkPath); err != nil {
			fmt.Printf("  %s -> %s (%v)\n", alias, red("ERROR"), err)
			continue
		}

		fmt.Printf("  %s -> removed\n", alias)
		removed++
	}

	fmt.Println()
	if removed > 0 {
		fmt.Printf("%s Removed %d symlink(s) from %s\n", green("OK"), removed, targetDir)
	} else if notFound == len(config.AliasMapping) {
		fmt.Printf("No symlinks found in %s\n", targetDir)
	}

	return nil
}
