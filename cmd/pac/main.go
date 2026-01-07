// Package main provides the entry point for the PAC CLI.
package main

import (
	"fmt"
	"os"

	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/cmd"
)

// Version is set at build time via -ldflags.
var Version = "dev"

func main() {
	cmd.SetVersion(Version)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
