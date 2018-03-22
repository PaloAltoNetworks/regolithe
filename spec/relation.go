package spec

// An Relation represents a specification Relation.
type Relation struct {

	// NOTE: Order of attributes matters!
	// The YAML will be dumped respecting this order.

	RestName     string            `yaml:"rest_name,omitempty"      json:"rest_name,omitempty"`
	Descriptions map[string]string `yaml:"descriptions,omitempty"   json:"descriptions,omitempty"`
	AllowsGet    bool              `yaml:"get,omitempty"            json:"get,omitempty"`
	AllowsCreate bool              `yaml:"create,omitempty"         json:"create,omitempty"`
	AllowsUpdate bool              `yaml:"update,omitempty"         json:"update,omitempty"`
	AllowsDelete bool              `yaml:"delete,omitempty"         json:"delete,omitempty"`
	Deprecated   bool              `yaml:"deprecated,omitempty"     json:"deprecated,omitempty"`

	linkedSpecification *Specification
}

// Specification returns the Specification the API links to.
func (a *Relation) Specification() *Specification {
	return a.linkedSpecification
}
