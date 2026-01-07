package task

import (
	"testing"

	"github.com/Praqma/Praqmatic-Automated-Changelog/internal/config"
)

func TestParseRubyRegex(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedPattern string
		expectedFlags   string
	}{
		{
			name:            "simple pattern with flags",
			input:           "/Issue:\\s*(\\d+)/i",
			expectedPattern: "Issue:\\s*(\\d+)",
			expectedFlags:   "i",
		},
		{
			name:            "pattern without flags",
			input:           "/(#\\d+)/",
			expectedPattern: "(#\\d+)",
			expectedFlags:   "",
		},
		{
			name:            "pattern with multiple flags",
			input:           "/pattern/im",
			expectedPattern: "pattern",
			expectedFlags:   "im",
		},
		{
			name:            "non-ruby pattern",
			input:           "Issue:\\s*(\\d+)",
			expectedPattern: "Issue:\\s*(\\d+)",
			expectedFlags:   "",
		},
		{
			name:            "delimiter pattern",
			input:           "/,|\\s/",
			expectedPattern: ",|\\s",
			expectedFlags:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pattern, flags := parseRubyRegex(tt.input)
			if pattern != tt.expectedPattern {
				t.Errorf("pattern: got %q, want %q", pattern, tt.expectedPattern)
			}
			if flags != tt.expectedFlags {
				t.Errorf("flags: got %q, want %q", flags, tt.expectedFlags)
			}
		})
	}
}

func TestExtractTaskIDs(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		configs  []config.RegexConfig
		expected []TaskMatch
	}{
		{
			name:    "extract issue number",
			message: "Fix bug Issue: 123",
			configs: []config.RegexConfig{
				{Pattern: "/Issue:\\s*(\\d+)/i", Label: "issue"},
			},
			expected: []TaskMatch{
				{TaskID: "123", Label: "issue"},
			},
		},
		{
			name:    "extract github issue",
			message: "Fix #456 and #789",
			configs: []config.RegexConfig{
				{Pattern: "/(#\\d+)/", Label: "github"},
			},
			expected: []TaskMatch{
				{TaskID: "#456", Label: "github"},
				{TaskID: "#789", Label: "github"},
			},
		},
		{
			name:    "extract jira key",
			message: "PRJ-123: Implement feature",
			configs: []config.RegexConfig{
				{Pattern: "/(PRJ-\\d+)/i", Label: "jira"},
			},
			expected: []TaskMatch{
				{TaskID: "PRJ-123", Label: "jira"},
			},
		},
		{
			name:    "case insensitive match",
			message: "issue: 999",
			configs: []config.RegexConfig{
				{Pattern: "/Issue:\\s*(\\d+)/i", Label: "issue"},
			},
			expected: []TaskMatch{
				{TaskID: "999", Label: "issue"},
			},
		},
		{
			name:    "no match",
			message: "Regular commit message",
			configs: []config.RegexConfig{
				{Pattern: "/(#\\d+)/", Label: "github"},
			},
			expected: nil,
		},
		{
			name:    "multiple patterns",
			message: "Fix #123 for PRJ-456",
			configs: []config.RegexConfig{
				{Pattern: "/(#\\d+)/", Label: "github"},
				{Pattern: "/(PRJ-\\d+)/", Label: "jira"},
			},
			expected: []TaskMatch{
				{TaskID: "#123", Label: "github"},
				{TaskID: "PRJ-456", Label: "jira"},
			},
		},
		{
			name:    "duplicate task IDs are deduplicated",
			message: "#123 fixes #123",
			configs: []config.RegexConfig{
				{Pattern: "/(#\\d+)/", Label: "github"},
			},
			expected: []TaskMatch{
				{TaskID: "#123", Label: "github"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			matches := ExtractTaskIDs(tt.message, tt.configs)

			if len(matches) != len(tt.expected) {
				t.Fatalf("got %d matches, want %d", len(matches), len(tt.expected))
			}

			for i, match := range matches {
				if match.TaskID != tt.expected[i].TaskID {
					t.Errorf("match[%d].TaskID: got %q, want %q", i, match.TaskID, tt.expected[i].TaskID)
				}
				if match.Label != tt.expected[i].Label {
					t.Errorf("match[%d].Label: got %q, want %q", i, match.Label, tt.expected[i].Label)
				}
			}
		})
	}
}

func TestSplitByDelimiter(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		delimiter string
		expected  []string
	}{
		{
			name:      "comma delimiter",
			input:     "a,b,c",
			delimiter: "/,/",
			expected:  []string{"a", "b", "c"},
		},
		{
			name:      "comma or whitespace",
			input:     "a, b c",
			delimiter: "/,|\\s/",
			expected:  []string{"a", "b", "c"},
		},
		{
			name:      "empty delimiter",
			input:     "abc",
			delimiter: "",
			expected:  []string{"abc"},
		},
		{
			name:      "trims whitespace",
			input:     " a , b , c ",
			delimiter: "/,/",
			expected:  []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitByDelimiter(tt.input, tt.delimiter)

			if len(result) != len(tt.expected) {
				t.Fatalf("got %d parts, want %d: %v", len(result), len(tt.expected), result)
			}

			for i, part := range result {
				if part != tt.expected[i] {
					t.Errorf("part[%d]: got %q, want %q", i, part, tt.expected[i])
				}
			}
		})
	}
}
