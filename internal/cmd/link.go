package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/nanmi/cc-portkey/internal/config"
	"github.com/spf13/cobra"
)

var linkCmd = &cobra.Command{
	Use:    "link [directory]",
	Short:  "Create shortcut symlinks (ccc, ds, glm, mm)",
	Hidden: true, // 'init' does this automatically
	Long: `Create symbolic links for quick profile switching.

Creates symlinks in the specified directory (default ~/.local/bin/):
  ccc -> cc-portkey (switches to claude)
  ds  -> cc-portkey (switches to deepseek)
  glm -> cc-portkey (switches to glm)
  mm  -> cc-portkey (switches to minimax)

After creating links, you can quickly switch profiles:
  $ ds    # switches to DeepSeek
  $ glm   # switches to GLM`,
	Args: cobra.MaximumNArgs(1),
	RunE: runLink,
}

func init() {
	rootCmd.AddCommand(linkCmd)
}

func runLink(cmd *cobra.Command, args []string) error {
	var targetDir string
	if len(args) > 0 {
		targetDir = args[0]
	} else {
		targetDir = ""
	}
	return CreateSymlinks(targetDir, true)
}

// CreateSymlinks creates symlinks for all aliases
// If targetDir is empty, uses ~/.local/bin/
// If verbose is true, prints detailed output
func CreateSymlinks(targetDir string, verbose bool) error {
	// Determine target directory
	if targetDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return fmt.Errorf("failed to get home directory: %w", err)
		}
		targetDir = filepath.Join(home, ".local", "bin")
	}

	// Create target directory if it doesn't exist
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", targetDir, err)
	}

	// Get the path to the current executable
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return fmt.Errorf("failed to resolve executable path: %w", err)
	}

	// Create symlinks for each alias
	created := 0
	skipped := 0

	for alias := range config.AliasMapping {
		linkPath := filepath.Join(targetDir, alias)

		// On Windows, add .exe extension
		if runtime.GOOS == "windows" {
			linkPath += ".exe"
		}

		// Check if link already exists
		if _, err := os.Lstat(linkPath); err == nil {
			// Link exists, check if it points to our executable
			target, err := os.Readlink(linkPath)
			if err == nil && target == execPath {
				if verbose {
					fmt.Printf("  %s -> already exists\n", alias)
				}
				skipped++
				continue
			}
			// Remove existing link/file
			if err := os.Remove(linkPath); err != nil {
				if verbose {
					fmt.Printf("  %s -> %s (failed to remove existing)\n", alias, red("ERROR"))
				}
				continue
			}
		}

		// Create symlink
		if err := os.Symlink(execPath, linkPath); err != nil {
			if verbose {
				if runtime.GOOS == "windows" {
					fmt.Printf("  %s -> %s (symlinks may require admin rights on Windows)\n", alias, red("ERROR"))
				} else {
					fmt.Printf("  %s -> %s (%v)\n", alias, red("ERROR"), err)
				}
			}
			continue
		}

		if verbose {
			fmt.Printf("  %s -> created\n", green(alias))
		}
		created++
	}

	if verbose {
		fmt.Println()
		if created > 0 {
			fmt.Printf("%s Created %d symlink(s) in %s\n", green("OK"), created, targetDir)
		}
		if skipped > 0 {
			fmt.Printf("   Skipped %d existing symlink(s)\n", skipped)
		}

		// Check if directory is in PATH
		if !isInPath(targetDir) {
			fmt.Println()
			fmt.Printf("%s %s is not in your PATH.\n", yellow("Note:"), targetDir)
			fmt.Printf("Add it to your shell config:\n")
			fmt.Printf("  %s\n", cyan(fmt.Sprintf("export PATH=\"%s:$PATH\"", targetDir)))
		}
	}

	return nil
}

// GetDefaultLinkDir returns the default directory for symlinks
func GetDefaultLinkDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".local", "bin")
}

// isInPath checks if a directory is in the PATH environment variable
func isInPath(dir string) bool {
	pathDirs := filepath.SplitList(os.Getenv("PATH"))
	for _, d := range pathDirs {
		if d == dir {
			return true
		}
	}
	return false
}
