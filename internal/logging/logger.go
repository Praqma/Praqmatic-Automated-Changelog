// Package logging provides verbosity-based logging for PAC.package logging

package logging

import (
	"fmt"
	"log"
)

var currentVerbosity int

// SetVerbosity sets the global verbosity level.
// Higher values mean more verbose output.
func SetVerbosity(level int) {
	currentVerbosity = level
}

// GetVerbosity returns the current verbosity level.
func GetVerbosity() int {
	return currentVerbosity
}

// Verboseprint prints a message if the current verbosity is at or above the required level.
func Verboseprint(requiredLevel int, message string) {
	if currentVerbosity >= requiredLevel {
		log.Println(message)
	}
}

// Verbosef prints a formatted message if the current verbosity is at or above the required level.
func Verbosef(requiredLevel int, format string, args ...any) {
	if currentVerbosity >= requiredLevel {
		log.Printf(format, args...)
	}
}

// Debug prints a debug message (requires verbosity >= 2).
func Debug(format string, args ...any) {
	Verbosef(2, fmt.Sprintf("[DEBUG] %s", format), args...)
}

// Info prints an info message (requires verbosity >= 1).
func Info(format string, args ...any) {
	Verbosef(1, fmt.Sprintf("[PAC] %s", format), args...)
}

// Warn prints a warning message (always shown unless quiet).
func Warn(format string, args ...any) {
	Verbosef(0, fmt.Sprintf("[PAC] WARNING: %s", format), args...)
}
