package spec

// A Model holds generic information about a specification.
type Model struct {

	// NOTE: Order of attributes matters!
	// The YAML will be dumped respecting this order.

	RestName     string   `yaml:"rest_name,omitempty"       json:"rest_name,omitempty"`
	ResourceName string   `yaml:"resource_name,omitempty"   json:"resource_name,omitempty"`
	EntityName   string   `yaml:"entity_name,omitempty"     json:"entity_name,omitempty"`
	Package      string   `yaml:"package,omitempty"         json:"package,omitempty"`
	Description  string   `yaml:"description,omitempty"     json:"description,omitempty"`
	Aliases      []string `yaml:"aliases,omitempty"         json:"aliases,omitempty"`
	Private      bool     `yaml:"private,omitempty"         json:"private,omitempty"`

	AllowsCreate bool `yaml:"create,omitempty"  json:"create,omitempty"`
	AllowsGet    bool `yaml:"get,omitempty"     json:"get,omitempty"`
	AllowsUpdate bool `yaml:"update,omitempty"  json:"update,omitempty"`
	AllowsDelete bool `yaml:"delete,omitempty"  json:"delete,omitempty"`

	Extends []string `yaml:"extends,omitempty"  json:"extends,omitempty"`
	IsRoot  bool     `yaml:"root,omitempty"     json:"root,omitempty"`

	EntityNamePlural string `yaml:"-" json:"-"`
}
