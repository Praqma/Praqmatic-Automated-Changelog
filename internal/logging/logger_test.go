package logging

import (
	"bytes"
	"strings"
	"testing"
)

func TestSetVerbosity(t *testing.T) {
	// Reset to known state
	SetVerbosity(0)

	tests := []struct {
		level    int
		expected int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{-1, -1},
	}

	for _, tt := range tests {
		SetVerbosity(tt.level)
		if got := GetVerbosity(); got != tt.expected {
			t.Errorf("SetVerbosity(%d): GetVerbosity() = %d, want %d", tt.level, got, tt.expected)
		}
	}
}

func TestVerboseprint(t *testing.T) {
	buf := &bytes.Buffer{}
	SetOutput(buf)
	defer SetOutput(nil)

	tests := []struct {
		name          string
		verbosity     int
		requiredLevel int
		message       string
		shouldPrint   bool
	}{
		{"prints when level met", 1, 1, "test message", true},
		{"prints when level exceeded", 2, 1, "test message", true},
		{"skips when level not met", 0, 1, "test message", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf.Reset()
			SetVerbosity(tt.verbosity)

			Verboseprint(tt.requiredLevel, tt.message)

			output := buf.String()
			if tt.shouldPrint && !strings.Contains(output, tt.message) {
				t.Errorf("expected output to contain %q, got %q", tt.message, output)
			}
			if !tt.shouldPrint && strings.Contains(output, tt.message) {
				t.Errorf("expected no output, got %q", output)
			}
		})
	}
}

func TestVerbosef(t *testing.T) {
	buf := &bytes.Buffer{}
	SetOutput(buf)
	defer SetOutput(nil)

	SetVerbosity(1)
	buf.Reset()

	Verbosef(1, "formatted %s %d", "message", 42)

	output := buf.String()
	if !strings.Contains(output, "formatted message 42") {
		t.Errorf("expected formatted output, got %q", output)
	}
}

func TestError(t *testing.T) {
	buf := &bytes.Buffer{}
	SetOutput(buf)
	defer SetOutput(nil)

	// Error should print even when quiet
	SetVerbosity(-1)
	buf.Reset()

	Error("something went wrong: %s", "details")

	output := buf.String()
	if !strings.Contains(output, "ERROR") {
		t.Errorf("expected ERROR in output, got %q", output)
	}
	if !strings.Contains(output, "something went wrong") {
		t.Errorf("expected message in output, got %q", output)
	}
}

func TestWarn(t *testing.T) {
	buf := &bytes.Buffer{}
	SetOutput(buf)
	defer SetOutput(nil)

	SetVerbosity(0)
	buf.Reset()

	Warn("warning message")

	output := buf.String()
	if !strings.Contains(output, "WARNING") {
		t.Errorf("expected WARNING in output, got %q", output)
	}
}

func TestInfo(t *testing.T) {
	buf := &bytes.Buffer{}
	SetOutput(buf)
	defer SetOutput(nil)

	tests := []struct {
		verbosity   int
		shouldPrint bool
	}{
		{1, true},
		{2, true},
		{0, false},
		{-1, false},
	}

	for _, tt := range tests {
		buf.Reset()
		SetVerbosity(tt.verbosity)

		Info("info message")

		output := buf.String()
		if tt.shouldPrint && !strings.Contains(output, "info message") {
			t.Errorf("verbosity %d: expected output, got %q", tt.verbosity, output)
		}
		if !tt.shouldPrint && strings.Contains(output, "info message") {
			t.Errorf("verbosity %d: expected no output, got %q", tt.verbosity, output)
		}
	}
}

func TestDebug(t *testing.T) {
	buf := &bytes.Buffer{}
	SetOutput(buf)
	defer SetOutput(nil)

	tests := []struct {
		verbosity   int
		shouldPrint bool
	}{
		{2, true},
		{3, true},
		{1, false},
		{0, false},
	}

	for _, tt := range tests {
		buf.Reset()
		SetVerbosity(tt.verbosity)

		Debug("debug message")

		output := buf.String()
		if tt.shouldPrint && !strings.Contains(output, "DEBUG") {
			t.Errorf("verbosity %d: expected DEBUG output, got %q", tt.verbosity, output)
		}
		if !tt.shouldPrint && strings.Contains(output, "DEBUG") {
			t.Errorf("verbosity %d: expected no output, got %q", tt.verbosity, output)
		}
	}
}

func TestTrace(t *testing.T) {
	buf := &bytes.Buffer{}
	SetOutput(buf)
	defer SetOutput(nil)

	SetVerbosity(3)
	buf.Reset()

	Trace("trace message")

	output := buf.String()
	if !strings.Contains(output, "TRACE") {
		t.Errorf("expected TRACE in output, got %q", output)
	}

	// Should not print at lower verbosity
	SetVerbosity(2)
	buf.Reset()

	Trace("trace message")

	output = buf.String()
	if strings.Contains(output, "TRACE") {
		t.Errorf("expected no TRACE output at verbosity 2, got %q", output)
	}
}

func TestPrintAlways(t *testing.T) {
	buf := &bytes.Buffer{}
	SetOutput(buf)
	defer SetOutput(nil)

	// Should print even when quiet
	SetVerbosity(-10)
	buf.Reset()

	PrintAlways("always printed: %d", 42)

	output := buf.String()
	if !strings.Contains(output, "always printed: 42") {
		t.Errorf("expected message in output, got %q", output)
	}
}

func TestIsDebug(t *testing.T) {
	tests := []struct {
		verbosity int
		expected  bool
	}{
		{0, false},
		{1, false},
		{2, true},
		{3, true},
	}

	for _, tt := range tests {
		SetVerbosity(tt.verbosity)
		if got := IsDebug(); got != tt.expected {
			t.Errorf("verbosity %d: IsDebug() = %v, want %v", tt.verbosity, got, tt.expected)
		}
	}
}

func TestIsVerbose(t *testing.T) {
	tests := []struct {
		verbosity int
		expected  bool
	}{
		{0, false},
		{1, true},
		{2, true},
		{-1, false},
	}

	for _, tt := range tests {
		SetVerbosity(tt.verbosity)
		if got := IsVerbose(); got != tt.expected {
			t.Errorf("verbosity %d: IsVerbose() = %v, want %v", tt.verbosity, got, tt.expected)
		}
	}
}

func TestIsQuiet(t *testing.T) {
	tests := []struct {
		verbosity int
		expected  bool
	}{
		{0, false},
		{1, false},
		{-1, true},
		{-2, true},
	}

	for _, tt := range tests {
		SetVerbosity(tt.verbosity)
		if got := IsQuiet(); got != tt.expected {
			t.Errorf("verbosity %d: IsQuiet() = %v, want %v", tt.verbosity, got, tt.expected)
		}
	}
}
