// Package main provides the entry point for the PAC CLI.
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version is set at build time via -ldflags.
var Version = "dev"

func main() {
	if err := execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func execute() error {
	rootCmd := &cobra.Command{
		Use:     "pac",
		Short:   "Praqmatic Automated Changelog",
		Long:    "Generate automated changelogs from git commits and task systems",
		Version: Version,
	}

	// Commands will be added in Phase 7
	return rootCmd.Execute()
}
