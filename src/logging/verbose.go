package logging

import (
	"fmt"
	"os"
	"time"
)

var verboseEnabled bool

// SetVerbose enables or disables verbose logging
func SetVerbose(enabled bool) {
	verboseEnabled = enabled
}

// IsVerbose returns whether verbose logging is enabled
func IsVerbose() bool {
	return verboseEnabled
}

// Verbose prints a message if verbose mode is enabled
func Verbose(format string, args ...interface{}) {
	if verboseEnabled {
		timestamp := time.Now().Format("15:04:05")
		message := fmt.Sprintf(format, args...)
		fmt.Fprintf(os.Stderr, "[%s] %s\n", timestamp, message)
	}
}

// VerboseSection prints a section header if verbose mode is enabled
func VerboseSection(section string) {
	if verboseEnabled {
		fmt.Fprintf(os.Stderr, "\n=== %s ===\n", section)
	}
}
