package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/nanmi/cc-portkey/internal/config"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Open config file in editor",
	Long: `Open the configuration file in your default editor.

Uses $EDITOR environment variable, falling back to:
  - vim (Unix)
  - notepad (Windows)`,
	RunE: runEdit,
}

func init() {
	rootCmd.AddCommand(editCmd)
}

func runEdit(cmd *cobra.Command, args []string) error {
	path, err := config.ConfigPath()
	if err != nil {
		return err
	}

	if !config.Exists() {
		return fmt.Errorf("config file not found. Run 'cc-portkey init' first")
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		// Default editors
		if _, err := exec.LookPath("vim"); err == nil {
			editor = "vim"
		} else if _, err := exec.LookPath("nano"); err == nil {
			editor = "nano"
		} else if _, err := exec.LookPath("notepad"); err == nil {
			editor = "notepad"
		} else {
			return fmt.Errorf("no editor found. Set $EDITOR environment variable")
		}
	}

	execCmd := exec.Command(editor, path)
	execCmd.Stdin = os.Stdin
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr

	return execCmd.Run()
}
