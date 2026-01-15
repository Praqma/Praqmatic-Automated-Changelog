package task

import (
	"regexp"
	"strings"

	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/config"
)

// TaskMatch represents a matched task ID with its associated label.
type TaskMatch struct {
	TaskID string
	Label  string
}

// ExtractTaskIDs extracts task IDs from a commit message using the configured regex patterns.
// Returns a slice of TaskMatch containing the task ID and associated label.
func ExtractTaskIDs(message string, regexConfigs []config.RegexConfig) []TaskMatch {
	var matches []TaskMatch
	seen := make(map[string]bool) // Avoid duplicate task IDs

	for _, regexCfg := range regexConfigs {
		pattern, flags := parseRubyRegex(regexCfg.Pattern)
		if pattern == "" {
			continue
		}

		// Build Go regex with flags
		var re *regexp.Regexp
		var err error

		if strings.Contains(flags, "i") {
			// Case-insensitive
			re, err = regexp.Compile("(?i)" + pattern)
		} else {
			re, err = regexp.Compile(pattern)
		}

		if err != nil {
			continue // Skip invalid patterns
		}

		// Find all matches
		allMatches := re.FindAllStringSubmatch(message, -1)
		for _, match := range allMatches {
			if len(match) > 1 {
				taskID := match[1] // First capture group
				if !seen[taskID] {
					seen[taskID] = true
					matches = append(matches, TaskMatch{
						TaskID: taskID,
						Label:  regexCfg.Label,
					})
				}
			}
		}
	}

	return matches
}

// parseRubyRegex parses a Ruby-style regex pattern: /pattern/flags
// Returns the pattern and flags separately.
func parseRubyRegex(rubyPattern string) (pattern, flags string) {
	rubyPattern = strings.TrimSpace(rubyPattern)

	// Check if it's in Ruby format: /pattern/flags
	if !strings.HasPrefix(rubyPattern, "/") {
		// Not Ruby format, use as-is
		return rubyPattern, ""
	}

	// Find the last slash (separates pattern from flags)
	lastSlash := strings.LastIndex(rubyPattern, "/")
	if lastSlash <= 0 {
		// Only one slash or slash at start, invalid
		return rubyPattern, ""
	}

	pattern = rubyPattern[1:lastSlash]
	flags = rubyPattern[lastSlash+1:]

	return pattern, flags
}

// SplitByDelimiter splits a string by the configured delimiter pattern.
// The delimiter can be a regex pattern in Ruby format (e.g., '/,|\s/').
func SplitByDelimiter(input, delimiterPattern string) []string {
	if delimiterPattern == "" {
		return []string{input}
	}

	pattern, _ := parseRubyRegex(delimiterPattern)
	if pattern == "" {
		pattern = delimiterPattern
	}

	re, err := regexp.Compile(pattern)
	if err != nil {
		return []string{input}
	}

	parts := re.Split(input, -1)

	// Filter out empty strings
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
