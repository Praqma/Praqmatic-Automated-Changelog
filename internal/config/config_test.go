package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSettings_DefaultSettings(t *testing.T) {
	// Find the settings file relative to the project root
	settingsPath := findSettingsFile(t, "settings/default_settings.yml")

	settings, err := LoadSettings(settingsPath, nil)
	if err != nil {
		t.Fatalf("failed to load default_settings.yml: %v", err)
	}

	// Verify general settings
	if settings.General.Strict != false {
		t.Errorf("expected strict=false, got %v", settings.General.Strict)
	}

	// Verify templates are loaded
	if len(settings.Templates) != 2 {
		t.Errorf("expected 2 templates, got %d", len(settings.Templates))
	}

	// Verify task systems are loaded
	if len(settings.TaskSystems) != 3 {
		t.Errorf("expected 3 task systems, got %d", len(settings.TaskSystems))
	}

	// Verify first task system is "none"
	if len(settings.TaskSystems) > 0 && settings.TaskSystems[0].Name != "none" {
		t.Errorf("expected first task system name='none', got %q", settings.TaskSystems[0].Name)
	}

	// Verify VCS settings
	if settings.VCS.Type != "git" {
		t.Errorf("expected vcs type='git', got %q", settings.VCS.Type)
	}

	if settings.VCS.RepoLocation != "" {
		t.Errorf("expected repo_location='', got %q", settings.VCS.RepoLocation)
	}
}

func TestLoadSettings_MinimalSettings(t *testing.T) {
	settingsPath := findSettingsFile(t, "settings/minimal_settings.yml")

	settings, err := LoadSettings(settingsPath, nil)
	if err != nil {
		t.Fatalf("failed to load minimal_settings.yml: %v", err)
	}

	// Verify strict mode is enabled in minimal settings
	if settings.General.Strict != true {
		t.Errorf("expected strict=true, got %v", settings.General.Strict)
	}

	// Verify templates are loaded
	if len(settings.Templates) != 3 {
		t.Errorf("expected 3 templates, got %d", len(settings.Templates))
	}

	// Verify task system
	if len(settings.TaskSystems) != 1 {
		t.Errorf("expected 1 task system, got %d", len(settings.TaskSystems))
	}
}

func TestLoadSettings_TaskSystemRegex(t *testing.T) {
	settingsPath := findSettingsFile(t, "settings/default_settings.yml")

	settings, err := LoadSettings(settingsPath, nil)
	if err != nil {
		t.Fatalf("failed to load settings: %v", err)
	}

	// Find the "none" task system
	noneTS := settings.GetTaskSystemByName("none")
	if noneTS == nil {
		t.Fatal("expected to find 'none' task system")
	}

	// Verify regex patterns are loaded
	if len(noneTS.Regex) != 4 {
		t.Errorf("expected 4 regex patterns for 'none', got %d", len(noneTS.Regex))
	}

	// Verify first regex pattern
	if len(noneTS.Regex) > 0 {
		if noneTS.Regex[0].Pattern != "/Issue:\\s*(\\d+)/i" {
			t.Errorf("unexpected first regex pattern: %q", noneTS.Regex[0].Pattern)
		}
		if noneTS.Regex[0].Label != "none" {
			t.Errorf("expected label='none', got %q", noneTS.Regex[0].Label)
		}
	}
}

func TestLoadSettings_JiraTaskSystem(t *testing.T) {
	settingsPath := findSettingsFile(t, "settings/default_settings.yml")

	settings, err := LoadSettings(settingsPath, nil)
	if err != nil {
		t.Fatalf("failed to load settings: %v", err)
	}

	jiraTS := settings.GetTaskSystemByName("jira")
	if jiraTS == nil {
		t.Fatal("expected to find 'jira' task system")
	}

	// Verify Jira configuration
	if jiraTS.QueryString == "" {
		t.Error("expected jira query_string to be set")
	}

	if jiraTS.Username != "user" {
		t.Errorf("expected jira username='user', got %q", jiraTS.Username)
	}

	if jiraTS.Password != "password" {
		t.Errorf("expected jira password='password', got %q", jiraTS.Password)
	}
}

func TestParseCredentialFlags(t *testing.T) {
	tests := []struct {
		name        string
		credentials []string
		wantLen     int
		wantErr     bool
	}{
		{
			name:        "empty credentials",
			credentials: []string{},
			wantLen:     0,
			wantErr:     false,
		},
		{
			name:        "single credential set",
			credentials: []string{"myuser", "mypass", "jira"},
			wantLen:     1,
			wantErr:     false,
		},
		{
			name:        "multiple credential sets",
			credentials: []string{"user1", "pass1", "jira", "user2", "pass2", "github"},
			wantLen:     2,
			wantErr:     false,
		},
		{
			name:        "incomplete credentials",
			credentials: []string{"user", "pass"},
			wantLen:     0,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creds, err := ParseCredentialFlags(tt.credentials)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCredentialFlags() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(creds) != tt.wantLen {
				t.Errorf("ParseCredentialFlags() returned %d credentials, want %d", len(creds), tt.wantLen)
			}
		})
	}
}

func TestSettings_Validate(t *testing.T) {
	tests := []struct {
		name     string
		settings Settings
		wantErr  bool
	}{
		{
			name: "valid settings",
			settings: Settings{
				Templates: []TemplateConfig{{Location: "template.md"}},
				VCS:       VCSConfig{Type: "git"},
			},
			wantErr: false,
		},
		{
			name: "no templates",
			settings: Settings{
				Templates: []TemplateConfig{},
				VCS:       VCSConfig{Type: "git"},
			},
			wantErr: true,
		},
		{
			name: "template without location",
			settings: Settings{
				Templates: []TemplateConfig{{Output: "out.md"}},
				VCS:       VCSConfig{Type: "git"},
			},
			wantErr: true,
		},
		{
			name: "unsupported vcs type",
			settings: Settings{
				Templates: []TemplateConfig{{Location: "template.md"}},
				VCS:       VCSConfig{Type: "svn"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.settings.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Settings.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// findSettingsFile locates a settings file relative to the project root.
func findSettingsFile(t *testing.T, relativePath string) string {
	t.Helper()

	// Start from current directory and walk up to find the settings folder
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	for {
		candidate := filepath.Join(dir, relativePath)
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("could not find %s in any parent directory", relativePath)
		}
		dir = parent
	}
}
