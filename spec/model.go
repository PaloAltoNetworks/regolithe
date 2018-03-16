package spec

// A Model holds generic information about a specification.
type Model struct {
	Aliases          []string `yaml:"aliases,omitempty"`
	AllowsCreate     bool     `yaml:"create,omitempty"`
	AllowsDelete     bool     `yaml:"delete,omitempty"`
	AllowsGet        bool     `yaml:"get,omitempty"`
	AllowsUpdate     bool     `yaml:"update,omitempty"`
	Description      string   `yaml:"description,omitempty"`
	EntityName       string   `yaml:"entity_name,omitempty"`
	EntityNamePlural string   `yaml:"-"`
	Extends          []string `yaml:"extends,omitempty"`
	InstanceName     string   `yaml:"instance_name,omitempty"`
	IsRoot           bool     `yaml:"root,omitempty"`
	Package          string   `yaml:"package,omitempty"`
	ResourceName     string   `yaml:"resource_name,omitempty"`
	RestName         string   `yaml:"rest_name,omitempty"`
	Private          bool     `yaml:"private,omitempty"`
}
