package model

// TaskSystem represents a task tracking system configuration
type TaskSystem struct {
	Name        string      `yaml:"name" json:"name"`
	Regex       []RegexRule `yaml:"regex" json:"regex"`
	Delimiter   string      `yaml:"delimiter" json:"delimiter,omitempty"`
	QueryString string      `yaml:"query_string" json:"query_string,omitempty" mapstructure:"query_string"`
}

// RegexRule represents a regex pattern and label for task identification
type RegexRule struct {
	Pattern string `yaml:"pattern" json:"pattern"`
}

// NewTaskSystem creates a new TaskSystem with default values
func NewTaskSystem(name string) TaskSystem {
	return TaskSystem{
		Name:  name,
		Regex: []RegexRule{},
	}
}
