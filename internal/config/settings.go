// Package config handles configuration loading and parsing for PAC.
package config

// Settings represents the complete PAC configuration.
// Supports both legacy Ruby format (with colon prefixes like :general) and
// modern format (without prefixes like general). Keys are normalized during loading.
type Settings struct {
	General     GeneralSettings    `mapstructure:"general"`
	Templates   []TemplateConfig   `mapstructure:"templates"`
	TaskSystems []TaskSystemConfig `mapstructure:"task_systems"`
	VCS         VCSConfig          `mapstructure:"vcs"`
	Properties  map[string]any     `mapstructure:"properties"`
	Verbosity   int                `mapstructure:"-"` // Set from CLI, not config file
}

// GeneralSettings contains general configuration options.
type GeneralSettings struct {
	Strict bool `mapstructure:"strict"`
}

// TemplateConfig defines a single template output.
type TemplateConfig struct {
	Location string `mapstructure:"location"`
	Output   string `mapstructure:"output"`
	PDF      bool   `mapstructure:"pdf"`
}

// TaskSystemConfig configures a task tracking system integration.
type TaskSystemConfig struct {
	Name        string        `mapstructure:"name"`
	QueryString string        `mapstructure:"query_string"`
	Username    string        `mapstructure:"usr"`
	Password    string        `mapstructure:"pw"`
	Regex       []RegexConfig `mapstructure:"regex"`
	Delimiter   string        `mapstructure:"delimiter"`
}

// RegexConfig defines a regex pattern for task ID extraction.
type RegexConfig struct {
	Pattern string `mapstructure:"pattern"`
	Label   string `mapstructure:"label"`
}

// VCSConfig configures the version control system.
type VCSConfig struct {
	Type         string   `mapstructure:"type"`
	RepoLocation string   `mapstructure:"repo_location"`
	FilterPaths  []string `mapstructure:"filter_paths"`
}
