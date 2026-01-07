// Package logging provides verbosity-based logging for PAC.
package logging

import (
	"fmt"
	"io"
	"log"
	"os"
)

var (
	currentVerbosity int
	logger           *log.Logger
	output           io.Writer = os.Stderr
)

func init() {
	logger = log.New(output, "", log.LstdFlags)
}

// SetVerbosity sets the global verbosity level.
// Higher values mean more verbose output.
// Typical levels:
//
//	-1: quiet (only errors)
//	 0: normal (warnings and above)
//	 1: verbose (info and above)
//	 2: debug (all messages)
func SetVerbosity(level int) {
	currentVerbosity = level
}

// GetVerbosity returns the current verbosity level.
func GetVerbosity() int {
	return currentVerbosity
}

// SetOutput sets the output destination for log messages.
// Defaults to os.Stderr.
func SetOutput(w io.Writer) {
	output = w
	logger = log.New(output, "", log.LstdFlags)
}

// Verboseprint prints a message if the current verbosity is at or above the required level.
func Verboseprint(requiredLevel int, message string) {
	if currentVerbosity >= requiredLevel {
		logger.Println(message)
	}
}

// Verbosef prints a formatted message if the current verbosity is at or above the required level.
func Verbosef(requiredLevel int, format string, args ...any) {
	if currentVerbosity >= requiredLevel {
		logger.Printf(format, args...)
	}
}

// Error prints an error message (always shown).
func Error(format string, args ...any) {
	logger.Printf("[PAC] ERROR: "+format, args...)
}

// Warn prints a warning message (shown at verbosity >= 0).
func Warn(format string, args ...any) {
	Verbosef(0, "[PAC] WARNING: "+format, args...)
}

// Info prints an info message (shown at verbosity >= 1).
func Info(format string, args ...any) {
	Verbosef(1, "[PAC] "+format, args...)
}

// Debug prints a debug message (shown at verbosity >= 2).
func Debug(format string, args ...any) {
	Verbosef(2, "[DEBUG] "+format, args...)
}

// Trace prints a trace message (shown at verbosity >= 3).
func Trace(format string, args ...any) {
	Verbosef(3, "[TRACE] "+format, args...)
}

// PrintAlways prints a message regardless of verbosity level.
// Used for output that should always be shown (like final results).
func PrintAlways(format string, args ...any) {
	fmt.Fprintf(output, format+"\n", args...)
}

// IsDebug returns true if debug logging is enabled (verbosity >= 2).
func IsDebug() bool {
	return currentVerbosity >= 2
}

// IsVerbose returns true if verbose logging is enabled (verbosity >= 1).
func IsVerbose() bool {
	return currentVerbosity >= 1
}

// IsQuiet returns true if quiet mode is enabled (verbosity < 0).
func IsQuiet() bool {
	return currentVerbosity < 0
}
