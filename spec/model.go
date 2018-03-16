package spec

// A Model holds generic information about a specification.
type Model struct {

	// NOTE: Order of attributes matters!
	// The YAML will be dumped respecting this order.

	RestName     string   `yaml:"rest_name,omitempty"`
	ResourceName string   `yaml:"resource_name,omitempty"`
	EntityName   string   `yaml:"entity_name,omitempty"`
	Package      string   `yaml:"package,omitempty"`
	Description  string   `yaml:"description,omitempty"`
	Aliases      []string `yaml:"aliases,omitempty"`
	Private      bool     `yaml:"private,omitempty"`

	AllowsCreate bool `yaml:"create,omitempty"`
	AllowsGet    bool `yaml:"get,omitempty"`
	AllowsUpdate bool `yaml:"update,omitempty"`
	AllowsDelete bool `yaml:"delete,omitempty"`

	Extends []string `yaml:"extends,omitempty"`
	IsRoot  bool     `yaml:"root,omitempty"`

	EntityNamePlural string `yaml:"-"`
}
