package model

// TaskSystem represents a task tracking system configuration
type TaskSystem struct {
	Name        string      `yaml:":name" json:"name"`
	Token 	 	string      `yaml:":token" json:"token,omitempty"`
	Regex       []RegexRule `yaml:":regex" json:"regex"`
	Delimiter   string      `yaml:":delimiter" json:"delimiter,omitempty"`
	QueryString string      `yaml:":query_string" json:"query_string,omitempty"`
	Username    string      `yaml:":usr" json:"usr,omitempty"`
	Password    string      `yaml:":pw" json:"pw,omitempty"`
}

// RegexRule represents a regex pattern and label for task identification
type RegexRule struct {
	Pattern string `yaml:"pattern" json:"pattern"`
	Label   string `yaml:"label" json:"label"`
}

// NewTaskSystem creates a new TaskSystem with default values
func NewTaskSystem(name string) TaskSystem {
	return TaskSystem{
		Name:  name,
		Regex: []RegexRule{},
	}
}
