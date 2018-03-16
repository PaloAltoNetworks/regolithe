package spec

// An Relation represents a specification Relation.
type Relation struct {

	// NOTE: Order of attributes matters!
	// The YAML will be dumped respecting this order.

	RestName     string `yaml:"rest_name,omitempty"`
	AllowsGet    bool   `yaml:"get,omitempty"`
	AllowsCreate bool   `yaml:"create,omitempty"`
	AllowsUpdate bool   `yaml:"update,omitempty"`
	AllowsDelete bool   `yaml:"delete,omitempty"`
	Deprecated   bool   `yaml:"deprecated,omitempty"`

	linkedSpecification *Specification
}

// Specification returns the Specification the API links to.
func (a *Relation) Specification() *Specification {
	return a.linkedSpecification
}
