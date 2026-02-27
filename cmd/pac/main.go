// Package main provides the entry point for the PAC CLI.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/cmd"
	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/logging"
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
	const maxExtraValues = 2

	var result []string
	i := 0
	for i < len(args) {
		if !isCredentialFlag(args[i]) {
			result = append(result, args[i])
			i++
			continue
		}

		flag := args[i]
		result = append(result, flag)
		i++

		if i >= len(args) {
			break
		}
		result = append(result, args[i])
		i++

		// Re-insert the flag before up to maxExtraValues additional bare values.
		for extra := 0; extra < maxExtraValues && i < len(args) && !strings.HasPrefix(args[i], "-"); extra++ {
			if extra == 0 {
				logging.Warn("deprecated: passing multiple values after %q is deprecated, use repeated %q flags instead (e.g. %s v1 %s v2 %s v3)", flag, flag, flag, flag, flag)
			}
			result = append(result, flag, args[i])
			i++
		}
	}
	return result
}

// isCredentialFlag reports whether arg is a credential flag (-c or --credentials).
func isCredentialFlag(arg string) bool {
	return arg == "-c" || arg == "--credentials"
}
