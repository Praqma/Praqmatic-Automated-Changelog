package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
)

func TestRootCommand_Help(t *testing.T) {
	// Capture output
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetArgs([]string{"--help"})

	err := rootCmd.Execute()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	output := buf.String()
	if output == "" {
		t.Error("expected help output, got empty string")
	}
}

func TestRootCommand_Version(t *testing.T) {
	SetVersion("1.2.3")
	if Version != "1.2.3" {
		t.Errorf("expected Version to be 1.2.3, got %s", Version)
	}
	if rootCmd.Version != "1.2.3" {
		t.Errorf("expected rootCmd.Version to be 1.2.3, got %s", rootCmd.Version)
	}
}

func TestFromCommand_Registered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "from <oldest-ref>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'from' command to be registered")
	}
}

func TestFromLatestTagCommand_Registered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "from-latest-tag <pattern>" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'from-latest-tag' command to be registered")
	}
}

func TestVersionCommand_Registered(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "version" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'version' command to be registered")
	}
}

func TestGlobalFlags_Defined(t *testing.T) {
	tests := []struct {
		name     string
		flagName string
	}{
		{"settings flag", "settings"},
		{"properties flag", "properties"},
		{"verbose flag", "verbose"},
		{"quiet flag", "quiet"},
		{"credentials flag", "credentials"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			flag := rootCmd.PersistentFlags().Lookup(tt.flagName)
			if flag == nil {
				t.Errorf("expected flag %s to be defined", tt.flagName)
			}
		})
	}
}

func TestFromCommand_Flags(t *testing.T) {
	var fromCmd *cobra.Command
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "from <oldest-ref>" {
			fromCmd = cmd
			break
		}
	}
	if fromCmd == nil {
		t.Fatal("from command not found")
	}

	toFlag := fromCmd.Flags().Lookup("to")
	if toFlag == nil {
		t.Error("expected --to flag to be defined on from command")
	}
}

func TestFromLatestTagCommand_Flags(t *testing.T) {
	var cmd *cobra.Command
	for _, c := range rootCmd.Commands() {
		if c.Use == "from-latest-tag <pattern>" {
			cmd = c
			break
		}
	}
	if cmd == nil {
		t.Fatal("from-latest-tag command not found")
	}

	toFlag := cmd.Flags().Lookup("to")
	if toFlag == nil {
		t.Error("expected --to flag to be defined on from-latest-tag command")
	}
}
