// Package dto contains all data transfer objects for the application
package model

// Settings represents the root configuration structure for PAC
type Settings struct {
	General     General        `yaml:":general" json:"general"`
	Templates   []Template     `yaml:":templates" json:"templates"`
	TaskSystems []TaskSystem   `yaml:":task_systems" json:"task_systems"`
	VCS         VCS            `yaml:":vcs" json:"vcs"`
	Properties  map[string]any `yaml:":properties" json:"properties"`
	Verbosity   int            `yaml:"verbosity" json:"verbosity"`
}

// General contains global application settings
type General struct {
	Strict bool `yaml:":strict" json:"strict"`
}

// Template represents an output template configuration
type Template struct {
	Location string `yaml:"location" json:"location"`
	Output   string `yaml:"output" json:"output"`
}

// VCS contains version control system configuration
type VCS struct {
	Repo string `yaml:":repo" json:"repo_location"`
}

// NewSettings creates a default Settings struct
func NewSettings() *Settings {
	return &Settings{
		General: General{
			Strict: true,
		},
		Templates:   []Template{},
		TaskSystems: []TaskSystem{},
		VCS: VCS{
			Repo: ".",
		},
		Properties: make(map[string]any),
		Verbosity:  1,
	}
}
