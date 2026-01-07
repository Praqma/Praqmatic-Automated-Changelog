// Package config handles configuration loading and parsing for PAC.
package config

// Settings represents the complete PAC configuration.
// Field names use mapstructure tags to match the Ruby YAML format with colon prefixes.
type Settings struct {
	General     GeneralSettings    `mapstructure:":general" yaml:"general"`
	Templates   []TemplateConfig   `mapstructure:":templates" yaml:"templates"`
	TaskSystems []TaskSystemConfig `mapstructure:":task_systems" yaml:"task_systems"`
	VCS         VCSConfig          `mapstructure:":vcs" yaml:"vcs"`
	Properties  map[string]any     `mapstructure:":properties" yaml:"properties"`
	Verbosity   int                `mapstructure:"-"` // Set from CLI, not config file
}

// GeneralSettings contains general configuration options.
type GeneralSettings struct {
	Strict bool `mapstructure:":strict" yaml:"strict"`
}

// TemplateConfig defines a single template output.
type TemplateConfig struct {
	Location string `mapstructure:"location" yaml:"location"`
	Output   string `mapstructure:"output" yaml:"output"`
	PDF      bool   `mapstructure:"pdf" yaml:"pdf"`
}

// TaskSystemConfig configures a task tracking system integration.
type TaskSystemConfig struct {
	Name        string        `mapstructure:":name" yaml:"name"`
	QueryString string        `mapstructure:":query_string" yaml:"query_string"`
	Username    string        `mapstructure:":usr" yaml:"usr"`
	Password    string        `mapstructure:":pw" yaml:"pw"`
	Regex       []RegexConfig `mapstructure:":regex" yaml:"regex"`
	Delimiter   string        `mapstructure:":delimiter" yaml:"delimiter"`
}

// RegexConfig defines a regex pattern for task ID extraction.
type RegexConfig struct {
	Pattern string `mapstructure:"pattern" yaml:"pattern"`
	Label   string `mapstructure:"label" yaml:"label"`
}

// VCSConfig configures the version control system.
type VCSConfig struct {
	Type         string   `mapstructure:":type" yaml:"type"`
	RepoLocation string   `mapstructure:":repo_location" yaml:"repo_location"`
	FilterPaths  []string `mapstructure:":filter_paths" yaml:"filter_paths"`
}
