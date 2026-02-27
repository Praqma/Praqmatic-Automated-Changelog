package main

import (
	"testing"
)

func Test_normalizeLegacyCredentialArgs(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want []string
	}{
		{
			name: "no credentials flag",
			args: []string{"pac", "from", "--settings", "s.yml"},
			want: []string{"pac", "from", "--settings", "s.yml"},
		},
		{
			name: "already new format with -c",
			args: []string{"pac", "from", "-c", "user", "-c", "pass", "-c", "jira"},
			want: []string{"pac", "from", "-c", "user", "-c", "pass", "-c", "jira"},
		},
		{
			name: "legacy three values after single -c",
			args: []string{"pac", "from", "-c", "user", "pass", "jira"},
			want: []string{"pac", "from", "-c", "user", "-c", "pass", "-c", "jira"},
		},
		{
			name: "legacy two values after single -c (token flow)",
			args: []string{"pac", "from", "-c", "mytoken", "github"},
			want: []string{"pac", "from", "-c", "mytoken", "-c", "github"},
		},
		{
			name: "legacy with --credentials long flag",
			args: []string{"pac", "from", "--credentials", "user", "pass", "jira"},
			want: []string{"pac", "from", "--credentials", "user", "--credentials", "pass", "--credentials", "jira"},
		},
		{
			name: "legacy -c followed by other flags",
			args: []string{"pac", "from", "-c", "user", "pass", "jira", "-s", "settings.yml"},
			want: []string{"pac", "from", "-c", "user", "-c", "pass", "-c", "jira", "-s", "settings.yml"},
		},
		{
			name: "mixed legacy and flags in between",
			args: []string{"pac", "-v", "-c", "user", "pass", "jira", "from", "abc123", "def456"},
			want: []string{"pac", "-v", "-c", "user", "-c", "pass", "-c", "jira", "from", "abc123", "def456"},
		},
		{
			name: "single value after -c stays unchanged",
			args: []string{"pac", "from", "-c", "user", "-s", "s.yml"},
			want: []string{"pac", "from", "-c", "user", "-s", "s.yml"},
		},
		{
			name: "-c with no value at end",
			args: []string{"pac", "from", "-c"},
			want: []string{"pac", "from", "-c"},
		},
		{
			name: "legacy -c with positional arg before credentials",
			args: []string{"pac", "from", "HEAD~100", "-c", "token", "target"},
			want: []string{"pac", "from", "HEAD~100", "-c", "token", "-c", "target"},
		},
		{
			name: "empty args",
			args: []string{},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeLegacyCredentialArgs(tt.args)
			if len(got) != len(tt.want) {
				t.Fatalf("length mismatch: got %v, want %v", got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("index %d: got %q, want %q\n  full got:  %v\n  full want: %v", i, got[i], tt.want[i], got, tt.want)
				}
			}
		})
	}
}
