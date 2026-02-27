// Package main provides the entry point for the PAC CLI.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/cmd"
)

// Version is set at build time via -ldflags.
var Version = "dev"

func main() {
	os.Args = normalizeLegacyCredentialArgs(os.Args)

	cmd.SetVersion(Version)

	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// normalizeLegacyCredentialArgs rewrites legacy "-c user password target"
// style arguments into the Cobra-compatible "-c user -c password -c target"
// form. This preserves backward compatibility with the old PAC CLI.
//
// Legacy credentials were passed as one -c with up to 3 positional values
// (e.g. "-c user password jira" or "-c token github"). This function
// detects consecutive bare (non-flag) values after -c and inserts the
// flag before each one, capping at 2 extra values to avoid consuming
// subcommands or positional arguments.
func normalizeLegacyCredentialArgs(args []string) []string {
	const maxExtraValues = 2 // legacy format: -c val1 val2 val3 → at most 2 extras

	var result []string
	i := 0
	for i < len(args) {
		arg := args[i]

		if arg == "-c" || arg == "--credentials" {
			result = append(result, arg)
			i++
			if i < len(args) {
				// First value (required).
				result = append(result, args[i])
				i++
			}
			// Gather up to maxExtraValues additional bare values.
			extra := 0
			for extra < maxExtraValues && i < len(args) && !strings.HasPrefix(args[i], "-") {
				result = append(result, arg, args[i])
				i++
				extra++
			}
		} else {
			result = append(result, arg)
			i++
		}
	}
	return result
}
