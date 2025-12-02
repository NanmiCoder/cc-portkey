package cmd

import (
	"fmt"
	"sort"

	"github.com/nanmi/cc-portkey/internal/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all configured profiles",
	RunE:    runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func runList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	if len(cfg.Profiles) == 0 {
		fmt.Println("No profiles configured.")
		fmt.Printf("Run %s to add a profile.\n", cyan("cc-portkey add <name>"))
		return nil
	}

	// Sort profile names for consistent output
	names := make([]string, 0, len(cfg.Profiles))
	for name := range cfg.Profiles {
		names = append(names, name)
	}
	sort.Strings(names)

	fmt.Println(bold("Profiles:"))
	fmt.Println()

	for _, name := range names {
		profile := cfg.Profiles[name]
		marker := "  "
		if name == cfg.Current {
			marker = green("* ")
		}

		displayName := profile.DisplayName
		if displayName == "" {
			displayName = name
		}

		// Show current marker
		if name == cfg.Current {
			fmt.Printf("%s%-12s  %s  %s\n", marker, cyan(name), displayName, yellow("[current]"))
		} else {
			fmt.Printf("%s%-12s  %s\n", marker, name, displayName)
		}
	}

	fmt.Println()
	fmt.Printf("Use %s to switch profiles.\n", cyan("cc-portkey use <profile>"))

	return nil
}
